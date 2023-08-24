package tcpmon

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/dgraph-io/badger/v4"
)

// GetPrefix method returns all pairs in *reversed* order, the key starts with prefix
// maxCount indicates the maximum number of entries to be returned, <= 0 means unlimited
// hasValue indicates whether the returned value contains value
// prefix can be null, in this case, all key-value pairs are returned.
func GetPrefix(db *badger.DB, prefix []byte, maxCount int, hasValue bool) ([]KVPair, error) {
	if maxCount <= 0 {
		maxCount = math.MaxInt
	}

	r := make([]KVPair, 0)
	err := db.View(func(txn *badger.Txn) (err error) {
		err = nil
		options := badger.DefaultIteratorOptions
		options.Reverse = true
		if len(prefix) != 0 {
			options.Prefix = prefix
		}
		itr := txn.NewIterator(options)
		defer itr.Close()

		count := 0
		// append 0xff is a trick for reverse iteration
		// see more: https://github.com/dgraph-io/badger/issues/436#issuecomment-400095559
		for itr.Seek(append(prefix, 0xff)); itr.Valid() && count < maxCount; itr.Next() {
			count++
			item := itr.Item()
			key := item.Key()
			var value []byte
			if hasValue {
				value, err = item.ValueCopy(nil)
				if err != nil {
					return errors.WithStack(err)
				}
			} else {
				value = nil
			}
			r = append(r, KVPair{
				Key:   string(key),
				Value: value,
			})
		}
		return
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return r, nil
}

func (d *DataStore) GetPrefix(prefix []byte, maxCount int, hasValue bool) ([]KVPair, error) {
	return GetPrefix(d.db, prefix, maxCount, hasValue)
}

// GetKeys returns all keys in the database
func (d *DataStore) GetKeys() ([]string, error) {
	keys := make([]string, 0)
	err := d.db.View(func(txn *badger.Txn) error {
		options := badger.DefaultIteratorOptions
		itr := txn.NewIterator(options)
		defer itr.Close()

		for itr.Rewind(); itr.Valid(); itr.Next() {
			item := itr.Item()
			key := item.Key()
			keys = append(keys, string(key))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (d *DataStore) GetMetrics() ([]string, error) {
	keys, err := d.GetKeys()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	metrics := make([]string, 0)
	for _, k := range keys {
		if !strings.HasPrefix(k, PrefixMember) && !strings.HasPrefix(k, PrefixMetadata) {
			metrics = append(metrics, k)
		}
	}

	return metrics, nil
}

func (d *DataStore) Backup(w io.Writer, since uint64) (uint64, error) {
	return d.db.Backup(w, since)
}

func (d *DataStore) Get(key string) (*KVPair, error) {
	p := &KVPair{
		Key: key,
	}
	err := d.db.View(func(txn *badger.Txn) error {
		it, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		p.Value, err = it.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return p, nil
}

func (d *DataStore) GetKeyCount(kind string) (uint32, error) {
	if !ValidCountKind(kind) {
		return 0, errors.Newf("invalid type: '%s'", kind)
	}

	val := uint32(0)
	key := fmt.Sprintf("%s/count/%s", PrefixMetadata, kind)

	err := d.db.View(func(txn *badger.Txn) error {
		it, err := txn.Get([]byte(key))
		if err != nil {
			return errors.WithStack(err)
		}
		buf, err := it.ValueCopy(nil)
		if err != nil {
			return errors.WithStack(err)
		}

		v, err := strconv.ParseUint(string(buf), 10, 32)
		if err != nil {
			return errors.WithStack(err)
		}

		val = uint32(v)
		return nil
	})
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (d *DataStore) GetTotalCount() (uint32, error) {
	return d.GetKeyCount("total")
}

func (d *DataStore) GetTcpKeyCount() (uint32, error) {
	return d.GetKeyCount("tcp")
}

func (d *DataStore) GetNicKeyCount() (uint32, error) {
	return d.GetKeyCount("nic")
}

func (d *DataStore) GetNetKeyCount() (uint32, error) {
	return d.GetKeyCount("net")
}
