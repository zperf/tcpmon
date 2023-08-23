package test

import (
	"os"
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
			ReclaimBatch:    100000,
			ReclaimInterval: 2 * time.Second,
			GcInterval:      5 * time.Minute,
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

	// check GetPrefix all
	pairs, err = ds.GetPrefix([]byte{}, 10, true)
	assert.NoError(err)
	assert.Equal(9, len(pairs))
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

	keys, err := ds.GetKeys()
	assert.NoError(err)
	assert.Equal(3, len(keys))
}

func (s *StorageTestSuite) TestPeriodicReclaim() {
	ds := NewDataStore(0, s.config)
	defer ds.Close()

	tx := ds.Tx()
	for i := 0; i < s.config.MaxSize; i++ {
		tx <- &KVPair{
			Key:   NewKey(PrefixNicMetric),
			Value: nil,
		}
		tx <- &KVPair{
			Key:   NewKey(PrefixNetMetric),
			Value: nil,
		}
		tx <- &KVPair{
			Key:   NewKey(PrefixTcpMetric),
			Value: nil,
		}
	}
	s.writeBarrier(tx)

	size := ds.GetSize(nil)
	log.Trace().Int("size", size).Msg("insert")

	// wait for reclaim trigger
	time.Sleep(s.config.ReclaimInterval + time.Second)

	size = ds.GetSize(nil)
	log.Info().Int("size", size).Msg("reclaim done")

	s.Assert().GreaterOrEqual(s.config.MaxSize, size)
}

// writeBarrier waits for write complete
func (s *StorageTestSuite) writeBarrier(tx chan<- *KVPair) {
	var wg sync.WaitGroup
	wg.Add(1)
	tx <- &KVPair{
		Key:   NewKey(PrefixSignal),
		Value: nil,
		Callback: func(err error) {
			s.Assert().NoError(err)
			wg.Done()
		},
	}
	wg.Wait()
}
