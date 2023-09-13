package tcpmon

import (
	"context"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"

	"github.com/zperf/tcpmon/tcpmon/tutils"
)

func ParseIfconfigOutput(nics *NicMetric, out []string) {
	r := &IfaceMetric{}
	for _, line := range out {
		if strings.Contains(line, ": flags=") {
			fields := strings.FieldsFunc(line, func(c rune) bool {
				return c == ':'
			})
			r = &IfaceMetric{}
			r.Name = fields[0]
		} else if strings.Contains(line, "RX errors ") {
			fields := strings.FieldsFunc(line, tutils.SplitSpace)
			r.RxErrors, _ = tutils.ParseUint64(fields[2])
			r.RxDropped, _ = tutils.ParseUint64(fields[4])
			r.RxOverruns, _ = tutils.ParseUint64(fields[6])
			r.RxFrame, _ = tutils.ParseUint64(fields[8])
		} else if strings.Contains(line, "TX errors ") {
			fields := strings.FieldsFunc(line, tutils.SplitSpace)
			r.TxErrors, _ = tutils.ParseUint64(fields[2])
			r.TxDropped, _ = tutils.ParseUint64(fields[4])
			r.TxOverruns, _ = tutils.ParseUint64(fields[6])
			r.TxCarrier, _ = tutils.ParseUint64(fields[8])
			r.TxCollisions, _ = tutils.ParseUint64(fields[10])
			nics.Ifaces = append(nics.Ifaces, r)
		}
	}
}

func (m *NicMonitor) RunIfconfig(now time.Time) (*NicMetric, string, error) {
	c := cmd.NewCmd(m.config.PathIfconfig)
	ctx, cancel := context.WithTimeout(context.Background(), m.config.Timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, "", errors.Wrap(ctx.Err(), "ifconfig timeout")
	case st := <-c.Start():
		var nics NicMetric
		nics.Type = MetricType_NIC
		nics.Timestamp = now.Unix()

		ParseIfconfigOutput(&nics, st.Stdout)
		return &nics, strings.Join(st.Stdout, ""), nil
	}
}
