package test

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	v1 "github.com/zperf/tcpmon/tcpmon/storage/v1"
)

type StorageTestSuite struct {
	suite.Suite
	config *v1.DataStoreConfig
}

func TestStorage(t *testing.T) {
	s := &StorageTestSuite{
		config: &v1.DataStoreConfig{
			Path:          "./tcpmon-test",
			MaxSize:       10000,
			WriteInterval: 10 * time.Second,
		},
	}
	suite.Run(t, s)
}

func (s *StorageTestSuite) SetupTest() {
	err := os.RemoveAll(s.config.Path)
	s.Assert().NoError(err)
}

func (s *StorageTestSuite) TestDelayedWriterWorking() {
	ds := v1.NewDataStore(0, s.config)
	defer ds.Close()

	now := time.Now()

	tx := ds.Tx()
	tx <- v1.NewMetricContext(v1.NewKey(v1.PrefixNetMetric)+"1", nil)
	tx <- v1.NewMetricContext(v1.NewKey(v1.PrefixNetMetric)+"2", nil)
	tx <- v1.NewMetricContext(v1.NewKey(v1.PrefixNetMetric)+"3", nil)
	s.writeBarrier(tx)

	elapsed := time.Since(now)
	s.Assert().GreaterOrEqual(elapsed, s.config.WriteInterval, "Writing should take more than 10s")
	s.Assert().Less(elapsed, s.config.WriteInterval*2)
}

//func printMemStats(msg string) {
//	var stat runtime.MemStats
//	runtime.ReadMemStats(&stat)
//	log.Info().Float32("alloc", float32(stat.Alloc)/1024/1024).
//		Float32("sys", float32(stat.Sys)/1024/1024).
//		Float32("frees", float32(stat.Frees)/1024/1024).Msg(msg)
//}

//func (s *StorageTestSuite) TestMemoryUsage() {
//	printMemStats("Start")
//	ds := v1.NewDataStore(0, s.config)
//
//	for i := 0; i < 10; i++ {
//		const writeCount = 180
//		tx := ds.Tx()
//		printMemStats("Before write")
//
//		for i := 0; i < writeCount; i++ {
//			tx <- v1.NewMetricContext(fmt.Sprintf("%s/%v", v1.NewKey(v1.PrefixNetMetric), i), randBuf(1<<20))
//			if i == writeCount/2 {
//				printMemStats("During write")
//			}
//		}
//		s.writeBarrier(tx)
//		debug.FreeOSMemory()
//		printMemStats("After write")
//	}
//	ds.Close()
//	printMemStats("DB closed")
//}

// writeBarrier waits for write complete
func (s *StorageTestSuite) writeBarrier(tx chan<- *v1.MetricContext) {
	var wg sync.WaitGroup
	wg.Add(1)
	tx <- &v1.MetricContext{
		Key:   v1.NewKey(v1.PrefixSignal),
		Value: nil,
		Callback: func() {
			wg.Done()
		},
	}
	wg.Wait()
}
