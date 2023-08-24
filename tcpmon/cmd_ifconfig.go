package tcpmon

import (
	"context"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
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
			fields := strings.FieldsFunc(line, SplitSpace)
			r.RxErrors, _ = ParseUint64(fields[2])
			r.RxDropped, _ = ParseUint64(fields[4])
			r.RxOverruns, _ = ParseUint64(fields[6])
			r.RxFrame, _ = ParseUint64(fields[8])
		} else if strings.Contains(line, "TX errors ") {
			fields := strings.FieldsFunc(line, SplitSpace)
			r.TxErrors, _ = ParseUint64(fields[2])
			r.TxDropped, _ = ParseUint64(fields[4])
			r.TxOverruns, _ = ParseUint64(fields[6])
			r.TxCarrier, _ = ParseUint64(fields[8])
			r.TxCollisions, _ = ParseUint64(fields[10])
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
