package tcpmon

import (
	"io"
	"math"

	"github.com/cockroachdb/errors"
	"github.com/dgraph-io/badger/v4"
)

// GetPrefix method returns all pairs in *reversed* order, the key starts with prefix
// maxCount indicates the maximum number of entries to be returned, <= 0 means unlimited
// hasValue indicates whether the returned value contains value
// prefix can be null, in this case, all key-value pairs are returned.
func (d *Datastore) GetPrefix(prefix []byte, maxCount int, hasValue bool) ([]KVPair, error) {
	if maxCount <= 0 {
		maxCount = math.MaxInt
	}

	r := make([]KVPair, 0)
	err := d.db.View(func(txn *badger.Txn) (err error) {
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

// GetKeys returns all keys in the database
func (d *Datastore) GetKeys() ([]string, error) {
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

func (d *Datastore) Backup(w io.Writer, since uint64) (uint64, error) {
	return d.db.Backup(w, since)
}
