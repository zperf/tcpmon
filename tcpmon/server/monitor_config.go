package server

import (
	"time"

	"github.com/spf13/viper"

	"github.com/zperf/tcpmon/tcpmon/tutils"
)

type CmdConfig struct {
	PathSS string
	ArgSS  string

	PathIfconfig string
	ArgIfconfig  string

	Timeout time.Duration
}

func NewCmdConfig() *CmdConfig {
	return &CmdConfig{
		PathSS: tutils.FileFallback(
			viper.GetString("cmd-ss"),
			viper.GetString("cmd-ss2")),
		ArgSS: viper.GetString("cmd-ss-arg"),

		PathIfconfig: tutils.FileFallback(
			viper.GetString("cmd-ifconfig"),
			viper.GetString("cmd-ifconfig2")),
		ArgIfconfig: viper.GetString("cmd-ifconfig-arg"),

		Timeout: viper.GetDuration("cmd-timeout"),
	}
}
