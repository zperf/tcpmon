package test

import (
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"

	. "github.com/zperf/tcpmon/tcpmon"
)

type StorageTestSuite struct {
	suite.Suite
	config *DataStoreConfig
}

func TestStorage(t *testing.T) {
	s := &StorageTestSuite{
		config: &DataStoreConfig{
			Path:            "/tmp/tcpmon-test",
			MaxSize:         10000,
			WriteInterval:   time.Second,
			ExpectedRatio:   200 << 20,
			MinOpenInterval: 10 * time.Second,
		},
	}
	suite.Run(t, s)
}

func (s *StorageTestSuite) SetupTest() {
	err := os.RemoveAll(s.config.Path)
	s.Assert().NoError(err)
}

func (s *StorageTestSuite) TestGetPrefix() {
	assert := s.Assert()

	ds := NewDataStore(0, s.config)
	defer ds.Close()

	tx := ds.Tx()

	for i := 0; i < 3; i++ {
		tx <- &KVPair{
			Key:   NewKey(PrefixNicMetric),
			Value: []byte("test-nic"),
		}
		tx <- &KVPair{
			Key:   NewKey(PrefixTcpMetric),
			Value: []byte("test-tcp"),
		}
		tx <- &KVPair{
			Key:   NewKey(PrefixNetMetric),
			Value: []byte("test-net"),
		}
	}
	s.writeBarrier(tx)

	// check GetPrefix
	pairs, err := ds.GetPrefix([]byte(PrefixNicMetric), 10, false)
	assert.NoError(err)
	assert.Equal(3, len(pairs))
	for _, p := range pairs {
		if !strings.HasPrefix(p.Key, PrefixNicMetric) {
			assert.Failf("Key: %s don't have the nic prefix", p.Key)
		}
	}
	pairs, err = ds.GetPrefix([]byte(PrefixTcpMetric), 10, false)
	assert.NoError(err)
	assert.Equal(3, len(pairs))
	for _, p := range pairs {
		if !strings.HasPrefix(p.Key, PrefixTcpMetric) {
			assert.Failf("Key: %s don't have the tcp prefix", p.Key)
		}
	}
	pairs, err = ds.GetPrefix([]byte(PrefixNetMetric), 10, false)
	assert.NoError(err)
	assert.Equal(3, len(pairs))
	for _, p := range pairs {
		if !strings.HasPrefix(p.Key, PrefixNetMetric) {
			assert.Failf("Key: %s don't have the net prefix", p.Key)
		}
	}
}

func (s *StorageTestSuite) TestGetKeys() {
	assert := s.Assert()

	ds := NewDataStore(0, s.config)
	defer ds.Close()

	tx := ds.Tx()

	for i := 0; i < 3; i++ {
		tx <- &KVPair{
			Key:   NewKey(PrefixNicMetric),
			Value: nil,
		}
	}
	s.writeBarrier(tx)

	keys, err := ds.GetPrefix([]byte(PrefixNicMetric), 0, false)
	assert.NoError(err)
	assert.Equal(3, len(keys))
}

func (s *StorageTestSuite) TestMemoryFootprint() {
	ds := NewDataStore(0, s.config)
	tx := ds.Tx()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	log.Info().Float32("memStats.Sys(MiB)", float32(memStats.Sys)/(1<<20)).
		Float32("memStats.Alloc(MiB)", float32(memStats.Alloc)/(1<<20)).
		Msg("Memory footprint, before write")

	// for a day
	for i := 0; i < 3600; i++ {
		tx <- &KVPair{
			Key:   NewKey(PrefixNicMetric),
			Value: make([]byte, 4096),
		}
		tx <- &KVPair{
			Key:   NewKey(PrefixTcpMetric),
			Value: make([]byte, 4096),
		}
		tx <- &KVPair{
			Key:   NewKey(PrefixNetMetric),
			Value: make([]byte, 4096),
		}
	}
	s.writeBarrier(tx)

	runtime.ReadMemStats(&memStats)
	log.Info().Float32("memStats.Sys(MiB)", float32(memStats.Sys)/(1<<20)).
		Float32("memStats.Alloc(MiB)", float32(memStats.Alloc)/(1<<20)).
		Msg("Memory footprint, after write")

	_, err := ds.GetKeys()
	s.Require().NoError(err)

	runtime.ReadMemStats(&memStats)
	log.Info().Float32("memStats.Sys(MiB)", float32(memStats.Sys)/(1<<20)).
		Float32("memStats.Alloc(MiB)", float32(memStats.Alloc)/(1<<20)).
		Msg("Memory footprint, after read")
}

// writeBarrier waits for write complete
func (s *StorageTestSuite) writeBarrier(tx chan<- *KVPair) {
	var wg sync.WaitGroup
	wg.Add(1)
	tx <- &KVPair{
		Key:   NewKey(PrefixSignal),
		Value: nil,
		Callback: func() {
			wg.Done()
		},
	}
	wg.Wait()
}
