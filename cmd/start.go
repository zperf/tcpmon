package cmd

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zperf/tcpmon/tcpmon"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start monitoring",
	Run: func(cmd *cobra.Command, args []string) {
		level, _ := zerolog.ParseLevel(viper.GetString("log-level"))
		config := tcpmon.LogConfig{
			Level:                 level,
			ConsoleLoggingEnabled: viper.GetBool("verbose"),
			FileLoggingEnabled:    true,
			Directory:             viper.GetString("log-dir"),
			Filename:              viper.GetString("log-filename"),
			MaxSize:               viper.GetInt("log-max-size"),
			MaxBackups:            viper.GetInt("log-max-backups"),
			MaxAge:                viper.GetInt("log-max-age"),
		}
		tcpmon.InitLogger(&config)
		if level == zerolog.NoLevel {
			log.Warn().Str("level", viper.GetString("log-level")).
				Msg("invalid level, default to NoLevel")
		}
		log.Info().Str("configFile", viper.ConfigFileUsed()).Str("logDir", viper.GetString("log-dir")).Msg("Config loaded")

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

func init() {
	rootCmd.AddCommand(startCmd)

	flag := "listen"
	startCmd.Flags().StringP(flag, "l", "0.0.0.0:6789", "HTTP server listening at this address")
	fatalIf(viper.BindPFlag(flag, startCmd.Flags().Lookup(flag)))
}
