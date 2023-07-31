package tcpmon

import (
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/golang/protobuf/proto"
)

type SocketMonitor struct{}

func (m *SocketMonitor) Collect(now time.Time) (*StoreRequest, error) {
	r, _, err := ss(now)
	if err != nil {
		return nil, err
	}

	val, err := proto.Marshal(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &StoreRequest{
		Key:   fmt.Sprintf("%s/%v/", PrefixSocketRecord, now.UnixMilli()),
		Value: val,
	}, nil
}
