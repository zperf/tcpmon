package tcpmon

import (
	"net"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/rs/zerolog/log"
)

type Quorum struct {
	mlist *memberlist.Memberlist
	ds    *Datastore
}

func NewQuorum(ds *Datastore) *Quorum {
	q := &Quorum{
		ds: ds,
	}

	config := memberlist.DefaultLANConfig()
	config.Events = q
	config.LogOutput = &MemberlistLogger{}

	m, err := memberlist.Create(config)
	if err != nil {
		log.Fatal().Err(err).Msg("create memberlist failed")
	}
	q.mlist = m

	address := m.LocalNode().Address()
	log.Info().Str("hostname", m.LocalNode().String()).
		Str("address", address).
		Msg("Quorum created")

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
	err := q.ds.AddMember(node.Address())
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
	log.Info().Str("node", node.Address()).Msgf("Member data updated")
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
