package tcpmon

import (
	"fmt"
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

	for {
		select {
		case <-d.done:
			return
		case now := <-ticker.C:
			d.openDatabase()
			log.Trace().Time("now", now).Msg("Write trigger")
			err := d.doWrite(&epoch, maxCountPerType)
			if err != nil {
				log.Warn().Err(errors.WithStack(err)).Msg("Write failed")
			}
			d.closeDatabase()
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
}

func (d *DataStore) closeDatabase() {
	if d.db != nil {
		err := d.db.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Close db failed")
		}
		d.db = nil
	}
	debug.FreeOSMemory()
}
