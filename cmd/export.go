package cmd

import (
	"github.com/cockroachdb/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	. "github.com/zperf/tcpmon/tcpmon"
	storagev2 "github.com/zperf/tcpmon/tcpmon/storage/v2"
	"github.com/zperf/tcpmon/tcpmon/tutils"
)

var FlagExportFormat = exportFormatTsdb
var FlagHostname string

var exportCmd = &cobra.Command{
	Use:   "export [BASE_DIR]",
	Short: "export backup file to txt file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		baseDir := args[0]

		if FlagHostname == "" {
			log.Fatal().Msg("hostname is empty. Try adds `--hostname` to set")
		}

		var printer MetricPrinter
		switch FlagExportFormat.String() {
		case "tsdb":
			printer = TSDBMetricPrinter{}
		}

		reader, err := storagev2.NewDataStoreReader(storagev2.NewReaderConfig(baseDir))
		if err != nil {
			log.Fatal().Err(err).Msg("Open datastore failed")
		}
		defer reader.Close()

		err = reader.Iterate(func(buf []byte) {
			log.Info().Int("bufLen", len(buf)).Msg("Read buffer")

			var msg Metric
			err := proto.Unmarshal(buf, &msg)
			if err != nil {
				log.Fatal().Err(err).Msg("Unmarshal failed")
			}

			switch m := msg.Body.(type) {
			case *Metric_Tcp:
				printer.PrintTcpMetric(m.Tcp, FlagHostname)
			case *Metric_Net:
				printer.PrintNetstatMetric(m.Net, FlagHostname)
			case *Metric_Nic:
				printer.PrintNicMetric(m.Nic, FlagHostname)
			}
		})
		if err != nil {
			log.Err(err).Msg("Read db failed")
		}
	},
}

type ExportFormat string

const (
	exportFormatTsdb ExportFormat = "tsdb"
)

func (f *ExportFormat) String() string {
	return string(*f)
}

func (f *ExportFormat) Set(v string) error {
	switch v {
	case "tsdb":
		*f = ExportFormat(v)
		return nil
	default:
		return errors.New("Export to tsdb format is the only supported")
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
	rootCmd.AddCommand(exportCmd)
}
