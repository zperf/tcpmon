package cmd

import (
	"github.com/spf13/cobra"
)

var dbSetCmd = &cobra.Command{
	Use:     "set [KEY] [VALUE-FILE]",
	Short:   "Set value to the db",
	Example: "  set --db ./db tcp/1692679337381 1.json",
	Args:    cobra.ExactArgs(2),
	//Run: func(cmd *cobra.Command, args []string) {
	//	path := viper.GetString("db")
	//	key := args[0]
	//	valueFile := args[1]
	//
	//	prefixType := GetPrefixType(key)
	//	if !ValidPrefix(prefixType) {
	//		log.Fatal().Err(errors.Newf("invalid metric type '%s'", prefixType)).Msg("validation failed")
	//	}
	//
	//	fh, err := os.Open(valueFile)
	//	if err != nil {
	//		log.Fatal().Err(err).Msg("Open value file failed")
	//	}
	//	defer fh.Close()
	//
	//	db := openBadger(path)
	//	defer db.Close()
	//
	//	buf, err := io.ReadAll(fh)
	//	if err != nil {
	//		log.Fatal().Err(err).Msg("Read value file failed")
	//	}
	//
	//	if prefixType == PrefixMember {
	//		err = db.Update(func(txn *badger.Txn) error {
	//			return txn.Set([]byte(key), buf)
	//		})
	//		if err != nil {
	//			log.Fatal().Err(err).Msg("Set value failed")
	//		}
	//		return
	//	}
	//
	//	var m proto.Message
	//	switch prefixType {
	//	case PrefixNetMetric:
	//		var metric NetstatMetric
	//		m = &metric
	//
	//	case PrefixTcpMetric:
	//		var metric TcpMetric
	//		m = &metric
	//
	//	case PrefixNicMetric:
	//		var metric NicMetric
	//		m = &metric
	//	}
	//	err = protojson.Unmarshal(buf, m)
	//	if err != nil {
	//		log.Fatal().Err(err).Msg("Unmarshal failed")
	//	}
	//
	//	val, err := proto.Marshal(m)
	//	if err != nil {
	//		log.Fatal().Err(err).Msg("Marshal failed")
	//	}
	//
	//	err = db.Update(func(txn *badger.Txn) error {
	//		return txn.Set([]byte(key), val)
	//	})
	//	if err != nil {
	//		log.Fatal().Err(err).Msg("Update failed")
	//	}
	//},
}

func init() {
	dbCmd.AddCommand(dbSetCmd)
}
