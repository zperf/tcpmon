package tcpmon

import (
	"strings"
	"time"
	"unicode"

	"github.com/hashicorp/memberlist"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type GossipServer struct {
	cluster *memberlist.Memberlist
}

func NewGossipServer() *GossipServer {
	m, err := memberlist.Create(memberlist.DefaultLocalConfig())
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create memberlist")
	}

	gossipCluster := strings.FieldsFunc(viper.GetString("gossip-cluster"), func(c rune) bool {
		return unicode.IsSpace(c) || c == ','
	})

	for {
		n, err := m.Join(gossipCluster)
		if err != nil {
			log.Err(err).Msg("fail to join cluster")
			time.Sleep(1 * time.Second)
			continue
		}
		log.Info().Int("member number", n).Msg("success to join cluster")
		break
	}

	return &GossipServer{
		cluster: m,
	}
}
