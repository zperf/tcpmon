package export

import (
	"github.com/zperf/tcpmon/tcpmon/tproto"
)

type Exporter interface {
	ExportMetric(metric *tproto.Metric)
}
