package cmd

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	v1 "github.com/zperf/tcpmon/tcpmon/storage/v1"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database management",
}

func init() {
	rootCmd.AddCommand(dbCmd)
}

func openBadgerForRead(path string) *badger.DB {
	db, err := badger.Open(badger.DefaultOptions(path).
		WithLogger(v1.NewBadgerLogger()).
		WithReadOnly(true))
	if err != nil {
		log.Fatal().Err(err).Msg("Open db failed")
	}
	return db
}

func openBadger(path string) *badger.DB {
	db, err := badger.Open(badger.DefaultOptions(path).
		WithLogger(v1.NewBadgerLogger()).
		WithCompactL0OnClose(true))
	if err != nil {
		log.Fatal().Err(err).Msg("Open db failed")
	}
	return db
}
