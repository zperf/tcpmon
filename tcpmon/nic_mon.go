package tcpmon

import (
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"google.golang.org/protobuf/proto"
)

type NicMonitor struct{}

func (m *NicMonitor) Collect(now time.Time) (*StoreRequest, error) {
	r, _, err := RunIfconfig(now)
	if err != nil {
		return nil, err
	}

	val, err := proto.Marshal(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &StoreRequest{
		Key:   fmt.Sprintf("%s/%s/", PrefixNicRecord, now.UnixMilli()),
		Value: val,
	}, nil
}
