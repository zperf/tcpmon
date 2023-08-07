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
		config := tcpmon.Config{
			ConsoleLoggingEnabled: true,
			FileLoggingEnabled:    true,
			Directory:             viper.GetString("log-dir"),
			Filename:              viper.GetString("log-filename"),
			MaxSize:               viper.GetInt("log-max-size"),
			MaxBackups:            viper.GetInt("log-max-backups"),
			MaxAge:                viper.GetInt("log-max-age"),
		}
		tcpmon.InitLogger(config)

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
