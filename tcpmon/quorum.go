package tcpmon

import (
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type Quorum struct {
	mlist *memberlist.Memberlist
	ds    *Datastore
}

func NewQuorum(ds *Datastore, mconfig *MonitorConfig) *Quorum {
	q := &Quorum{
		ds: ds,
	}

	config := memberlist.DefaultLANConfig()
	config.Events = q
	config.LogOutput = NewMemberlistLogger()
	config.BindPort = mconfig.QuorumPort
	config.AdvertisePort = mconfig.QuorumPort

	m, err := memberlist.Create(config)
	if err != nil {
		log.Fatal().Err(err).Msg("create memberlist failed")
	}
	q.mlist = m

	my := m.LocalNode()
	log.Info().Str("hostname", my.String()).
		Str("address", my.Address()).
		Msg("Quorum created")

	// update local meta
	var memberInfo MemberInfo
	ipaddr := ParseIpAddr(mconfig.HttpListen)
	ipaddr.Address = my.Addr.String()
	memberInfo.HttpListen = fmt.Sprintf("http://%s", ipaddr.String())
	buf, err := proto.Marshal(&memberInfo)
	if err != nil {
		log.Fatal().Err(err).Msg("Marshal member info failed")
	}
	my.Meta = buf
	err = ds.UpdateMember(my.Address(), buf)
	if err != nil {
		log.Fatal().Err(err).Msg("Update my member info failed")
	}

	return q
}

func (q *Quorum) Close() {
	err := q.mlist.Shutdown()
	if err != nil {
		log.Warn().Err(err).Msg("Shutdown quorum failed")
	}
}

func (q *Quorum) NotifyJoin(node *memberlist.Node) {
	log.Info().Str("node", node.Address()).Msgf("Member joined")
	err := q.ds.AddMember(node.Address(), node.Meta)
	if err != nil {
		log.Warn().Err(err).Str("node", node.Address()).Msg("Save member failed")
	}
}

func (q *Quorum) NotifyLeave(node *memberlist.Node) {
	log.Info().Str("node", node.Address()).Msgf("Member left")
	err := q.ds.DeleteMember(node.Address())
	if err != nil {
		log.Warn().Err(err).Str("node", node.Address()).Msg("Save member failed")
	}
}

func (q *Quorum) NotifyUpdate(node *memberlist.Node) {
	log.Info().Str("node", node.Address()).Msgf("Update member meta data")
	err := q.ds.UpdateMember(node.Address(), node.Meta)
	if err != nil {
		log.Warn().Err(err).Str("node", node.Address()).Msg("Update member failed")
	}
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

func (q *Quorum) TryJoin(members []string) (int, error) {
	return q.mlist.Join(members)
}

func (q *Quorum) Join(members []string) {
	for i := 0; i < 3; i++ {
		_, err := q.TryJoin(members)
		if err != nil {
			log.Err(err).Strs("members", members).Msg("join quorum failed")
		} else {
			log.Info().Strs("members", members).Msg("join quorum success")
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func (q *Quorum) MyIP() net.IP {
	return net.ParseIP(GetIpFromAddress(q.mlist.LocalNode().Address()))
}

func (q *Quorum) GetMemberMeta(member string) (map[string]any, error) {
	return q.ds.GetMemberMeta(member)
}
