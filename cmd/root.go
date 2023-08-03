package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "tcpmon",
	Short: "Tcpmon is a portable local netowrk monitor for Linux",
}

func fatalIf(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal error")
	}
}

func Execute() {
	initViper()

	rootCmd.Flags().Bool("verbose", false, "Verbose mode")

	// cmd: start
	startCmd.Flags().BoolP("foreground", "f", false, "Run in foreground")
	fatalIf(viper.BindPFlag("foreground", startCmd.Flags().Lookup("foreground")))
	startCmd.Flags().String("listen", "0.0.0.0:6789", "HTTP server listening at this address")
	fatalIf(viper.BindPFlag("listen", startCmd.Flags().Lookup("listen")))
	rootCmd.AddCommand(startCmd)

	// path
	rootCmd.PersistentFlags().String("db", "/tmp/tcpmon", "Database path")
	fatalIf(viper.BindPFlag("db", rootCmd.PersistentFlags().Lookup("db")))
	rootCmd.PersistentFlags().String("ifconfig", "/usr/bin/ifconfig", "Command 'ifconfig' path")
	fatalIf(viper.BindPFlag("ifconfig", rootCmd.PersistentFlags().Lookup("ifconfig")))
	rootCmd.PersistentFlags().String("ss", "/usr/bin/ss", "Command 'ss' path")
	fatalIf(viper.BindPFlag("ss", rootCmd.PersistentFlags().Lookup("ss")))
	rootCmd.PersistentFlags().String("ss-arg", "-4ntipmoHOna", "Set the arg for ss")
	fatalIf(viper.BindPFlag("ss-arg", rootCmd.PersistentFlags().Lookup("ss-arg")))
	rootCmd.PersistentFlags().String("netstat", "/usr/bin/netstat", "Command 'netstat' path")
	fatalIf(viper.BindPFlag("netstat", rootCmd.PersistentFlags().Lookup("netstat")))
	rootCmd.PersistentFlags().String("netstat-arg", "-s", "Set the arg for netstat")
	fatalIf(viper.BindPFlag("netstat-arg", rootCmd.PersistentFlags().Lookup("netstat-arg")))

	// timeout
	rootCmd.PersistentFlags().DurationP("command-timeout", "c", time.Second,
		"Command execution timeout")
	fatalIf(viper.BindPFlag("command-timeout", rootCmd.PersistentFlags().Lookup("command-timeout")))

	// interval
	rootCmd.PersistentFlags().DurationP("interval", "i", time.Second,
		"Interval between two metric collections")
	fatalIf(viper.BindPFlag("interval", rootCmd.PersistentFlags().Lookup("interval")))

	// size
	rootCmd.PersistentFlags().Int("max-size", 10000, "badger db max size")
	fatalIf(viper.BindPFlag("max-size", rootCmd.PersistentFlags().Lookup("max-size")))
	rootCmd.PersistentFlags().Int("delete-size", 100, "badger db delete size for once")
	fatalIf(viper.BindPFlag("delete-size", rootCmd.PersistentFlags().Lookup("delete-size")))

	// period
	rootCmd.PersistentFlags().Int("period-second", 1, "period of badger db delete")
	fatalIf(viper.BindPFlag("period-second", rootCmd.PersistentFlags().Lookup("period-second")))
	rootCmd.PersistentFlags().Int("period-minute", 5, "period of badger db GC")
	fatalIf(viper.BindPFlag("period-minute", rootCmd.PersistentFlags().Lookup("period-minute")))

	err := viper.ReadInConfig()
	if err != nil {
		var expected viper.ConfigFileNotFoundError
		if errors.As(err, &expected) {
			log.Warn().Err(err).Msg("config file not found")
			// TODO(fanyang) create default config file
		} else {
			log.Fatal().Err(err).Msg("failed to read config file")
		}
	}

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initViper() {
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/tcpmon/")
	viper.AddConfigPath("$HOME/.tcpmon")
	viper.AddConfigPath(".")
}
