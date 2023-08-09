package cmd

import (
	"encoding/json"
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zperf/tcpmon/tcpmon"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "get db backup json data",
	Run: func(cmd *cobra.Command, args []string) {
		// open backup file
		backupFile, err := os.Open(viper.GetString("backup-file"))
		if err != nil {
			log.Err(err).Msg("cannot open backup file")
			return
		}
		defer backupFile.Close()

		// open data file
		dataFile, err := os.OpenFile(viper.GetString("data-file"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Err(err).Msg("cannot open data file")
			return
		}
		defer dataFile.Close()

		// open db
		opts := badger.DefaultOptions(viper.GetString("backup-db"))
		db, err := badger.Open(opts)
		if err != nil {
			log.Err(err).Msg("cannot open db")
			return
		}
		defer db.Close()

		// clear db
		err = db.DropAll()
		if err != nil {
			log.Err(err).Msg("cannot clear db")
			return
		}

		// load backup file to db
		err = db.Load(backupFile, 32)
		if err != nil {
			log.Err(err).Msg("error load backup file to db")
			return
		}

		// write data to data file from db
		err = db.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			it := txn.NewIterator(opts)
			defer it.Close()
			encoder := json.NewEncoder(dataFile)
			for it.Rewind(); it.Valid(); it.Next() {
				key := string(it.Item().KeyCopy(nil))
				value, err := it.Item().ValueCopy(nil)
				if err != nil {
					log.Warn().Err(err).Str("Key", key).Msg("fail to get value")
					continue
				}
				kvp := tcpmon.KVPair{
					Key:   key,
					Value: value,
				}
				msg, _ := kvp.ToProto()
				err = encoder.Encode(msg)
				if err != nil {
					log.Warn().Err(err).Str("Key", key).Msg("fail to write data")
					continue
				}
			}
			return nil
		})
		if err != nil {
			log.Err(err).Msg("write data to data file from db failed")
		}
	},
}
