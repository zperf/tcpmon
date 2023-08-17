package tcpmon

import (
	"strings"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/rs/zerolog/log"
)

type GossipServer struct {
	cluster *memberlist.Memberlist
}

func (g *GossipServer) Join(clusterAddr []string) {
	for {
		_, err := g.cluster.Join(clusterAddr)
		if err != nil {
			log.Err(err).Str("clusterAddr", strings.Join(clusterAddr, ", ")).Msg("fail to join cluster")
			time.Sleep(1 * time.Second)
			continue
		}
		log.Info().Str("clusterAddr", strings.Join(clusterAddr, ", ")).Msg("success to join cluster")
		break
	}
}

func NewGossipServer() *GossipServer {
	m, err := memberlist.Create(memberlist.DefaultLocalConfig())
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create memberlist")
	}

	return &GossipServer{
		cluster: m,
	}
}
