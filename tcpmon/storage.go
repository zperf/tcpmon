package tcpmon

import (
	"fmt"
	"strings"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/dgraph-io/badger/v4"
	boptions "github.com/dgraph-io/badger/v4/options"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/ugurcsen/gods-generic/queues"
	"github.com/ugurcsen/gods-generic/queues/circularbuffer"
	"google.golang.org/protobuf/proto"
)

const PrefixTcpRecord = "tcp"
const PrefixNicRecord = "nic"
const PrefixNetRecord = "net"

type Datastore struct {
	tx   chan *KVPair
	done chan struct{}
	db   *badger.DB

	// peek window
	windowMu   sync.RWMutex
	windowSize int
	window     queues.Queue[string]

	wait sync.WaitGroup
}

type KVPair struct {
	Key   string
	Value []byte
}

func NewDatastore(epoch uint64, windowSize int) *Datastore {
	options := badger.DefaultOptions(viper.GetString("db")).
		// TODO(fanyang) log with zerolog
		WithLoggingLevel(badger.WARNING).
		WithInMemory(false).
		WithCompression(boptions.ZSTD).
		WithValueLogFileSize(100 * 1000 * 1000) // MB

	db, err := badger.Open(options)
	if err != nil {
		log.Fatal().Err(errors.WithStack(err)).Msg("failed open database")
	}

	log.Info().Uint64("epoch", epoch).Msg("datastore created")
	tx := make(chan *KVPair, 256)

	d := &Datastore{
		done:       make(chan struct{}),
		tx:         tx,
		db:         db,
		windowSize: windowSize,
		window:     circularbuffer.New[string](windowSize),
	}
	d.wait.Add(1)
	go d.writer(epoch, d.window)
	return d
}

// Tx returns a send-only channel
func (d *Datastore) Tx() chan<- *KVPair {
	return d.tx
}

func (d *Datastore) Close() {
	close(d.done)
	d.wait.Wait()
}

func (d *Datastore) LastKeys(batch int) []string {
	if batch > d.windowSize {
		log.Fatal().Err(errors.New("invalid batch size")).Int("batch", batch).
			Msg("please increase the window while creating datastore")
	}
	d.windowMu.RLock()
	defer d.windowMu.RUnlock()

	if d.window.Size() > batch {
		size := d.window.Size()
		return d.window.Values()[size-batch : size]
	}
	return d.window.Values()
}

func (d *Datastore) GetBatch(keys []string) ([]map[string]any, error) {
	rawValues := make([][]byte, len(keys))
	err := d.db.View(func(txn *badger.Txn) error {
		for i, key := range keys {
			itr, err := txn.Get([]byte(key))
			if err != nil {
				rawValues[i] = []byte(ErrorStr(err))
			} else {
				rawValues[i], err = itr.ValueCopy(nil)
				if err != nil {
					rawValues[i] = []byte(ErrorStr(err))
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	values := make([]map[string]any, len(keys))
	for i, buf := range rawValues {
		key := keys[i]
		if strings.HasPrefix(key, PrefixTcpRecord) {
			var metric TcpMetric
			err := proto.Unmarshal(buf, &metric)
			if err != nil {
				values[i] = ErrorJSON(err)
				break
			}
			values[i] = ToProtojson(&metric)
		} else if strings.HasPrefix(key, PrefixNicRecord) {
			var metric NicMetric
			err := proto.Unmarshal(buf, &metric)
			if err != nil {
				values[i] = ErrorJSON(err)
				break
			}
			values[i] = ToProtojson(&metric)
		} else {
			log.Fatal().Str("key", key).Msg("unknown key prefix")
		}
	}
	return values, nil
}

func (d *Datastore) writer(initialEpoch uint64, window queues.Queue[string]) {
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

		if len(req.Key) <= 0 || len(req.Value) <= 0 {
			log.Warn().Str("Key", req.Key).Int("ValueLen", len(req.Value)).
				Msg("ignore invalid KVPair")
			continue
		}

		key := fmt.Sprintf("%s%v", req.Key, epoch)
		req.Key = key
		epoch++

		// TODO(fanyang) batch write txn
		err := d.db.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(req.Key), req.Value)
		})
		log.Trace().Str("key", req.Key).Int("valueLen", len(req.Value)).Msg("write new item")

		if err != nil {
			log.Warn().Err(err).Msg("failed to insert db")
		}

		d.windowMu.Lock()
		if window.Size() >= d.windowSize {
			window.Dequeue()
		}
		window.Enqueue(key)
		d.windowMu.Unlock()
	}
}
