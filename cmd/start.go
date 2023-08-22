package cmd

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/cockroachdb/errors"
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
		// read config file
		err := viper.ReadInConfig()
		if err != nil {
			// config file not found, is must be a dev env
			// write a default config file to $HOME/.tcpmon/config.yaml
			var expected viper.ConfigFileNotFoundError
			if errors.As(err, &expected) {
				log.Warn().Err(err).Msg("config file not found, creating default config file")
				err = writeDefaultConfig()
				if err != nil {
					log.Fatal().Err(err).Msg("create default config file failed")
				}
			} else {
				log.Fatal().Err(err).Msg("failed to read config file")
			}
		}

		// init logger
		level, _ := zerolog.ParseLevel(viper.GetString("log-level"))
		config := tcpmon.LogConfig{
			Level:                 level,
			ConsoleLoggingEnabled: viper.GetBool("verbose"),
			FileLoggingEnabled:    true,
			Directory:             viper.GetString("log-dir"),
			Filename:              viper.GetString("log-filename"),
			MaxSize:               viper.GetInt("log-max-size"),
			MaxBackups:            viper.GetInt("log-max-count"),
		}
		tcpmon.InitLogger(&config)

		// print warnings after logger initialized
		if level == zerolog.NoLevel {
			log.Warn().Str("level", viper.GetString("log-level")).
				Msg("invalid level, default to NoLevel")
		}
		log.Info().Str("configFile", viper.ConfigFileUsed()).
			Str("logDir", viper.GetString("log-dir")).
			Msg("Config loaded")

		// create and start monitor
		m, err := tcpmon.New()
		if err != nil {
			log.Fatal().Err(err).Msg("Create tcpmon failed")
		}

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		err = m.Run(ctx, viper.GetDuration("collect-interval"), viper.GetString("listen"))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to run")
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// monitor flags
	startCmd.PersistentFlags().DurationP("collect-interval", "i", time.Second, "Metric collection interval")
	startCmd.PersistentFlags().StringP("listen", "l", "0.0.0.0:6789", "HTTP server listening at this address")

	// logging flags
	startCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode")
	startCmd.PersistentFlags().String("log-level", "info", "log level")
	startCmd.PersistentFlags().String("log-dir", "/tmp/tcpmon/log", "The log directory")
	startCmd.PersistentFlags().String("log-filename", "tcpmon.log", "The file name of logs")
	startCmd.PersistentFlags().Int("log-max-size", 10, "Maximum size of each log file")
	startCmd.PersistentFlags().Int("log-max-count", 5, "Maximum log files to keep")

	fatalIf(viper.BindPFlags(startCmd.PersistentFlags()))
}
