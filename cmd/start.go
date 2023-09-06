package cmd

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zperf/tcpmon/tcpmon"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start monitoring",
	Run: func(cmd *cobra.Command, args []string) {
		// create and start monitor
		m, err := tcpmon.New(tcpmon.MonitorConfig{
			CollectInterval: viper.GetDuration("collect-interval"),
			HttpListen:      viper.GetString("listen"),
			QuorumPort:      viper.GetInt("quorum-port"),
			DataStoreConfig: tcpmon.DataStoreConfig{
				Path:            viper.GetString("db"),
				MaxSize:         viper.GetUint32("db-max-size"),
				WriteInterval:   viper.GetDuration("db-write-interval"),
				ExpectedRss:     viper.GetUint64("db-expected-rss"),
				MinOpenInterval: viper.GetDuration("db-min-open-interval"),
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
		fmt.Printf("%s %s\n", viper.GetString("cmd-ifconfig2"), viper.GetString("cmd-ifconfig-arg"))
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
	startCmd.PersistentFlags().String("cmd-ifconfig2", "/usr/sbin/ifconfig", "The path of 'ifconfig'")
	startCmd.PersistentFlags().String("cmd-ifconfig-arg", "", "Parameters when executing 'ifconfig'")
	startCmd.PersistentFlags().String("cmd-ss", "/usr/bin/ss", "The path of 'ss'")
	startCmd.PersistentFlags().String("cmd-ss2", "/usr/sbin/ss", "The path of 'ss'")
	startCmd.PersistentFlags().String("cmd-ss-arg", "-4ntipmona", "Parameters when executing 'ss'")
	startCmd.PersistentFlags().String("cmd-netstat", "/usr/bin/netstat", "The path of 'netstat'")
	startCmd.PersistentFlags().String("cmd-netstat-arg", "-s", "Parameters when executing 'netstat'")
	startCmd.PersistentFlags().DurationP("cmd-timeout", "c", time.Second, "Command execution timeout")

	// db flags
	startCmd.PersistentFlags().String("db", "/tmp/tcpmon/db", "Database path")
	startCmd.PersistentFlags().Uint32("db-max-size", 2000000, "Maximum number of records in the database")
	startCmd.PersistentFlags().Duration("db-write-interval", 3*time.Second, "Write interval")
	startCmd.PersistentFlags().Uint64("db-expected-rss", 200<<20, "Expected maximum RSS")
	startCmd.PersistentFlags().Duration("db-min-open-interval", 1*time.Minute, "Database reopen interval")

	fatalIf(viper.BindPFlags(startCmd.PersistentFlags()))
	rootCmd.AddCommand(startCmd)
	startCmd.AddCommand(startTestCmd)
}
