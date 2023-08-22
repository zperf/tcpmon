package cmd

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zperf/tcpmon/tcpmon"
)

var dbGetCmd = &cobra.Command{
	Use:     "get [KEY]",
	Short:   "Get record by key",
	Example: "  get --db ./db tcp/1692679337381",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := viper.GetString("db")
		key := []byte(args[0])

		db := openBadgerForRead(path)
		defer db.Close()

		err := db.View(func(txn *badger.Txn) error {
			it, err := txn.Get(key)
			if err != nil {
				log.Fatal().Err(err).Send()
			}

			val, err := it.ValueCopy(nil)
			if err != nil {
				log.Fatal().Err(err).Send()
			}

			fmt.Println(tcpmon.NewKVPair(string(key), val).ToJSONString())
			return nil
		})
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	},
}

func init() {
	dbCmd.AddCommand(dbGetCmd)
}
