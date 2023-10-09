package parsing

import (
	"strings"

	"github.com/zperf/tcpmon/tcpmon/tproto"
	"github.com/zperf/tcpmon/tcpmon/tutils"
)

func ParseIfconfigOutput(nics *tproto.NicMetric, out []string) {
	r := &tproto.IfaceMetric{}
	for _, line := range out {
		if strings.Contains(line, ": flags=") {
			fields := strings.FieldsFunc(line, func(c rune) bool {
				return c == ':'
			})
			r = &tproto.IfaceMetric{}
			r.Name = fields[0]
		} else if strings.Contains(line, "RX errors ") {
			fields := strings.Fields(line)
			r.RxErrors, _ = tutils.ParseUint64(fields[2])
			r.RxDropped, _ = tutils.ParseUint64(fields[4])
			r.RxOverruns, _ = tutils.ParseUint64(fields[6])
			r.RxFrame, _ = tutils.ParseUint64(fields[8])
		} else if strings.Contains(line, "TX errors ") {
			fields := strings.Fields(line)
			r.TxErrors, _ = tutils.ParseUint64(fields[2])
			r.TxDropped, _ = tutils.ParseUint64(fields[4])
			r.TxOverruns, _ = tutils.ParseUint64(fields[6])
			r.TxCarrier, _ = tutils.ParseUint64(fields[8])
			r.TxCollisions, _ = tutils.ParseUint64(fields[10])
			nics.Ifaces = append(nics.Ifaces, r)
		}
	}
}
