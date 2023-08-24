package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
	"time"

	"github.com/dgraph-io/badger/v4"
	boptions "github.com/dgraph-io/badger/v4/options"
)

func main1() {
	src, err := badger.Open(badger.DefaultOptions("/Users/fanyang/Downloads/tcpmon-data-17.95"))
	if err != nil {
		log.Fatalf("failed to open src")
	}
	defer src.Close()

	options := badger.DefaultOptions("/Users/fanyang/Downloads/tcpmon-data-17.95-out").
		WithInMemory(false).
		WithCompression(boptions.ZSTD).
		WithZSTDCompressionLevel(2).
		WithMaxLevels(1).
		WithBaseTableSize(200 << 20).
		WithValueLogFileSize(1 << 20)

	dst, err := badger.Open(options)
	if err != nil {
		log.Fatalf("failed to open src")
	}
	defer dst.Close()

	err = src.View(func(txn *badger.Txn) error {
		itr := txn.NewIterator(badger.DefaultIteratorOptions)
		defer itr.Close()

		for itr.Rewind(); itr.Valid(); itr.Next() {
			it := itr.Item()
			val, err := it.ValueCopy(nil)
			if err != nil {
				return err
			}

			err = dst.Update(func(txn *badger.Txn) error {
				return txn.Set(it.Key(), val)
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("failed to iterate")
	}

	err = dst.RunValueLogGC(0.5)
	if err != nil {
		log.Printf("run value log GC failed, %v", err)
	}
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6789", nil))
	}()

	src, err := badger.Open(badger.DefaultOptions("/Users/fanyang/Downloads/tcpmon-data-17.95"))
	if err != nil {
		log.Fatalf("failed to open src")
	}

	count := 0
	err = src.View(func(txn *badger.Txn) error {
		itr := txn.NewIterator(badger.DefaultIteratorOptions)
		defer itr.Close()
		for itr.Rewind(); itr.Valid(); itr.Next() {
			count++
		}
		return nil
	})
	if err != nil {
		log.Fatalf("failed to iterate")
	}
	log.Println("ok")

	src.Close()
	debug.FreeOSMemory()

	for {
		time.Sleep(time.Second)
	}
}
