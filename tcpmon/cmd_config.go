package tcpmon

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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

func NewCmdConfig() *CmdConfig {
	ss := viper.GetString("cmd-ss")
	ok, err := FileExists(ss)
	if err != nil {
		log.Fatal().Err(err).Msg("Stat command ss failed")
	}

	if !ok {
		ss = viper.GetString("cmd-ss2")
		ok, err = FileExists(ss)
		if err != nil {
			log.Fatal().Err(err).Msg("Stat command ss2 failed")
		}
		if !ok {
			log.Fatal().Err(errors.New("command ss not found")).Msg("Please install iproute or iproute2")
		}
	}

	return &CmdConfig{
		PathSS:       viper.GetString("cmd-ss"),
		ArgSS:        viper.GetString("cmd-ss-arg"),
		PathIfconfig: viper.GetString("cmd-ifconfig"),
		ArgIfconfig:  viper.GetString("cmd-ifconfig-arg"),
		PathNetstat:  viper.GetString("cmd-netstat"),
		ArgNetstat:   viper.GetString("cmd-netstat-arg"),
		Timeout:      viper.GetDuration("cmd-timeout"),
	}
}
