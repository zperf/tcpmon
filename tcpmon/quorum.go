package tcpmon

import (
	"net"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/memberlist"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type Quorum struct {
	mlist *memberlist.Memberlist
	ds    *DataStore
	// listenAddr IpAddr
}

func NewQuorum(ds *DataStore, monitorConfig *MonitorConfig) *Quorum {
	q := &Quorum{
		ds: ds,
	}

	// create memberlist
	config := memberlist.DefaultLANConfig()
	config.Events = q
	config.LogOutput = NewMemberlistLogger()
	config.BindPort = monitorConfig.QuorumPort
	config.AdvertisePort = monitorConfig.QuorumPort

	m, err := memberlist.Create(config)
	if err != nil {
		log.Fatal().Err(err).Msg("create memberlist failed")
	}
	q.mlist = m

	return q
}

func (q *Quorum) Close() {
	err := q.mlist.Shutdown()
	if err != nil {
		log.Warn().Err(err).Msg("Shutdown quorum failed")
	}
}

func writeConfig() {
	err := viper.WriteConfig()
	if err != nil {
		log.Warn().Err(err).Msg("Write members to config failed")
	}
}

func (q *Quorum) configMemberJoin(member string, meta string) {
	members := viper.GetStringMapString("members")
	_, ok := members[member]
	if ok {
		log.Warn().Str("member", member).Msg("Already in the quorum")
		return
	}

	members[member] = meta
	viper.Set("members", members)
	writeConfig()
}

func (q *Quorum) configMemberLeave(member string) {
	members := viper.GetStringMapString("members")
	delete(members, member)
	viper.Set("members", members)
	writeConfig()
}

func (q *Quorum) configMemberUpdate(member string, meta string) {
	members := viper.GetStringMapString("members")
	members[member] = meta
	viper.Set("members", members)
	writeConfig()
}

func (q *Quorum) NotifyJoin(node *memberlist.Node) {
	log.Info().Str("node", node.Address()).
		Str("meta", string(node.Meta)).
		Msgf("Member joined")
	q.configMemberJoin(node.Address(), string(node.Meta))

}

func (q *Quorum) NotifyLeave(node *memberlist.Node) {
	log.Info().Str("node", node.Address()).
		Str("meta", string(node.Meta)).
		Msgf("Member left quorum")
	q.configMemberLeave(node.Address())
}

func (q *Quorum) NotifyUpdate(node *memberlist.Node) {
	log.Info().Str("node", node.Address()).
		Str("meta", string(node.Meta)).
		Msgf("Update member meta data")
	q.configMemberUpdate(node.Address(), string(node.Meta))
}

func (q *Quorum) Local() *memberlist.Node {
	return q.mlist.LocalNode()
}

func (q *Quorum) Members() []*memberlist.Node {
	return q.mlist.Members()
}

func (q *Quorum) Leave(timeout time.Duration) error {
	return q.mlist.Leave(timeout)
}

func (q *Quorum) TryJoin(members map[string]string) (int, error) {
	m := lo.Keys(members)
	if len(m) > 0 {
		return q.mlist.Join(m)
	}
	return 0, errors.New("members is empty")
}

func (q *Quorum) Join(members map[string]string, retry int, delay time.Duration) {
	if retry == 0 {
		retry = 3
	}
	if delay == 0 {
		delay = time.Second
	}

	for i := 0; i < retry; i++ {
		_, err := q.TryJoin(members)
		if err != nil {
			log.Err(err).Msg("join quorum failed")
		} else {
			log.Info().Msg("join quorum success")
			break
		}
		time.Sleep(delay)
	}
}

func (q *Quorum) MyIP() net.IP {
	return net.ParseIP(GetIpFromAddress(q.mlist.LocalNode().Address()))
}

func (q *Quorum) My() *memberlist.Node {
	return q.mlist.LocalNode()
}

func (q *Quorum) GetMemberMeta(member string) (string, error) {
	members := viper.GetStringMapString("members")
	if members == nil {
		return "", errors.New("Get members failed")
	}

	m, ok := members[member]
	if !ok {
		return "", errors.Newf("Member %s not in the cluster", member)
	}

	return m, nil
}
