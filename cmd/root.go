package cmd

import (
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zperf/tcpmon/tcpmon"
)

var rootCmd = &cobra.Command{
	Use:   "tcpmon",
	Short: "Tcpmon is a portable local network monitor for Linux",
}

func Execute(cmdline string) {
	// init viper
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/tcpmon/")
	viper.AddConfigPath("$HOME/.tcpmon")

	// read config file
	if !strings.HasPrefix(cmdline, "config default") {
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
		if strings.HasPrefix(cmdline, "start") {
			log.Info().Str("configFile", viper.ConfigFileUsed()).
				Str("logDir", viper.GetString("log-dir")).
				Msg("Config loaded")
		}
	}

	// init logger
	level, _ := zerolog.ParseLevel(viper.GetString("log-level"))
	logConfig := &tcpmon.LogConfig{
		Level:                 level,
		FileLoggingEnabled:    true,
		ConsoleLoggingEnabled: !viper.GetBool("disable-console"),
		Directory:             viper.GetString("log-dir"),
		Filename:              viper.GetString("log-filename"),
		MaxSize:               viper.GetInt("log-max-size"),
		MaxBackups:            viper.GetInt("log-max-count"),
	}
	tcpmon.InitLogger(logConfig)

	// print warnings after logger initialized
	if level == zerolog.NoLevel {
		log.Warn().Str("level", viper.GetString("log-level")).
			Msg("invalid level, default to NoLevel")
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("disable-console", false, "Disable log to console")
	rootCmd.PersistentFlags().String("log-level", "info", "log level")
	rootCmd.PersistentFlags().String("log-dir", "/tmp/tcpmon/log", "The log directory")
	rootCmd.PersistentFlags().String("log-filename", "tcpmon.log", "The file name of logs")
	rootCmd.PersistentFlags().Int("log-max-size", 10, "Maximum size of each log file")
	rootCmd.PersistentFlags().Int("log-max-count", 5, "Maximum log files to keep")
	fatalIf(viper.BindPFlags(rootCmd.PersistentFlags()))
}
