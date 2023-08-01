package tcpmon

import (
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"google.golang.org/protobuf/proto"
)

type NetstatMonitor struct{}

func (m *NetstatMonitor) Collect(now time.Time) (*StoreRequest, error) {
	r, _, err := RunNetstat()
	if err != nil {
		return nil, err
	}

	val, err := proto.Marshal(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &StoreRequest{
		Key:   fmt.Sprintf("%s/%v/", PrefixNetRecord, now.UnixMilli()),
		Value: val,
	}, nil
}
