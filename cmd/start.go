package cmd

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zperf/tcpmon/tcpmon"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start monitoring",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := tcpmon.New()
		if err != nil {
			log.Fatal().Err(err).Msg("Create tcpmon failed")
		}

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		err = m.Run(ctx, viper.GetDuration("interval"))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to run")
		}
	},
}
