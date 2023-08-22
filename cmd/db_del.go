package cmd

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dbDelCmd = &cobra.Command{
	Use:     "del [KEY]",
	Short:   "Delete record by key",
	Example: "  del --db ./db tcp/1692679337381",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		path := viper.GetString("db")

		db := openBadger(path)
		defer db.Close()

		err := db.Update(func(txn *badger.Txn) error {
			return txn.Delete([]byte(key))
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Delete failed")
		}
	},
}

func init() {
	dbCmd.AddCommand(dbDelCmd)
}
