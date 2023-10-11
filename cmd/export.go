package cmd

import (
	"os"

	"github.com/cockroachdb/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/zperf/tcpmon/tcpmon/export/influxdb"
	"github.com/zperf/tcpmon/tcpmon/storage"
	"github.com/zperf/tcpmon/tcpmon/tproto"
	"github.com/zperf/tcpmon/tcpmon/tutils"
)

var FlagExportFormat = exportFormatInfluxdb
var FlagHostname string
var FlagOutput string

var exportCmd = &cobra.Command{
	Use:   "export [-o output] [-h hostname] BASE_DIR",
	Short: "export a backup to influxdb line protocol file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		baseDir := args[0]

		writer := os.Stdout
		if FlagOutput != "-" {
			w, err := os.OpenFile(FlagOutput, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal().Err(err).Msg("Open file for output writing failed")
			}
			writer = w
		}

		var exporter *influxdb.Exporter
		switch FlagExportFormat.String() {
		case "influxdb":
			exporter = influxdb.New(FlagHostname, writer)
		default:
			log.Fatal().Str("format", FlagExportFormat.String()).Msg("")
		}

		reader, err := storage.NewDataStoreReader(storage.NewReaderConfig(baseDir))
		if err != nil {
			log.Fatal().Err(err).Msg("Open datastore failed")
		}
		defer reader.Close()

		err = reader.Iterate(func(buf []byte) {
			var m tproto.Metric
			err := proto.Unmarshal(buf, &m)
			if err != nil {
				log.Fatal().Err(err).Msg("Unmarshal failed")
			}

			exporter.ExportMetric(&m)
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Read db failed")
		}
	},
}

type ExportFormat string

const (
	exportFormatInfluxdb ExportFormat = "influxdb"
)

func (f *ExportFormat) String() string {
	return string(*f)
}

func (f *ExportFormat) Set(v string) error {
	switch v {
	case "influxdb":
		*f = ExportFormat(v)
		return nil
	default:
		return errors.New("Export to influxdb format is the only supported")
	}
}

func (f *ExportFormat) Type() string {
	return "ExportFormat"
}

func init() {
	exportCmd.Flags().Var(&FlagExportFormat, "format",
		"export backup to txt in this format")
	exportCmd.Flags().StringVarP(&FlagHostname, "hostname", "n", tutils.Hostname(),
		"export backup to txt in this format")
	exportCmd.Flags().StringVarP(&FlagOutput, "output", "o", "-",
		"Output to file")
	rootCmd.AddCommand(exportCmd)
}
