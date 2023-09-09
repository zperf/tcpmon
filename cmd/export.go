package cmd

import (
	"github.com/cockroachdb/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/zperf/tcpmon/tcpmon"
	storagev2 "github.com/zperf/tcpmon/tcpmon/storage/v2"
)

var FlagExportFormat = exportFormatTsdb

var exportCmd = &cobra.Command{
	Use:   "export [BASE_DIR] [HOSTNAME]",
	Short: "export backup file to txt file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		baseDir := args[0]
		hostname := args[1]

		var printer tcpmon.MetricPrinter
		switch FlagExportFormat.String() {
		case "tsdb":
			printer = tcpmon.TSDBMetricPrinter{}
		}

		reader, err := storagev2.NewDataStoreReader(baseDir, nil)
		if err != nil {
			log.Fatal().Err(err).Msg("Open datastore failed")
		}
		defer reader.Close()

		err = reader.Iterate(func(buf []byte) {
			log.Info().Int("bufLen", len(buf)).Msg("Read buffer")

			var msg tcpmon.Metric
			err := proto.Unmarshal(buf, &msg)
			if err != nil {
				log.Fatal().Err(err).Msg("Unmarshal failed")
			}

			switch m := msg.Body.(type) {
			case *tcpmon.Metric_Tcp:
				printer.PrintTcpMetric(m.Tcp, hostname)
			case *tcpmon.Metric_Net:
				printer.PrintNetstatMetric(m.Net, hostname)
			case *tcpmon.Metric_Nic:
				printer.PrintNicMetric(m.Nic, hostname)
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
	exportCmd.Flags().Var(&FlagExportFormat, "format", "export backup to txt in this format")
	rootCmd.AddCommand(exportCmd)
}
