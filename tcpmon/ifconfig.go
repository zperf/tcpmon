package tcpmon

import (
	"context"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func RunIfconfig(now time.Time) (*NicMetric, string, error) {
	cmd := cmd.NewCmd("/usr/sbin/ifconfig")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, "", errors.Wrap(ctx.Err(), "ifconfig timeout")
	case st := <-cmd.Start():
		var builder strings.Builder

		var nics NicMetric
		nics.Type = MetricType_Nic
		nics.Timestamp = timestamppb.New(now)

		var r IfaceMetric
		for _, line := range st.Stdout {
			builder.WriteString(line)
			if strings.Contains(line, ": flags=") {
				fields := strings.FieldsFunc(line, func(c rune) bool {
					return c == ':'
				})
				r.Name = fields[0]
			} else if strings.Contains(line, "RX errors ") {
				fields := strings.FieldsFunc(line, func(c rune) bool {
					return c == ' '
				})
				r.RXErrors, _ = parseUint32(fields[2])
				r.RXDropped, _ = parseUint32(fields[4])
				r.RXOverruns, _ = parseUint32(fields[6])
				r.RXFrame, _ = parseUint32(fields[8])
			} else if strings.Contains(line, "TX errors ") {
				fields := strings.FieldsFunc(line, func(c rune) bool {
					return c == ' '
				})
				r.TXErrors, _ = parseUint32(fields[2])
				r.TXDropped, _ = parseUint32(fields[4])
				r.TXOverruns, _ = parseUint32(fields[6])
				r.TXCarrier, _ = parseUint32(fields[8])
				r.TXCollisions, _ = parseUint32(fields[10])
				nics.Ifaces = append(nics.Ifaces, &r)
				r = IfaceMetric{}
			}
		}
		return &nics, builder.String(), nil
	}
}
