package server

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
	"github.com/gogo/protobuf/proto"

	"github.com/zperf/tcpmon/tcpmon/parsing"
	"github.com/zperf/tcpmon/tcpmon/tproto"
)

// SocketMonitor collect sockets statistics
type SocketMonitor struct {
	config *CmdConfig
}

func (m *SocketMonitor) Collect(now time.Time) ([]byte, error) {
	r, err := m.doCollect(now)
	if err != nil {
		return nil, err
	}

	metric := &tproto.Metric{Body: &tproto.Metric_Tcp{Tcp: r}}
	val, err := proto.Marshal(metric)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return val, nil
}

func (m *SocketMonitor) doCollect(now time.Time) (*tproto.TcpMetric, error) {
	c := cmd.NewCmd(m.config.PathSS, m.config.ArgSS)
	ctx, cancel := context.WithTimeout(context.Background(), m.config.Timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "ss timeout")

	case st := <-c.Start():
		var t tproto.TcpMetric
		t.Timestamp = now.Unix()
		t.Type = tproto.MetricType_TCP

		parsing.ParseSS(&t, st.Stdout)
		return &t, nil
	}
}
