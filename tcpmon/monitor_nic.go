package tcpmon

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gogo/protobuf/proto"

	v2 "github.com/zperf/tcpmon/tcpmon/storage/v2"
)

type NicMonitor struct{ config *CmdConfig }

func (m *NicMonitor) Collect(now time.Time) (*v2.MetricContext, error) {
	r, _, err := m.RunIfconfig(now)
	if err != nil {
		return nil, err
	}

	metric := &Metric{Body: &Metric_Nic{Nic: r}}
	val, err := proto.Marshal(metric)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &v2.MetricContext{
		Value: val,
	}, nil
}
