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

func Execute() {
	// init viper
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/tcpmon/")
	viper.AddConfigPath("$HOME/.tcpmon")
	viper.AddConfigPath(".")

	// init global flags
	initCommandFlags()
	initReclaimFlags()
	initLoggingFlags()

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

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initCommandFlags() {
	rootCmd.PersistentFlags().String("db", "/tmp/tcpmon/db", "Database path")
	fatalIf(viper.BindPFlag("db", rootCmd.PersistentFlags().Lookup("db")))

	rootCmd.PersistentFlags().String("ifconfig", "/usr/bin/ifconfig", "Command 'ifconfig' path")
	fatalIf(viper.BindPFlag("ifconfig", rootCmd.PersistentFlags().Lookup("ifconfig")))

	rootCmd.PersistentFlags().String("ss", "/usr/bin/ss", "Command 'ss' path")
	fatalIf(viper.BindPFlag("ss", rootCmd.PersistentFlags().Lookup("ss")))

	rootCmd.PersistentFlags().String("ss-arg", "-4ntipmoHna", "Set the arg for ss")
	fatalIf(viper.BindPFlag("ss-arg", rootCmd.PersistentFlags().Lookup("ss-arg")))

	rootCmd.PersistentFlags().String("netstat", "/usr/bin/netstat", "Command 'netstat' path")
	fatalIf(viper.BindPFlag("netstat", rootCmd.PersistentFlags().Lookup("netstat")))

	rootCmd.PersistentFlags().String("netstat-arg", "-s", "Set the arg for netstat")
	fatalIf(viper.BindPFlag("netstat-arg", rootCmd.PersistentFlags().Lookup("netstat-arg")))

	rootCmd.PersistentFlags().DurationP("command-timeout", "c", time.Second,
		"Command execution timeout")
	fatalIf(viper.BindPFlag("command-timeout", rootCmd.PersistentFlags().Lookup("command-timeout")))

	rootCmd.PersistentFlags().DurationP("interval", "i", time.Second,
		"Interval between two metric collections")
	fatalIf(viper.BindPFlag("interval", rootCmd.PersistentFlags().Lookup("interval")))
}

func initReclaimFlags() {
	rootCmd.PersistentFlags().Int("max-size", 10000, "badger db max size")
	fatalIf(viper.BindPFlag("max-size", rootCmd.PersistentFlags().Lookup("max-size")))

	rootCmd.PersistentFlags().Int("delete-size", 2000, "badger db delete size for once")
	fatalIf(viper.BindPFlag("delete-size", rootCmd.PersistentFlags().Lookup("delete-size")))

	rootCmd.PersistentFlags().Duration("reclaim-period", 10*time.Minute, "period of badger db delete")
	fatalIf(viper.BindPFlag("reclaim-period", rootCmd.PersistentFlags().Lookup("reclaim-period")))

	rootCmd.PersistentFlags().Duration("gc-period", 10*time.Minute, "period of badger db GC")
	fatalIf(viper.BindPFlag("gc-period", rootCmd.PersistentFlags().Lookup("gc-period")))
}

func initLoggingFlags() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode")
	fatalIf(viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")))

	rootCmd.PersistentFlags().String("log-dir", "/tmp/tcpmon/log", "dir to save log files")
	fatalIf(viper.BindPFlag("log-dir", rootCmd.PersistentFlags().Lookup("log-dir")))

	rootCmd.PersistentFlags().String("log-filename", "tcpmon.log", "filename of log files")
	fatalIf(viper.BindPFlag("log-filename", rootCmd.PersistentFlags().Lookup("log-filename")))

	rootCmd.PersistentFlags().Int("log-max-size", 10, "the max size in MB of the logfile before it's rolled")
	fatalIf(viper.BindPFlag("log-max-size", rootCmd.PersistentFlags().Lookup("log-max-size")))

	rootCmd.PersistentFlags().Int("log-max-backups", 5, "the max number of rolled files to keep")
	fatalIf(viper.BindPFlag("log-max-backups", rootCmd.PersistentFlags().Lookup("log-max-backups")))

	rootCmd.PersistentFlags().Int("log-max-age", 10, "the max age in days to keep a logfile")
	fatalIf(viper.BindPFlag("log-max-age", rootCmd.PersistentFlags().Lookup("log-max-age")))

	rootCmd.PersistentFlags().String("log-level", "info", "log level")
	fatalIf(viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level")))
}
