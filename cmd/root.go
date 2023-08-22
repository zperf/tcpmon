package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "tcpmon",
	Short: "Tcpmon is a portable local network monitor for Linux",
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
	initDbFlags()
	initMonitorFlags()
	fatalIf(viper.BindPFlags(rootCmd.PersistentFlags()))

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initMonitorFlags() {
	rootCmd.PersistentFlags().String("cmd-ifconfig", "/usr/bin/ifconfig", "The path of 'ifconfig'")
	rootCmd.PersistentFlags().String("cmd-ss", "/usr/bin/ss", "The path of 'ss'")
	rootCmd.PersistentFlags().String("cmd-ss-arg", "-4ntipmoHna", "Parameters when executing 'ss'")
	rootCmd.PersistentFlags().String("cmd-netstat", "/usr/bin/netstat", "The path of 'netstat'")
	rootCmd.PersistentFlags().String("cmd-netstat-arg", "-s", "Parameters when executing 'netstat'")
	rootCmd.PersistentFlags().DurationP("cmd-timeout", "c", time.Second, "Command execution timeout")
}

func initDbFlags() {
	rootCmd.PersistentFlags().String("db", "/tmp/tcpmon/db", "Database path")
	rootCmd.PersistentFlags().Int("db-max-size", 10000, "Maximum number of records in the database")
	rootCmd.PersistentFlags().Duration("db-gc-interval", 10*time.Minute, "BadgerDB value GC interval")
	rootCmd.PersistentFlags().Int("reclaim-delete-batch", 2000, "Maximum number of reclaiming per batch")
	rootCmd.PersistentFlags().Duration("reclaim-interval", 3*time.Minute, "Reclaiming interval")
}
