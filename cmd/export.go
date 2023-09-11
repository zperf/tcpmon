package cmd

import (
	"net"
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zperf/tcpmon/tcpmon"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export backup file to txt file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		backupFile := args[0]
		hostname := args[1]
		format := viper.GetString("format")
		dbDir := viper.GetString("db-dir")
		force := viper.GetBool("force")

		if net.ParseIP(hostname) == nil {
			log.Fatal().Msg("Invalid IP address")
		}

		err := os.MkdirAll(dbDir, 0755)
		if err != nil && !os.IsExist(err) {
			log.Fatal().Err(err).Msg("Create db directory failed")
		}

		isEmpty, err := IsDirEmpty(dbDir)
		if err != nil {
			log.Fatal().Err(err).Msg("Check db directory failed")
		}

		if force || isEmpty {
			db, err := badger.Open(badger.DefaultOptions(dbDir).
				WithLogger(&tcpmon.BadgerDbLogger{}).
				WithCompactL0OnClose(true))
			if err != nil {
				log.Fatal().Err(err).Msg("Open db for write failed")
			}
			defer db.Close()

			fh, err := os.Open(backupFile)
			if err != nil {
				log.Fatal().Err(err).Msg("Open backup file failed")
			}

			err = db.Load(fh, 256)
			if err != nil {
				log.Fatal().Err(err).Str("backupFile", backupFile).Str("db-dir", dbDir).Msg("Restore failed")
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
						switch format {
						case "tsdb":
							val.TSDBPrintNetstatMetric(hostname)
						default:
							log.Fatal().Msg("Format not supported")
						}
					case "nic":
						var val tcpmon.NicMetric
						err = proto.Unmarshal(valByte, &val)
						if err != nil {
							log.Err(err).Str("key", key).Msg("Unmarshal failed")
						}
						switch format {
						case "tsdb":
							val.TSDBPrintNicMetric(hostname)
						default:
							log.Fatal().Msg("Format not supported")
						}
					case "tcp":
						var val tcpmon.TcpMetric
						err = proto.Unmarshal(valByte, &val)
						if err != nil {
							log.Err(err).Str("key", key).Msg("Unmarshal failed")
						}
						switch format {
						case "tsdb":
							val.TSDBPrintTcpMetric(hostname)
						default:
							log.Fatal().Msg("Format not supported")
						}
					default:
						log.Warn().Str("key", key).Msg("wrong key format")
					}
				}
				return nil
			})
		} else {
			log.Warn().Msg("db-dir is not empty, please clear db-dir or use '--force'")
			return
		}
	},
}

func init() {
	exportCmd.PersistentFlags().StringP("format", "f", "tsdb", "export backup to txt in this format")
	exportCmd.PersistentFlags().StringP("db-dir", "d", "/tmp/tcpmon/export/db", "restore backup to db-dir")
	exportCmd.PersistentFlags().BoolP("force", "e", false, "force restore, may overwrite files")
	fatalIf(viper.BindPFlags(exportCmd.PersistentFlags()))
	rootCmd.AddCommand(exportCmd)
}
