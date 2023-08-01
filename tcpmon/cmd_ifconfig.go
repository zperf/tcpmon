package tcpmon

import (
	"context"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ifconfigOption struct {
	Path    string
	Timeout time.Duration
}

var ifconfigOptions *ifconfigOption

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
			fields := strings.FieldsFunc(line, func(c rune) bool {
				return c == ' '
			})
			r.RxErrors, _ = parseUint32(fields[2])
			r.RxDropped, _ = parseUint32(fields[4])
			r.RxOverruns, _ = parseUint32(fields[6])
			r.RxFrame, _ = parseUint32(fields[8])
		} else if strings.Contains(line, "TX errors ") {
			fields := strings.FieldsFunc(line, func(c rune) bool {
				return c == ' '
			})
			r.TxErrors, _ = parseUint32(fields[2])
			r.TxDropped, _ = parseUint32(fields[4])
			r.TxOverruns, _ = parseUint32(fields[6])
			r.TxCarrier, _ = parseUint32(fields[8])
			r.TxCollisions, _ = parseUint32(fields[10])
			nics.Ifaces = append(nics.Ifaces, r)
		}
	}
}

func RunIfconfig(now time.Time) (*NicMetric, string, error) {
	if ifconfigOptions == nil {
		ifconfigOptions = &ifconfigOption{
			Path:    viper.GetString("ifconfig"),
			Timeout: viper.GetDuration("command-timeout"),
		}
	}

	c := cmd.NewCmd(ifconfigOptions.Path)
	ctx, cancel := context.WithTimeout(context.Background(), ifconfigOptions.Timeout)
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
