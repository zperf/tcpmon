package cmd

import (
	"github.com/cockroachdb/errors"
	"github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zperf/tcpmon/tcpmon"
)

var dbPrintCmd = &cobra.Command{
	Use:   "print",
	Short: "Print database as JSON",
	Run: func(cmd *cobra.Command, args []string) {
		path := viper.GetString("db")
		prefix := viper.GetString("prefix")
		reversed := viper.GetBool("reversed")

		db, err := badger.Open(badger.DefaultOptions(path).
			WithLogger(&tcpmon.BadgerDbLogger{}).
			WithReadOnly(true))
		if err != nil {
			log.Fatal().Err(err).Msg("Open db failed")
		}
		defer db.Close()

		DoPrint(db, prefix, reversed, nil)
	},
}

func DoPrint(db *badger.DB, prefix string, reversed bool, printFn func(p tcpmon.KVPair)) {
	err := db.View(func(txn *badger.Txn) error {
		options := badger.DefaultIteratorOptions
		if len(prefix) != 0 {
			options.Prefix = []byte(prefix)
		}
		options.Reverse = reversed

		itr := txn.NewIterator(options)
		defer itr.Close()

		if prefix == "" {
			itr.Seek(nil)
		} else if reversed {
			itr.Seek(append([]byte(prefix), 0xff))
		} else {
			itr.Seek([]byte(prefix))
		}

		for ; itr.Valid(); itr.Next() {
			value, err := itr.Item().ValueCopy(nil)
			if err != nil {
				return errors.WithStack(err)
			}
			key := itr.Item().KeyCopy(nil)
			pair := tcpmon.KVPair{
				Key:   string(key),
				Value: value,
			}

			if printFn == nil {
				log.Info().Msg(pair.ToJSONString())
			} else {
				printFn(pair)
			}
		}

		return nil
	})
	if err != nil {
		log.Info().Err(err).Msg("")
	}
}

func addPrintFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("prefix", "p", "", "Prefix filter")
	cmd.Flags().BoolP("reversed", "r", false, "Reversed iterate order")
}

func init() {
	dbCmd.AddCommand(dbPrintCmd)
	addPrintFlags(dbPrintCmd)
}
