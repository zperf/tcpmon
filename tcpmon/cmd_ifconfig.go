package tcpmon

import (
	"context"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			fields := strings.FieldsFunc(line, splitSpace)
			r.RxErrors, _ = ParseUint32(fields[2])
			r.RxDropped, _ = ParseUint32(fields[4])
			r.RxOverruns, _ = ParseUint32(fields[6])
			r.RxFrame, _ = ParseUint32(fields[8])
		} else if strings.Contains(line, "TX errors ") {
			fields := strings.FieldsFunc(line, splitSpace)
			r.TxErrors, _ = ParseUint32(fields[2])
			r.TxDropped, _ = ParseUint32(fields[4])
			r.TxOverruns, _ = ParseUint32(fields[6])
			r.TxCarrier, _ = ParseUint32(fields[8])
			r.TxCollisions, _ = ParseUint32(fields[10])
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
		nics.Timestamp = timestamppb.New(now)

		ParseIfconfigOutput(&nics, st.Stdout)
		return &nics, strings.Join(st.Stdout, ""), nil
	}
}
