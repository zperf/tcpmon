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

const PrefixSocketRecord = "s"
const PrefixIfaceRecord = "if"

// SEP Separator between key segments. eg.
const SEP = "/"

func pathJoin(s ...string) string {
	return strings.Join(s, SEP)
}

type Datastore struct {
	tx   chan *StoreRequest
	done chan struct{}
	wait sync.WaitGroup
}

type StoreRequest struct {
	Key   string
	Value []byte
}

func NewDatastore(epoch uint64) *Datastore {
	log.Info().Uint64("epoch", epoch).Msg("datastore created")
	tx := make(chan *StoreRequest, 256)
	d := &Datastore{
		done: make(chan struct{}),
		tx:   tx,
	}
	d.wait.Add(1)
	go d.writer(epoch)
	return d
}

// Tx returns a send-only channel
func (d *Datastore) Tx() chan<- *StoreRequest {
	return d.tx
}

func (d *Datastore) Close() {
	close(d.done)
	d.wait.Wait()
}

func (d *Datastore) writer(initialEpoch uint64) {
	epoch := initialEpoch

	// TODO(fanyang) options for the db path
	options := badger.DefaultOptions("/tmp/tcpmon").
		// TODO(fanyang) log with zerolog
		WithLoggingLevel(badger.WARNING).
		WithInMemory(false).
		WithCompression(boptions.ZSTD).
		WithValueLogFileSize(100 * 1000 * 1000) // MB

	db, err := badger.Open(options)
	if err != nil {
		log.Fatal().Err(errors.WithStack(err)).Msg("failed open database")
	}
	defer func() {
		err := db.Close()
		log.Info().Err(err).Msg("db closed")
		d.wait.Done()
	}()

	log.Info().Msg("datastore writer started")
	for {
		var req *StoreRequest
		select {
		case req = <-d.tx:
			break
		case <-d.done:
			return
		}

		if len(req.Key) <= 0 || len(req.Value) <= 0 {
			log.Warn().Str("Key", req.Key).Int("ValueLen", len(req.Value)).
				Msg("ignore invalid StoreRequest")
			continue
		}

		req.Key = fmt.Sprintf("%s%v", req.Key, epoch)
		epoch++

		// TODO(fanyang) batch write txn
		err := db.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(req.Key), req.Value)
		})
		if err != nil {
			log.Warn().Err(err).Msg("failed to insert db")
		}
	}
}
