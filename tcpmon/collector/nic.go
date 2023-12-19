package collector

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
	"google.golang.org/protobuf/proto"

	"github.com/zperf/tcpmon/tcpmon/gproto"
	"github.com/zperf/tcpmon/tcpmon/parsing"
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

	metric := &gproto.Metric{Body: &gproto.Metric_Nic{Nic: r}}
	val, err := proto.Marshal(metric)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return val, nil
}

func (m *NicCollector) doCollect(now time.Time) (*gproto.NicMetric, error) {
	c := cmd.NewCmd(m.config.PathIfconfig)
	ctx, cancel := context.WithTimeout(context.Background(), m.config.Timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		err := c.Stop()
		return nil, errors.Wrap(errors.CombineErrors(ctx.Err(), err), "ifconfig timeout")
	case st := <-c.Start():
		var nics gproto.NicMetric
		nics.Type = gproto.MetricType_NIC
		nics.Timestamp = now.Unix()

		parsing.ParseIfconfigOutput(&nics, st.Stdout)
		return &nics, nil
	}
}
