package cmd

import (
	"context"
	"fmt"
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
		m, err := tcpmon.New(tcpmon.MonitorConfig{
			CollectInterval: viper.GetDuration("collect-interval"),
			HttpListen:      viper.GetString("listen"),
			QuorumPort:      viper.GetInt("quorum-port"),
			DataStoreConfig: tcpmon.DataStoreConfig{
				Path:            viper.GetString("db"),
				MaxSize:         viper.GetInt("db-max-size"),
				GcInterval:      viper.GetDuration("db-gc-interval"),
				ReclaimBatch:    viper.GetInt("reclaim-batch"),
				ReclaimInterval: viper.GetDuration("reclaim-interval"),
			},
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Create tcpmon failed")
		}

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		err = m.Run(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to run")
		}
	},
}

var startTestCmd = &cobra.Command{
	Use:     "test",
	Short:   "Test this machines can run daemon",
	Example: `  tcpmon-linux start test > test.sh; bash test.sh`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("#!/usr/bin/env bash")
		fmt.Println("set -x")
		fmt.Printf("%s %s\n", viper.GetString("cmd-ss"), viper.GetString("cmd-ss-arg"))
		fmt.Printf("%s %s\n", viper.GetString("cmd-ss2"), viper.GetString("cmd-ss-arg"))
		fmt.Printf("%s %s\n", viper.GetString("cmd-ifconfig"), viper.GetString("cmd-ifconfig-arg"))
		fmt.Printf("%s %s\n", viper.GetString("cmd-netstat"), viper.GetString("cmd-netstat-arg"))
		fmt.Println("echo \"$?\"")
	},
}

func init() {
	// monitor flags
	startCmd.PersistentFlags().DurationP("collect-interval", "i", time.Second, "Metric collection interval")
	startCmd.PersistentFlags().StringP("listen", "l", "0.0.0.0:6789", "HTTP server listening at this address")
	startCmd.PersistentFlags().IntP("quorum-port", "q", 6790, "Quorum bind and advertised port")

	// monitor command flags
	startCmd.PersistentFlags().String("cmd-ifconfig", "/usr/bin/ifconfig", "The path of 'ifconfig'")
	startCmd.PersistentFlags().String("cmd-ifconfig-arg", "", "Parameters when executing 'ifconfig'")
	startCmd.PersistentFlags().String("cmd-ss", "/usr/bin/ss", "The path of 'ss'")
	startCmd.PersistentFlags().String("cmd-ss2", "/usr/sbin/ss", "The path of 'ss'")
	startCmd.PersistentFlags().String("cmd-ss-arg", "-4ntipmoHna", "Parameters when executing 'ss'")
	startCmd.PersistentFlags().String("cmd-netstat", "/usr/bin/netstat", "The path of 'netstat'")
	startCmd.PersistentFlags().String("cmd-netstat-arg", "-s", "Parameters when executing 'netstat'")
	startCmd.PersistentFlags().DurationP("cmd-timeout", "c", time.Second, "Command execution timeout")

	// db flags
	startCmd.PersistentFlags().String("db", "/tmp/tcpmon/db", "Database path")
	startCmd.PersistentFlags().Int("db-max-size", 10000, "Maximum number of records in the database")
	startCmd.PersistentFlags().Duration("db-gc-interval", 10*time.Minute, "BadgerDB value GC interval")
	startCmd.PersistentFlags().Int("reclaim-batch", 2000, "Maximum number of reclaiming per batch")
	startCmd.PersistentFlags().Duration("reclaim-interval", 3*time.Minute, "Reclaiming interval")

	fatalIf(viper.BindPFlags(startCmd.PersistentFlags()))
	rootCmd.AddCommand(startCmd)
	startCmd.AddCommand(startTestCmd)
}
