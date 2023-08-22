package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get and set tcpmon options",
}

var configGetDefaultCmd = &cobra.Command{
	Use:   "default",
	Short: "Generate default config file",
	Run: func(cmd *cobra.Command, args []string) {
		options := viper.AllSettings()
		out, err := yaml.Marshal(options)
		if err != nil {
			log.Fatal().Err(err).Msg("marshal to yaml failed")
		}
		fmt.Println(string(out))
	},
}

func init() {
	configCmd.AddCommand(configGetDefaultCmd)
	rootCmd.AddCommand(configCmd)
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
