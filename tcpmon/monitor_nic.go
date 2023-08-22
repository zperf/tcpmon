package tcpmon

import (
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"google.golang.org/protobuf/proto"
)

type NicMonitor struct{ config *CmdConfig }

func (m *NicMonitor) Collect(now time.Time) (*KVPair, error) {
	r, _, err := m.RunIfconfig(now)
	if err != nil {
		return nil, err
	}

	val, err := proto.Marshal(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &KVPair{
		Key:   fmt.Sprintf("%s/%v/", PrefixNicMetric, now.UnixMilli()),
		Value: val,
	}, nil
}
