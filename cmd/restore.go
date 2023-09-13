package cmd

import (
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1 "github.com/zperf/tcpmon/tcpmon/storage/v1"
)

var restoreCmd = &cobra.Command{
	Use:   "restore FILE",
	Short: "Restore backup file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		src := args[0]
		output := viper.GetString("output")
		force := viper.GetBool("force")

		err := os.MkdirAll(output, 0755)
		if err != nil && !os.IsExist(err) {
			log.Fatal().Err(err).Msg("Create output directory failed")
		}

		isEmpty, err := IsDirEmpty(output)
		if err != nil {
			log.Fatal().Err(err).Msg("Check output dir failed")
		}

		if force || isEmpty {
			db, err := badger.Open(badger.DefaultOptions(output).
				WithLogger(v1.NewBadgerLogger()).
				WithCompactL0OnClose(true))
			if err != nil {
				log.Fatal().Err(err).Msg("Open db for write failed")
			}
			defer db.Close()

			fh, err := os.Open(src)
			if err != nil {
				log.Fatal().Err(err).Msg("Open backup file failed")
			}

			err = db.Load(fh, 256)
			if err != nil {
				log.Fatal().Err(err).Str("backupFile", src).Str("Output", output).Msg("Restore failed")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().StringP("output", "O", ".", "output database path")
	restoreCmd.Flags().BoolP("force", "f", false, "force restore, may overwrite files")
	restoreCmd.Flags().Bool("print", false, "print database as JSON after restore")
}
