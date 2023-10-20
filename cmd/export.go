package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zperf/tcpmon/tcpmon/export/influxdb"
	"github.com/zperf/tcpmon/tcpmon/storage"
	"github.com/zperf/tcpmon/tcpmon/tutils"
)

var exportCmd = &cobra.Command{
	Use:   "export [-o output] HOSTNAME DATA_FILE_OR_DIR",
	Short: "export a backup to influxdb line protocol file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		hostname := args[0]
		path := args[1]
		output := viper.GetString("export-output")
		showOnly := viper.GetBool("export-show")

		s, err := os.Stat(path)
		if err != nil {
			log.Fatal().Err(err).Str("path", path).Msg("Stat failed")
		}

		var target time.Time
		targetTime := viper.GetString("export-target-time")
		if targetTime != "" {
			target, err = time.Parse(tutils.TimeFormat, targetTime)
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

		bar := progressbar.Default(-1, "Exporting")
		defer bar.Close()

		token := viper.GetString("export-token")
		exportOption := influxdb.ExportOptions{
			Hostname:  hostname,
			Target:    target,
			ShowOnly:  showOnly,
			WriteDb:   token != "",
			Org:       viper.GetString("export-org"),
			Bucket:    viper.GetString("export-bucket"),
			Token:     token,
			DbAddress: viper.GetString("export-db"),
			Bar:       bar,
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

				err = exportFile(filepath.Join(path, f.Name()), writer, &exportOption)
				if err != nil {
					if errors.Is(err, influxdb.ErrTimePointNotIncluded) {
						if exportOption.Bar != nil {
							bar.Describe(fmt.Sprintf("Ignore %s", path))
						} else {
							log.Info().Str("file", f.Name()).
								Msg("File ignored since there is no target time point in it")
						}
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
			err = exportFile(path, writer, &exportOption)
			if err != nil {
				log.Fatal().Err(err).Msg("Export single data file failed")
			}
		}
	},
}

func exportFile(path string, w io.Writer, options *influxdb.ExportOptions) error {
	if options.Bar != nil {
		options.Bar.Describe(path)
	}

	exporter, err := influxdb.NewFastExporter(path, nil)
	if err != nil {
		return err
	}
	defer exporter.Close()

	return exporter.Export(w, options)
}

func init() {
	exportCmd.Flags().StringP("output", "o", "-",
		"Write output to file instead of stdout. "+
			"Specifying the output as '-' passes the output to stdout.")
	tutils.FatalIf(viper.BindPFlag("export-output",
		exportCmd.Flags().Lookup("output")))

	exportCmd.Flags().StringP("target-time", "t", "",
		"The target time point that needs to be imported")
	tutils.FatalIf(viper.BindPFlag("export-target-time", exportCmd.Flags().Lookup("target-time")))

	exportCmd.Flags().BoolP("show", "p", false,
		"Print timestamp only")
	tutils.FatalIf(viper.BindPFlag("export-show", exportCmd.Flags().Lookup("show")))

	// for db
	exportCmd.Flags().StringP("bucket", "b", "",
		"InfluxDB bucket name")
	tutils.FatalIf(viper.BindPFlag("export-bucket", exportCmd.Flags().Lookup("bucket")))

	exportCmd.Flags().String("org", "",
		"InfluxDB org name")
	tutils.FatalIf(viper.BindPFlag("export-org", exportCmd.Flags().Lookup("org")))

	exportCmd.Flags().String("token", "",
		"InfluxDB connection token")
	tutils.FatalIf(viper.BindPFlag("export-token", exportCmd.Flags().Lookup("token")))

	exportCmd.Flags().String("db", "http://127.0.0.1:8086",
		"InfluxDB connection address")
	tutils.FatalIf(viper.BindPFlag("export-db", exportCmd.Flags().Lookup("db")))

	rootCmd.AddCommand(exportCmd)
}
