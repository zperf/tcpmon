package tcpmon

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/dgraph-io/badger/v4"
	boptions "github.com/dgraph-io/badger/v4/options"
	"github.com/rs/zerolog/log"
)

type DataStore struct {
	db     *badger.DB
	tx     chan *KVPair
	config DataStoreConfig
	lastOpen time.Time

	done       chan struct{}
	waitExit   sync.WaitGroup // wait for all goroutines exit
	waitDbInit sync.WaitGroup
}

type DataStoreConfig struct {
	Path            string
	MaxSize         uint32
	WriteInterval   time.Duration
	ExpectedRatio   float32
	MinOpenInterval time.Duration
}

func (c *DataStoreConfig) MaxSizePerType() uint32 {
	return c.MaxSize / MetricTypeCount
}

func NewDataStore(initialEpoch uint64, config *DataStoreConfig) *DataStore {
	if config == nil {
		log.Fatal().Msgf("Config is nil")
		return nil // make linter happy
	}

	log.Info().Uint64("initialEpoch", initialEpoch).Msg("Created")
	tx := make(chan *KVPair, 256)
	d := &DataStore{
		done:   make(chan struct{}),
		config: *config,
		tx:     tx,
	}
	d.waitExit.Add(1)
	d.waitDbInit.Add(1)

	go d.writer(initialEpoch, config.WriteInterval)

	d.waitDbInit.Wait()
	return d
}

// Tx returns a send-only channel
func (d *DataStore) Tx() chan<- *KVPair {
	return d.tx
}

// Close the datastore and shutdown
func (d *DataStore) Close() {
	close(d.done)
	d.waitExit.Wait()
}

func (d *DataStore) writer(initialEpoch uint64, writeInterval time.Duration) {
	log.Info().Msg("DataStore writer started")
	defer func(d *DataStore) {
		if d.db != nil {
			err := d.db.Close()
			if err != nil {
				log.Warn().Err(err).Msg("Close db failed")
			}
		}
		log.Info().Msg("DataStore writer exited")
		d.waitExit.Done()
	}(d)

	d.openDatabase()
	err := d.db.Update(func(txn *badger.Txn) error {
		err := txnSetUint32(txn, KeyTcpCount, 0)
		if err != nil {
			return errors.WithStack(err)
		}
		err = txnSetUint32(txn, KeyNicCount, 0)
		if err != nil {
			return errors.WithStack(err)
		}
		err = txnSetUint32(txn, KeyNetCount, 0)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("db initialize write failed")
	}
	d.closeDatabase()

	maxCountPerType := d.config.MaxSizePerType()
	epoch := initialEpoch
	ticker := time.NewTicker(writeInterval)

	d.openDatabase()
	d.waitDbInit.Done()

	err := d.db.Update(func(txn *badger.Txn) error {
		err := txnEnsureExistsUint32(txn, KeyTotalCount)
		if err != nil {
			return err
		}

		err = txnEnsureExistsUint32(txn, KeyTcpCount)
		if err != nil {
			return err
		}

		err = txnEnsureExistsUint32(txn, KeyNetCount)
		if err != nil {
			return err
		}

		err = txnEnsureExistsUint32(txn, KeyNicCount)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("ensure metadata exists failed")
	}

	for {
		select {
		case <-d.done:
			return
		case now := <-ticker.C:
			log.Trace().Time("now", now).Msg("Write trigger")
			err := d.doWrite(&epoch, maxCountPerType)
			if err != nil {
				log.Warn().Err(err).Msg("Write failed")
			}

			coolDownReady := d.lastOpen.Add(d.config.MinOpenInterval).Before(time.Now())
			if coolDownReady {
				// check memory usage and reopen db if needed
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)

				ratio := 1 - float32(memStats.Alloc)/float32(memStats.Sys)
				if ratio-d.config.ExpectedRatio >= float32(1e-9) {
					log.Info().Float32("memStats.Sys(MiB)", float32(memStats.Sys)/(1<<20)).
						Float32("memStats.Alloc(MiB)", float32(memStats.Alloc)/(1<<20)).
						Float32("ratio", ratio).
						Msg("reopen database")
					d.reopenDatabase()
				}
			}
		}
	}
}

func (d *DataStore) doWrite(epoch *uint64, maxCountPerType uint32) error {
	// reap requests from queue
	toWrite := make([]*KVPair, 0)
reap:
	for {
		select {
		case p := <-d.tx:
			toWrite = append(toWrite, p)
		default:
			break reap
		}
	}

	// stat the queue
	tcpCount, nicCount, netCount, totalCount := uint32(0), uint32(0), uint32(0), uint32(0)
	for _, p := range toWrite {
		if strings.HasPrefix(p.Key, PrefixTcpMetric) {
			tcpCount++
		} else if strings.HasPrefix(p.Key, PrefixNetMetric) {
			netCount++
		} else if strings.HasPrefix(p.Key, PrefixNicMetric) {
			nicCount++
		}
		totalCount++
	}

	p, err := d.GetNetKeyCount()
	if err != nil {
		return err
	}
	netCount += p

	p, err = d.GetTcpKeyCount()
	if err != nil {
		return err
	}
	tcpCount += p

	p, err = d.GetNicKeyCount()
	if err != nil {
		return err
	}
	nicCount += p

	deleteNicCount := uint32(0)
	if nicCount > maxCountPerType {
		deleteNicCount = nicCount - maxCountPerType
	}

	deleteTcpCount := uint32(0)
	if tcpCount > maxCountPerType {
		deleteTcpCount = tcpCount - maxCountPerType
	}

	deleteNetCount := uint32(0)
	if netCount > maxCountPerType {
		deleteNetCount = netCount - maxCountPerType
	}

	deleteTotalCount := deleteNicCount + deleteTcpCount + deleteNetCount

	// submit txn
	signals := make([]int, 0)
	err = d.db.Update(func(txn *badger.Txn) error {
		// inserts
		for i, req := range toWrite {
			if req.IsSignal() {
				signals = append(signals, i)
				continue
			}

			key := fmt.Sprintf("%s%v", req.Key, *epoch)
			*epoch++

			err := txn.Set([]byte(key), req.Value)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		// reclaim
		err = txnDeleteOldestByPrefix(txn, []byte(PrefixTcpMetric), deleteTcpCount)
		if err != nil {
			return err
		}
		err = txnDeleteOldestByPrefix(txn, []byte(PrefixNicMetric), deleteTcpCount)
		if err != nil {
			return err
		}
		err = txnDeleteOldestByPrefix(txn, []byte(PrefixNetMetric), deleteTcpCount)
		if err != nil {
			return err
		}

		// update count
		err = txnSetUint32(txn, KeyTcpCount, tcpCount-deleteTcpCount)
		if err != nil {
			return err
		}
		err = txnSetUint32(txn, KeyNicCount, nicCount-deleteNicCount)
		if err != nil {
			return err
		}
		err = txnSetUint32(txn, KeyNetCount, netCount-deleteNetCount)
		if err != nil {
			return err
		}
		err = txnSetUint32(txn, KeyTotalCount, totalCount-deleteTotalCount)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	for _, i := range signals {
		req := toWrite[i]
		if req.Callback != nil {
			req.Callback()
		}
	}

	return nil
}

func (d *DataStore) openDatabase() {
	if d.db != nil {
		log.Fatal().Msg("db should be nil before open it")
	}

	options := badger.DefaultOptions(d.config.Path).
		WithLogger(NewBadgerLogger()).
		WithCompression(boptions.ZSTD).
		WithZSTDCompressionLevel(1).
		WithNumMemtables(1)

	db, err := badger.Open(options)
	if err != nil {
		log.Fatal().Err(errors.WithStack(err)).Msg("failed open database")
	}
	d.db = db
	d.lastOpen = time.Now()
}

func (d *DataStore) reopenDatabase() {
	if d.db != nil {
		err := d.db.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Close db failed")
		}
		d.db = nil
	}

	debug.FreeOSMemory()
	d.openDatabase()
}

func txnSetUint32(txn *badger.Txn, key string, value uint32) error {
	err := txn.Set([]byte(key), []byte(strconv.FormatUint(uint64(value), 10)))
	return errors.Wrap(err, "txn set uint32 failed")
}

func txnDeleteOldestByPrefix(txn *badger.Txn, prefix []byte, deleteCount uint32) error {
	options := badger.DefaultIteratorOptions
	options.Prefix = prefix
	options.PrefetchValues = false
	itr := txn.NewIterator(options)
	defer itr.Close()

	count := uint32(0)
	for itr.Seek(prefix); itr.ValidForPrefix(prefix); itr.Next() {
		if count >= deleteCount {
			break
		}

		el := itr.Item()
		err := txn.Delete(el.Key())
		if err != nil {
			return errors.WithStack(err)
		}
		count++
	}

	return nil
}

func txnEnsureExistsUint32(txn *badger.Txn, key string) error {
	_, err := txn.Get([]byte(key))
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			err = txnSetUint32(txn, key, 0)
			if err != nil {
				return err
			}
		} else {
			return errors.Wrapf(err, "get %s failed", KeyTcpCount)
		}
	}
	return nil
}
