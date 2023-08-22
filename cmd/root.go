package cmd

import (
	"os"
	"strings"

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
	initLoggingFlags()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initLoggingFlags() {
	startCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode")
	startCmd.PersistentFlags().String("log-level", "info", "log level")
	startCmd.PersistentFlags().String("log-dir", "/tmp/tcpmon/log", "The log directory")
	startCmd.PersistentFlags().String("log-filename", "tcpmon.log", "The file name of logs")
	startCmd.PersistentFlags().Int("log-max-size", 10, "Maximum size of each log file")
	startCmd.PersistentFlags().Int("log-max-count", 5, "Maximum log files to keep")
	fatalIf(viper.BindPFlags(startCmd.PersistentFlags()))
}
