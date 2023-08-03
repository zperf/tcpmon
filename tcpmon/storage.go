package tcpmon

import (
	"fmt"
	"strings"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/dgraph-io/badger/v4"
	boptions "github.com/dgraph-io/badger/v4/options"
	"github.com/rs/zerolog/log"
)

type Datastore struct {
	tx   chan *KVPair
	done chan struct{}
	db   *badger.DB
	wait sync.WaitGroup
}

func NewDatastore(initialEpoch uint64, path string) *Datastore {
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
		done: make(chan struct{}),
		tx:   tx,
		db:   db,
	}
	d.wait.Add(1)

	go d.writer(initialEpoch)
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
	defer func(wait *sync.WaitGroup) {
		wait.Done()
		log.Info().Msg("datastore writer exited")
	}(&d.wait)

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
