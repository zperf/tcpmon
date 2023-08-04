package tcpmon_test

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
	path          string
	periodOptions *PeriodOption
}

func TestStorage(t *testing.T) {
	s := &StorageTestSuite{
		path: "/tmp/tcpmon-test",
		periodOptions: &PeriodOption{
			MaxSize:       10000,
			DeleteSize:    100000,
			ReclaimPeriod: 2 * time.Second,
			GCPeriod:      5 * time.Minute,
		},
	}
	suite.Run(t, s)
}

func (s *StorageTestSuite) SetupTest() {
	err := os.RemoveAll(s.path)
	s.Assert().NoError(err)
}

func (s *StorageTestSuite) TestGetPrefix() {
	assert := s.Assert()

	ds := NewDatastore(0, s.path, s.periodOptions)
	defer ds.Close()

	tx := ds.Tx()

	for i := 0; i < 3; i++ {
		tx <- &KVPair{
			Key:   NewKey(PrefixNicRecord),
			Value: []byte("test-nic"),
		}
		tx <- &KVPair{
			Key:   NewKey(PrefixTcpRecord),
			Value: []byte("test-tcp"),
		}
		tx <- &KVPair{
			Key:   NewKey(PrefixNetRecord),
			Value: []byte("test-net"),
		}
	}
	s.writeBarrier(tx)

	// check GetPrefix
	pairs, err := ds.GetPrefix([]byte(PrefixNicRecord), 10, false)
	assert.NoError(err)
	assert.Equal(3, len(pairs))
	for _, p := range pairs {
		if !strings.HasPrefix(p.Key, PrefixNicRecord) {
			assert.Failf("Key: %s don't have the nic prefix", p.Key)
		}
	}
	pairs, err = ds.GetPrefix([]byte(PrefixTcpRecord), 10, false)
	assert.NoError(err)
	assert.Equal(3, len(pairs))
	for _, p := range pairs {
		if !strings.HasPrefix(p.Key, PrefixTcpRecord) {
			assert.Failf("Key: %s don't have the tcp prefix", p.Key)
		}
	}
	pairs, err = ds.GetPrefix([]byte(PrefixNetRecord), 10, false)
	assert.NoError(err)
	assert.Equal(3, len(pairs))
	for _, p := range pairs {
		if !strings.HasPrefix(p.Key, PrefixNetRecord) {
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

	ds := NewDatastore(0, s.path, s.periodOptions)
	defer ds.Close()

	tx := ds.Tx()

	for i := 0; i < 3; i++ {
		tx <- &KVPair{
			Key:   NewKey(PrefixNicRecord),
			Value: nil,
		}
	}
	s.writeBarrier(tx)

	keys, err := ds.GetKeys()
	assert.NoError(err)
	assert.Equal(3, len(keys))
}

func (s *StorageTestSuite) TestPeriodicReclaim() {
	ds := NewDatastore(0, s.path, s.periodOptions)
	defer ds.Close()

	tx := ds.Tx()
	for i := 0; i < s.periodOptions.MaxSize; i++ {
		tx <- &KVPair{
			Key:   NewKey(PrefixNicRecord),
			Value: nil,
		}
		tx <- &KVPair{
			Key:   NewKey(PrefixNetRecord),
			Value: nil,
		}
		tx <- &KVPair{
			Key:   NewKey(PrefixTcpRecord),
			Value: nil,
		}
	}
	s.writeBarrier(tx)

	size := ds.GetSize()
	log.Trace().Int("size", size).Msg("insert")

	// wait for reclaim trigger
	time.Sleep(s.periodOptions.ReclaimPeriod + time.Second)

	size = ds.GetSize()
	log.Info().Int("size", size).Msg("reclaim done")

	s.Assert().GreaterOrEqual(s.periodOptions.MaxSize, size)
}

// writeBarrier waits for write complete
func (s *StorageTestSuite) writeBarrier(tx chan<- *KVPair) {
	var wg sync.WaitGroup
	wg.Add(1)
	tx <- &KVPair{
		Key:   NewKey(PrefixSignalRecord),
		Value: nil,
		Callback: func(err error) {
			s.Assert().NoError(err)
			wg.Done()
		},
	}
	wg.Wait()
}
