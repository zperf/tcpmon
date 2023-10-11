package collector

import (
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gogo/protobuf/proto"

	"github.com/zperf/tcpmon/tcpmon/parsing"
	"github.com/zperf/tcpmon/tcpmon/tproto"
)

const procNet = "/proc/net/"

type NetstatCollector struct{ config *Config }

func NewNetstat(config *Config) *NetstatCollector {
	return &NetstatCollector{config: config}
}

func (m *NetstatCollector) Collect(now time.Time) ([]byte, error) {
	r, err := m.doCollect(now)
	if err != nil {
		return nil, err
	}

	buf, err := proto.Marshal(&tproto.Metric{Body: &tproto.Metric_Net{Net: r}})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return buf, nil
}

func (m *NetstatCollector) doCollect(now time.Time) (*tproto.NetstatMetric, error) {
	var metric tproto.NetstatMetric
	metric.Timestamp = now.Unix()

	err := CollectProc("snmp", &metric)
	if err != nil {
		return nil, err
	}

	err = CollectProc("netstat", &metric)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func CollectProc(t string, metric *tproto.NetstatMetric) error {
	path := procNet + t

	fd, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "open %s failed", path)
	}
	defer fd.Close()

	switch t {
	case "snmp":
		return parsing.ParseSnmp(fd, metric)
	case "netstat":
		return parsing.ParseNetstat(fd, metric)
	default:
		return errors.Newf("unrecognized procfs type: %s", t)
	}
}
