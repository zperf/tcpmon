package cmd

import (
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zperf/tcpmon/tcpmon"
)

var format string
var dbDir string
var force bool

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export backup file to txt file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		backupFile := args[0]
		hostname := args[1]

		var printer tcpmon.MetricPrinter
		switch format {
		case "tsdb":
			printer = tcpmon.TSDBMetricPrinter{}
		default:
			log.Fatal().Msg("Format not supported")
		}

		err := os.MkdirAll(dbDir, 0755)
		if err != nil && !os.IsExist(err) {
			log.Fatal().Err(err).Msg("Create db directory failed")
		}

		isEmpty, err := IsDirEmpty(dbDir)
		if err != nil {
			log.Fatal().Err(err).Msg("Check db directory failed")
		}

		if !force && !isEmpty {
			log.Fatal().Msg("db is not empty, please clear db or use '-f'")
		}

		db, err := badger.Open(badger.DefaultOptions(dbDir).
			WithLogger(&tcpmon.BadgerDbLogger{}))
		if err != nil {
			log.Fatal().Err(err).Msg("Open db for write failed")
		}
		defer db.Close()

		err = db.DropAll()
		if err != nil {
			log.Fatal().Err(err).Msg("Clear db failed")
		}

		fh, err := os.Open(backupFile)
		if err != nil {
			log.Fatal().Err(err).Msg("Open backup file failed")
		}

		err = db.Load(fh, 256)
		if err != nil {
			log.Fatal().Err(err).Str("backupFile", backupFile).Str("db", dbDir).Msg("Restore failed")
		}

		err = db.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			it := txn.NewIterator(opts)
			defer it.Close()
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				key := string(item.Key())
				valByte, err := item.ValueCopy(nil)
				if err != nil {
					log.Err(err).Str("key", key).Msg("Get value failed")
				}
				switch key[0:3] {
				case "net":
					var val tcpmon.NetstatMetric
					err = proto.Unmarshal(valByte, &val)
					if err != nil {
						log.Err(err).Str("key", key).Msg("Unmarshal failed")
					}
					printer.PrintNetstatMetric(&val, hostname)
				case "nic":
					var val tcpmon.NicMetric
					err = proto.Unmarshal(valByte, &val)
					if err != nil {
						log.Err(err).Str("key", key).Msg("Unmarshal failed")
					}
					printer.PrintNicMetric(&val, hostname)
				case "tcp":
					var val tcpmon.TcpMetric
					err = proto.Unmarshal(valByte, &val)
					if err != nil {
						log.Err(err).Str("key", key).Msg("Unmarshal failed")
					}
					printer.PrintTcpMetric(&val, hostname)
				default:
					log.Warn().Str("key", key).Msg("wrong key format")
				}
			}
			return nil
		})
		if err != nil {
			log.Err(err).Msg("Read db failed")
		}
	},
}

func init() {
	exportCmd.Flags().StringVar(&format, "format", "tsdb", "export backup to txt in this format")
	exportCmd.Flags().StringVarP(&dbDir, "db", "d", "/tmp/tcpmon/export/db", "db path to restore backup")
	exportCmd.Flags().BoolVarP(&force, "force", "f", false, "force restore, may overwrite files")
	rootCmd.AddCommand(exportCmd)
}
