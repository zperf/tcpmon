package cmd

import (
	"os"
	"path/filepath"
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

	// cmd: backup
	loadAndParseCmd.Flags().String("default-db", "/tmp/lap/defaultDB", "default db path for parse db data")
	fatalIf(viper.BindPFlag("default-db", loadAndParseCmd.Flags().Lookup("default-db")))
	loadAndParseCmd.Flags().String("load-db", "", "db path for recovering from backup, empty for not recovering")
	fatalIf(viper.BindPFlag("load-db", loadAndParseCmd.Flags().Lookup("load-db")))
	loadAndParseCmd.Flags().String("input", "/tmp/lap/input.txt", "input backup file")
	fatalIf(viper.BindPFlag("input", loadAndParseCmd.Flags().Lookup("input")))
	loadAndParseCmd.Flags().String("output", "", "output json format file, empty for stdout")
	fatalIf(viper.BindPFlag("output", loadAndParseCmd.Flags().Lookup("output")))
	loadAndParseCmd.Flags().String("prefix", "", "key prefix, empty for all key")
	fatalIf(viper.BindPFlag("prefix", loadAndParseCmd.Flags().Lookup("prefix")))
	rootCmd.AddCommand(loadAndParseCmd)

	// cmd: config
	configCmd.AddCommand(configGetDefaultCmd)
	rootCmd.AddCommand(configCmd)

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

	// garbage collection and key reclaiming
	rootCmd.PersistentFlags().Int("max-size", 10000, "badger db max size")
	fatalIf(viper.BindPFlag("max-size", rootCmd.PersistentFlags().Lookup("max-size")))
	rootCmd.PersistentFlags().Int("delete-size", 2000, "badger db delete size for once")
	fatalIf(viper.BindPFlag("delete-size", rootCmd.PersistentFlags().Lookup("delete-size")))

	rootCmd.PersistentFlags().Duration("reclaim-period", 10*time.Minute, "period of badger db delete")
	fatalIf(viper.BindPFlag("reclaim-period", rootCmd.PersistentFlags().Lookup("reclaim-period")))
	rootCmd.PersistentFlags().Duration("gc-period", 10*time.Minute, "period of badger db GC")
	fatalIf(viper.BindPFlag("gc-period", rootCmd.PersistentFlags().Lookup("gc-period")))

	// log
	rootCmd.PersistentFlags().String("log-dir", "/tmp/tcpmon-log", "dir to save log files")
	fatalIf(viper.BindPFlag("log-dir", rootCmd.PersistentFlags().Lookup("log-dir")))
	rootCmd.PersistentFlags().String("log-filename", "tcpmon.log", "filename of log files")
	fatalIf(viper.BindPFlag("log-filename", rootCmd.PersistentFlags().Lookup("log-filename")))
	rootCmd.PersistentFlags().Int("log-max-size", 10, "the max size in MB of the logfile before it's rolled")
	fatalIf(viper.BindPFlag("log-max-size", rootCmd.PersistentFlags().Lookup("log-max-size")))
	rootCmd.PersistentFlags().Int("log-max-backups", 5, "the max number of rolled files to keep")
	fatalIf(viper.BindPFlag("log-max-backups", rootCmd.PersistentFlags().Lookup("log-max-backups")))
	rootCmd.PersistentFlags().Int("log-max-age", 10, "the max age in days to keep a logfile")
	fatalIf(viper.BindPFlag("log-max-age", rootCmd.PersistentFlags().Lookup("log-max-age")))
	rootCmd.PersistentFlags().String("log-level", "Trace", "log level")
	fatalIf(viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level")))

	err := viper.ReadInConfig()
	if err != nil {
		var expected viper.ConfigFileNotFoundError
		if errors.As(err, &expected) {
			log.Warn().Err(err).Msg("config file not found, creating default config file")
			err = writeDefaultConfig()
			if err != nil {
				log.Warn().Err(err).Msg("create default config file failed")
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

func writeDefaultConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return errors.WithStack(err)
	}

	parentDir := filepath.Join(home, ".tcpmon")
	err = os.MkdirAll(parentDir, 0755)
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(viper.SafeWriteConfigAs(filepath.Join(parentDir, "config.yaml")))
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
