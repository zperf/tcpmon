package tcpmon

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/dgraph-io/badger/v4"
	boptions "github.com/dgraph-io/badger/v4/options"
	"github.com/rs/zerolog/log"
)

type DataStore struct {
	tx            chan *KVPair
	done          chan struct{}
	db            *badger.DB
	wait          sync.WaitGroup
	tickerReclaim *time.Ticker
	tickerGC      *time.Ticker
}

type DataStoreConfig struct {
	Path            string
	MaxSize         int
	ReclaimBatch    int
	ReclaimInterval time.Duration
	GcInterval      time.Duration
}

func (c *DataStoreConfig) WithDefault() {
	if c.MaxSize == 0 {
		c.MaxSize = 30000
	}
	if c.ReclaimBatch == 0 {
		c.ReclaimBatch = 2000
	}
	if c.ReclaimInterval == 0 {
		c.ReclaimInterval = time.Minute
	}
	if c.GcInterval == 0 {
		c.GcInterval = 5 * time.Minute
	}
}

func NewDataStore(initialEpoch uint64, config *DataStoreConfig) *DataStore {
	config.WithDefault()

	options := badger.DefaultOptions(config.Path).
		WithLogger(NewBadgerLogger()).
		WithInMemory(false).
		WithCompression(boptions.ZSTD).
		WithNumGoroutines(2).
		WithNumMemtables(1).
		WithMemTableSize(8 << 20).
		WithBlockSize(4 << 20).
		WithCompactL0OnClose(true)

	db, err := badger.Open(options)
	if err != nil {
		log.Fatal().Err(errors.WithStack(err)).Msg("failed open database")
	}

	log.Info().Uint64("initialEpoch", initialEpoch).Msg("Created")
	tx := make(chan *KVPair, 256)

	d := &DataStore{
		done:          make(chan struct{}),
		tx:            tx,
		db:            db,
		tickerReclaim: time.NewTicker(config.ReclaimInterval),
		tickerGC:      time.NewTicker(config.GcInterval),
	}
	d.wait.Add(3)

	go d.writer(initialEpoch)
	go d.autoReclaim(config.MaxSize, config.ReclaimBatch)
	go d.autoGC()
	return d
}

// Tx returns a send-only channel
func (d *DataStore) Tx() chan<- *KVPair {
	return d.tx
}

// Close the datastore and shutdown
func (d *DataStore) Close() {
	close(d.done)
	d.wait.Wait()
}

func (d *DataStore) writer(initialEpoch uint64) {
	log.Info().Msg("Writer started")
	defer func(wait *sync.WaitGroup, db *badger.DB) {
		err := db.Close()
		if err != nil {
			log.Warn().Err(err).Msg("close db failed")
		}
		wait.Done()
		log.Info().Msg("datastore writer exited")
	}(&d.wait, d.db)

	epoch := initialEpoch
	for {
		var req *KVPair
		select {
		case req = <-d.tx:
		case <-d.done:
			return
		}

		if len(req.Key) <= 0 {
			log.Warn().Str("Key", req.Key).Int("ValueLen", len(req.Value)).
				Msg("ignore invalid KVPair")
			continue
		}

		key := fmt.Sprintf("%s%v", req.Key, epoch)
		req.Key = key
		epoch++

		var err error
		if !strings.HasPrefix(key, PrefixSignal) {
			// TODO(fanyang) batch write txn
			err = d.db.Update(func(txn *badger.Txn) error {
				return txn.Set([]byte(req.Key), req.Value)
			})
			if err != nil {
				log.Warn().Err(err).Msg("failed to insert db")
			}
		}
		if req.Callback != nil {
			req.Callback(err)
		}
	}
}

func (d *DataStore) GetSize(prefix []byte) int {
	size := 0
	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			size++
		}
		return nil
	})
	if err != nil {
		log.Fatal().Err(errors.WithStack(err)).Msg("get size failed")
	}

	return size
}

func (d *DataStore) checkDeletePrefix(prefix []byte, maxSize int, deleteSize int) {
	size := d.GetSize(prefix)
	if size <= maxSize {
		return
	}

	err := d.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		deleted := 0
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			if deleted >= deleteSize || deleted >= size-maxSize {
				break
			}

			item := it.Item()
			key := item.KeyCopy(nil)
			err := txn.Delete(key)
			if err != nil {
				log.Warn().Err(err).Str("key", string(key)).Msg("failed to delete item")
			}
			deleted++
		}

		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("failed to delete")
	}
}

func (d *DataStore) autoReclaim(maxSize, deleteSize int) {
	defer func(wait *sync.WaitGroup) {
		wait.Done()
	}(&d.wait)

	maxSize, deleteSize = maxSize/MetricTypeCount, deleteSize/MetricTypeCount
	for {
		select {
		case <-d.done:
			return
		case <-d.tickerReclaim.C:
			d.checkDeletePrefix([]byte(PrefixNetMetric), maxSize, deleteSize)
			d.checkDeletePrefix([]byte(PrefixNicMetric), maxSize, deleteSize)
			d.checkDeletePrefix([]byte(PrefixTcpMetric), maxSize, deleteSize)
		}
	}
}

func (d *DataStore) autoGC() {
	defer func(wait *sync.WaitGroup) {
		wait.Done()
	}(&d.wait)

	for {
		select {
		case <-d.done:
			return
		case <-d.tickerGC.C:
			_ = d.db.RunValueLogGC(0.5)
		}
	}
}
