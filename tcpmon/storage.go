package tcpmon

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

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
	stop int32
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
		stop:       0,
	}
	d.wait.Add(1)
	go d.writer(epoch, d.window)
	go d.periodicDelete(10000, 100)
	return d
}

// Tx returns a send-only channel
func (d *Datastore) Tx() chan<- *KVPair {
	return d.tx
}

func (d *Datastore) Close() {
	close(d.done)
	atomic.StoreInt32(&d.stop, 1)
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

func (d *Datastore) getSize() int {
	txn := d.db.NewTransaction(false)
	defer txn.Discard()

	opts := badger.DefaultIteratorOptions
	it := txn.NewIterator(opts)
	defer it.Close()

	totalCount := 0
	for it.Rewind(); it.Valid(); it.Next() {
		totalCount++
	}
	return totalCount
}

func (d *Datastore) periodicDelete(maxSize, deleteSize int) {
	for atomic.LoadInt32(&d.stop) == 0 {
		totalCount := d.getSize()
		if totalCount > maxSize {
			txn := d.db.NewTransaction(true)
			opts := badger.DefaultIteratorOptions
			it := txn.NewIterator(opts)
			k := 0
			for it.Rewind(); it.Valid(); it.Next() {
				if k >= deleteSize {
					break
				}
				item := it.Item()
				key := item.Key()
				err := txn.Delete(key)
				if err != nil {
					log.Warn().Err(err).Str("key", string(key)).Msg("failed to delete item")
				} else {
					log.Trace().Str("key", string(key)).Msg("delete old item")
				}
				k++
			}
			it.Close()
			err := txn.Commit()
			if err != nil {
				log.Warn().Err(err).Msg("failed to commit")
			}
		}
		time.Sleep(1 * time.Second)
	}
}
