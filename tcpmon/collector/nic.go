package collector

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
	"github.com/gogo/protobuf/proto"

	"github.com/zperf/tcpmon/tcpmon/parsing"
	"github.com/zperf/tcpmon/tcpmon/tproto"
)

type NicCollector struct{ config *Config }

func NewNic(config *Config) *NicCollector {
	return &NicCollector{config: config}
}

func (m *NicCollector) Collect(now time.Time) ([]byte, error) {
	r, err := m.doCollect(now)
	if err != nil {
		return nil, err
	}

	metric := &tproto.Metric{Body: &tproto.Metric_Nic{Nic: r}}
	val, err := proto.Marshal(metric)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return val, nil
}

func (m *NicCollector) doCollect(now time.Time) (*tproto.NicMetric, error) {
	c := cmd.NewCmd(m.config.PathIfconfig)
	ctx, cancel := context.WithTimeout(context.Background(), m.config.Timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "ifconfig timeout")
	case st := <-c.Start():
		var nics tproto.NicMetric
		nics.Type = tproto.MetricType_NIC
		nics.Timestamp = now.Unix()

		parsing.ParseIfconfigOutput(&nics, st.Stdout)
		return &nics, nil
	}
}
