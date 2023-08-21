package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zperf/tcpmon/tcpmon"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "restore backup file",
	Run: func(cmd *cobra.Command, args []string) {
		// open input file
		inputFile, err := os.Open(viper.GetString("input"))
		if err != nil {
			log.Err(err).Msg("cannot open input file")
			return
		}
		defer inputFile.Close()

		// load backup to load-db
		loadDBPath := viper.GetString("load-db")
		if loadDBPath != "" {
			opts := badger.DefaultOptions(loadDBPath)
			db, err := badger.Open(opts)
			if err != nil {
				log.Err(err).Msg("cannot open load-db")
				return
			}
			defer db.Close()
			err = db.Load(inputFile, 32)
			if err != nil {
				log.Err(err).Msg("error load backup to load-db")
				return
			}
			return
		}

		// open default-db
		opts := badger.DefaultOptions(viper.GetString("default-db"))
		db, err := badger.Open(opts)
		if err != nil {
			log.Err(err).Msg("cannot open default-db")
			return
		}
		defer db.Close()
		// clear default-db
		err = db.DropAll()
		if err != nil {
			log.Err(err).Msg("cannot clear default-db")
			return
		}
		// load backup to default-db
		err = db.Load(inputFile, 32)
		if err != nil {
			log.Err(err).Msg("error load backup to default-db")
			return
		}

		// open output file
		outputFilePath := viper.GetString("output")
		var encoder *json.Encoder
		if outputFilePath != "" {
			OutputFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				log.Err(err).Msg("cannot open output file")
				return
			}
			defer OutputFile.Close()
			encoder = json.NewEncoder(OutputFile)
		}

		// parse default-db
		err = db.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			prefix := viper.GetString("prefix")
			if prefix != "" {
				opts.Prefix = []byte(prefix)
			}
			it := txn.NewIterator(opts)
			defer it.Close()

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
				if encoder != nil {
					err = encoder.Encode(msg)
					if err != nil {
						log.Warn().Err(err).Str("Key", key).Msg("fail to write data to output file")
						continue
					}
				} else {
					fmt.Printf("{%v}\n\n", msg)
				}
			}
			return nil
		})
		if err != nil {
			log.Err(err).Msg("error parse default-db")
		}
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().String("default-db", "/tmp/lap/defaultDB", "default db path for parse db data")
	fatalIf(viper.BindPFlag("default-db", restoreCmd.Flags().Lookup("default-db")))
	restoreCmd.Flags().String("load-db", "", "db path for recovering from backup, empty for not recovering")
	fatalIf(viper.BindPFlag("load-db", restoreCmd.Flags().Lookup("load-db")))
	restoreCmd.Flags().String("input", "/tmp/lap/input.txt", "input backup file")
	fatalIf(viper.BindPFlag("input", restoreCmd.Flags().Lookup("input")))
	restoreCmd.Flags().String("output", "", "output json format file, empty for stdout")
	fatalIf(viper.BindPFlag("output", restoreCmd.Flags().Lookup("output")))
	restoreCmd.Flags().String("prefix", "", "key prefix, empty for all key")
	fatalIf(viper.BindPFlag("prefix", restoreCmd.Flags().Lookup("prefix")))
}
