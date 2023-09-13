package tcpmon

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/zperf/tcpmon/tcpmon/tutils"
)

type CmdConfig struct {
	PathSS       string
	ArgSS        string
	PathIfconfig string
	ArgIfconfig  string
	PathNetstat  string
	ArgNetstat   string
	Timeout      time.Duration
}

func FileFallback(path ...string) string {
	for _, p := range path {
		ok, err := tutils.FileExists(p)
		if err != nil {
			log.Fatal().Err(err).Str("file", p).Msg("Stat file failed")
		}
		if ok {
			return p
		}
	}
	return ""
}

func NewCmdConfig() *CmdConfig {
	return &CmdConfig{
		PathSS:       FileFallback(viper.GetString("cmd-ss"), viper.GetString("cmd-ss2")),
		ArgSS:        viper.GetString("cmd-ss-arg"),
		PathIfconfig: FileFallback(viper.GetString("cmd-ifconfig"), viper.GetString("cmd-ifconfig2")),
		ArgIfconfig:  viper.GetString("cmd-ifconfig-arg"),
		PathNetstat:  viper.GetString("cmd-netstat"),
		ArgNetstat:   viper.GetString("cmd-netstat-arg"),
		Timeout:      viper.GetDuration("cmd-timeout"),
	}
}
