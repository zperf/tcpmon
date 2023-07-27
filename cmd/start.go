package cmd

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zperf/tcpmon/tcpmon"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start monitoring",
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}).
			Level(zerolog.InfoLevel).
			With().
			Timestamp().
			Caller().
			Logger()

		m := tcpmon.New()
		err := m.Run()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to run")
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
