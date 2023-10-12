package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zperf/tcpmon/tcpmon/storage"
	"github.com/zperf/tcpmon/tcpmon/tutils"
)

var exportCmd2 = &cobra.Command{
	Use:   "export2 [-o output] HOSTNAME DATA_FILE",
	Short: "export2 a backup to influxdb line protocol file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		hostname := args[0]
		path := args[1]
		output := viper.GetString("export2-output")
		showOnly := viper.GetBool("export2-show")

		s, err := os.Stat(path)
		if err != nil {
			log.Fatal().Err(err).Str("path", path).Msg("Stat failed")
		}

		var target time.Time
		targetTime := viper.GetString("export2-target-time")
		if targetTime != "" {
			target, err = time.Parse(storage.TimeFormat, targetTime)
			if err != nil {
				log.Fatal().Err(err).Msg("Invalid target time")
			}
		}

		writer := os.Stdout
		if output != "-" {
			fh, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				log.Fatal().Err(err).Msg("Open file for output failed")
			}
			writer = fh
		}

		if s.IsDir() {
			files, err := os.ReadDir(path)
			if err != nil {
				log.Fatal().Err(err).Msg("read dir files failed")
			}

			parsed := false

			for _, f := range files {
				if f.IsDir() {
					continue
				}
				if !strings.HasPrefix(f.Name(), storage.DataFilePrefix) {
					continue
				}

				err = exportFile(filepath.Join(path, f.Name()), writer, hostname, target, showOnly)
				if err != nil {
					if errors.Is(err, storage.ErrTimePointNotIncluded) {
						log.Info().Str("file", f.Name()).
							Msg("File ignored since there is no target time point in it")
					} else {
						log.Fatal().Err(err).Str("file", f.Name()).Msg("Export data files failed")
					}
				} else {
					parsed = true
				}

				if targetTime != "" && parsed {
					log.Info().Str("file", f.Name()).Msg("File exported with target time point")
					break
				}
			}
		} else {
			err = exportFile(path, writer, hostname, target, showOnly)
			if err != nil {
				log.Fatal().Err(err).Msg("Export single data file failed")
			}
		}
	},
}

func exportFile(path string, w io.Writer, hostname string, target time.Time, showOnly bool) error {
	exporter, err := storage.NewFastExporter(path, nil)
	if err != nil {
		return err
	}
	defer exporter.Close()

	return exporter.Export(w, &storage.ExportOptions{Hostname: hostname, Target: target, ShowOnly: showOnly})
}

func init() {
	exportCmd2.Flags().StringP("output", "o", "-",
		"Write output to file instead of stdout. "+
			"Specifying the output as '-' passes the output to stdout.")
	tutils.FatalIf(viper.BindPFlag("export2-output",
		exportCmd2.Flags().Lookup("output")))

	exportCmd2.Flags().StringP("target-time", "t", "",
		"")
	tutils.FatalIf(viper.BindPFlag("export2-target-time", exportCmd2.Flags().Lookup("target-time")))

	exportCmd2.Flags().BoolP("show", "p", false,
		"Print timestamp only")
	tutils.FatalIf(viper.BindPFlag("export2-show", exportCmd2.Flags().Lookup("show")))

	rootCmd.AddCommand(exportCmd2)
}
