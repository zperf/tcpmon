package tcpmon

import (
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"google.golang.org/protobuf/proto"
)

type SocketMonitor struct {
	config *CmdConfig
}

func (m *SocketMonitor) Collect(now time.Time) (*KVPair, error) {
	r, _, err := m.RunSS(now)
	if err != nil {
		return nil, err
	}

	val, err := proto.Marshal(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &KVPair{
		Key:   fmt.Sprintf("%s/%v/", PrefixTcpMetric, now.UnixMilli()),
		Value: val,
	}, nil
}
