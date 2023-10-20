package export

import "github.com/zperf/tcpmon/tcpmon/gproto"

type Exporter interface {
	ExportMetric(metric *gproto.Metric)
}
