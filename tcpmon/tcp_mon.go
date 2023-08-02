package tcpmon

import (
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"google.golang.org/protobuf/proto"
)

type SocketMonitor struct{}

func (m *SocketMonitor) Collect(now time.Time) (*KVPair, error) {
	r, _, err := RunSS(now)
	if err != nil {
		return nil, err
	}

	val, err := proto.Marshal(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &KVPair{
		Key:   fmt.Sprintf("%v/%s/", now.UnixMilli(), PrefixTcpRecord),
		Value: val,
	}, nil
}
