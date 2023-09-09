package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1 "github.com/zperf/tcpmon/tcpmon/storage/v1"
)

var dbListCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List keys with prefix in the database",
	Example: "  list --db ./database --prefix tcp",
	Run: func(cmd *cobra.Command, args []string) {
		path := viper.GetString("db")
		prefix := viper.GetString("prefix")

		db := openBadgerForRead(path)
		defer db.Close()

		pairs, err := v1.GetPrefix(db, []byte(prefix), 0, false)
		if err != nil {
			log.Fatal().Err(err).Msg("Get pairs by prefix failed")
		}

		for _, p := range pairs {
			fmt.Println(p.Key)
		}
	},
}

func init() {
	dbCmd.AddCommand(dbListCmd)
}
