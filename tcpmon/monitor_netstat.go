package tcpmon

import (
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gogo/protobuf/proto"
)

type NetstatMonitor struct{ config *CmdConfig }

func (m *NetstatMonitor) Collect(now time.Time) (*KVPair, error) {
	r, _, err := m.RunNetstat(now)
	if err != nil {
		return nil, err
	}

	val, err := proto.Marshal(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &KVPair{
		Key:   fmt.Sprintf("%s/%v/", PrefixNetMetric, now.UnixMilli()),
		Value: val,
	}, nil
}
