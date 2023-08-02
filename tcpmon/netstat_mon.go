package tcpmon

import (
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"google.golang.org/protobuf/proto"
)

type NetstatMonitor struct{}

func (m *NetstatMonitor) Collect(now time.Time) (*KVPair, error) {
	r, _, err := RunNetstat(now)
	if err != nil {
		return nil, err
	}

	val, err := proto.Marshal(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &KVPair{
		Key:   fmt.Sprintf("%v/%s/", now.UnixMilli(), PrefixNetRecord),
		Value: val,
	}, nil
}
