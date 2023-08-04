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

type Datastore struct {
	tx            chan *KVPair
	done          chan struct{}
	db            *badger.DB
	wait          sync.WaitGroup
	tickerReclaim *time.Ticker
	tickerGC      *time.Ticker
}

type PeriodOption struct {
	MaxSize       int
	DeleteSize    int
	ReclaimPeriod time.Duration
	GCPeriod      time.Duration
}

func NewDatastore(initialEpoch uint64, path string, periodOptions *PeriodOption) *Datastore {
	options := badger.DefaultOptions(path).
		// TODO(fanyang) log with zerolog
		WithLoggingLevel(badger.WARNING).
		WithInMemory(false).
		WithCompression(boptions.ZSTD).
		WithValueLogFileSize(100 * 1000 * 1000) // MB

	db, err := badger.Open(options)
	if err != nil {
		log.Fatal().Err(errors.WithStack(err)).Msg("failed open database")
	}

	log.Info().Uint64("initialEpoch", initialEpoch).Msg("datastore created")
	tx := make(chan *KVPair, 256)

	d := &Datastore{
		done:          make(chan struct{}),
		tx:            tx,
		db:            db,
		tickerReclaim: time.NewTicker(periodOptions.ReclaimPeriod),
		tickerGC:      time.NewTicker(periodOptions.GCPeriod),
	}
	d.wait.Add(3)

	go d.writer(initialEpoch)
	go d.periodicReclaim(periodOptions.MaxSize, periodOptions.DeleteSize)
	go d.periodicGC()
	return d
}

// Tx returns a send-only channel
func (d *Datastore) Tx() chan<- *KVPair {
	return d.tx
}

// Close the datastore and shutdown
func (d *Datastore) Close() {
	close(d.done)
	d.wait.Wait()
}

func (d *Datastore) writer(initialEpoch uint64) {
	log.Info().Msg("datastore writer started")
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
			break
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
		if !strings.HasPrefix(key, PrefixSignalRecord) {
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

func (d *Datastore) GetSize() int {
	count := 0
	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			count++
		}
		return nil
	})
	if err != nil {
		log.Fatal().Err(errors.WithStack(err)).Msg("get size failed")
	}

	return count
}

func (d *Datastore) checkDeletePrefix(prefix []byte, maxSize, deleteSize int) {
	err := d.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		count := 0
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			count++
		}

		if count > maxSize {
			deleted := 0
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				if deleted >= deleteSize || deleted >= count-maxSize {
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
		}

		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("failed to update db")
	}
}

func (d *Datastore) periodicReclaim(maxSize, deleteSize int) {
	log.Info().Msg("datastore periodic delete started")
	defer func(wait *sync.WaitGroup) {
		wait.Done()
		log.Info().Msg("datastore periodic delete exited")
	}(&d.wait)

	maxSize, deleteSize = maxSize/3, deleteSize/3
	for {
		select {
		case <-d.done:
			return
		case <-d.tickerReclaim.C:
			d.checkDeletePrefix([]byte(PrefixNetRecord), maxSize, deleteSize)
			d.checkDeletePrefix([]byte(PrefixNicRecord), maxSize, deleteSize)
			d.checkDeletePrefix([]byte(PrefixTcpRecord), maxSize, deleteSize)
		}
	}
}

func (d *Datastore) periodicGC() {
	log.Info().Msg("datastore periodic GC started")
	defer func(wait *sync.WaitGroup) {
		wait.Done()
		log.Info().Msg("datastore periodic GC exited")
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
