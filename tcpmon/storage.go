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
	db       *badger.DB
	tx       chan *KVPair
	config   DataStoreConfig
	lastOpen time.Time

	done     chan struct{}
	waitExit sync.WaitGroup // wait for all goroutines exit
}

type DataStoreConfig struct {
	Path            string
	MaxSize         uint32
	WriteInterval   time.Duration
	ExpectedRss     uint64
	MinOpenInterval time.Duration
}

func (c *DataStoreConfig) MaxSizePerType() uint32 {
	return c.MaxSize / MetricTypeCount
}

func NewDataStore(initialEpoch uint64, config *DataStoreConfig) *DataStore {
	if config == nil {
		log.Fatal().Msgf("Config is nil")
	}

	log.Info().Uint64("initialEpoch", initialEpoch).Msg("Created")
	tx := make(chan *KVPair, 256)
	d := &DataStore{
		done:   make(chan struct{}),
		config: *config,
		tx:     tx,
	}
	d.waitExit.Add(1)

	go d.writer(initialEpoch, config.WriteInterval)
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
	log.Info().Msg("Writer started")
	defer func(wait *sync.WaitGroup, db *badger.DB) {
		err := db.Close()
		if err != nil {
			log.Warn().Err(err).Msg("close db failed")
		}
		wait.Done()
		log.Info().Msg("datastore writer exited")
	}(&d.waitExit, d.db)

	maxCountPerType := d.config.MaxSizePerType()
	epoch := initialEpoch
	ticker := time.NewTicker(writeInterval)

	d.openDatabase()

	for {
		select {
		case now := <-ticker.C:
			log.Trace().Time("now", now).Msg("Write trigger")
			err := d.doWrite(&epoch, maxCountPerType)
			if err != nil {
				log.Warn().Err(errors.WithStack(err)).Msg("Write failed")
			}

			// check memory usage and reopen db if needed
			var memstats runtime.MemStats
			runtime.ReadMemStats(&memstats)

			coolDownReady := d.lastOpen.Add(d.config.MinOpenInterval).Before(time.Now())
			if coolDownReady && memstats.Sys > d.config.ExpectedRss {
				d.reopenDatabase()
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
		return errors.WithStack(err)
	}
	netCount += p

	p, err = d.GetTcpKeyCount()
	if err != nil {
		return errors.WithStack(err)
	}
	tcpCount += p

	p, err = d.GetNicKeyCount()
	if err != nil {
		return errors.WithStack(err)
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
	err = d.db.Update(func(txn *badger.Txn) error {
		// inserts
		for _, req := range toWrite {
			key := fmt.Sprintf("%s%v", req.Key, epoch)
			*epoch++

			err := txn.Set([]byte(key), req.Value)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		// reclaim
		err = txnDeleteOldestByPrefix(txn, []byte(PrefixTcpMetric), deleteTcpCount)
		if err != nil {
			return errors.WithStack(err)
		}
		err = txnDeleteOldestByPrefix(txn, []byte(PrefixNicMetric), deleteTcpCount)
		if err != nil {
			return errors.WithStack(err)
		}
		err = txnDeleteOldestByPrefix(txn, []byte(PrefixNetMetric), deleteTcpCount)
		if err != nil {
			return errors.WithStack(err)
		}

		// update count
		err = txnSetUint32(txn, KeyTcpCount, tcpCount-deleteTcpCount)
		if err != nil {
			return errors.WithStack(err)
		}
		err = txnSetUint32(txn, KeyNicCount, nicCount-deleteNicCount)
		if err != nil {
			return errors.WithStack(err)
		}
		err = txnSetUint32(txn, KeyNetCount, netCount-deleteNetCount)
		if err != nil {
			return errors.WithStack(err)
		}
		err = txnSetUint32(txn, KeyTotalCount, totalCount-deleteTotalCount)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return errors.WithStack(err)
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
		WithZSTDCompressionLevel(2).
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
	return txn.Set([]byte(key), []byte(strconv.FormatUint(uint64(value), 10)))
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
