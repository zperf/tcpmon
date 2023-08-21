package cmd

import (
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zperf/tcpmon/tcpmon"
)

var restoreCmd = &cobra.Command{
	Use:   "restore FILE",
	Short: "Restore backup file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		src := args[0]
		output := viper.GetString("output")
		force := viper.GetBool("force")
		needPrint := viper.GetBool("print")

		isEmpty, err := IsDirEmpty(output)
		if err != nil {
			log.Fatal().Err(err).Msg("Check output dir failed")
		}

		if force || isEmpty {
			db, err := badger.Open(badger.DefaultOptions(output).
				WithLogger(&tcpmon.BadgerLogger{}).
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

			if needPrint {
				DoPrint(db, viper.GetString("prefix"), viper.GetBool("reversed"), nil)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().StringP("output", "O", ".", "output database path")
	restoreCmd.Flags().BoolP("force", "f", false, "force restore, may overwrite files")
	restoreCmd.Flags().BoolP("print", "p", false, "print database as JSON after restore")
	addPrintFlags(restoreCmd)

	fatalIf(viper.BindPFlags(restoreCmd.Flags()))
}
