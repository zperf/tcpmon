package tcpmon

import (
	"context"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/spf13/viper"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
)

type netstatOption struct {
	Path    string
	Timeout time.Duration
	Arg     string
}

var netstatOptions *netstatOption

var headSet = map[string]struct{}{
	"Ip:":      struct{}{},
	"Icmp:":    struct{}{},
	"IcmpMsg:": struct{}{},
	"Tcp:":     struct{}{},
	"Udp:":     struct{}{},
	"UdpLite:": struct{}{},
	"TcpExt:":  struct{}{},
	"IpExt:":   struct{}{},
}

var splitFunc = func(c rune) bool {
	return c == ' '
}

func ParseNetstatOutput(r *NetstatMetric, out []string) {
	flag := ""
	for _, line := range out {
		if _, exist := headSet[line]; exist {
			flag = line
			if flag == "UdpLite:" {
				break
			}
			continue
		}
		if flag == "Ip:" {
			if strings.Contains(line, "total packets received") {
				r.IpTotalPacketsReceived, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "forwarded") {
				r.IpForwarded, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "incoming packets discarded") {
				r.IpIncomingPacketsDiscarded, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "incoming packets delivered") {
				r.IpIncomingPacketsDelivered, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "requests sent out") {
				r.IpRequestsSentOut, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "outgoing packets dropped") {
				r.IpOutgoingPacketsDropped, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			}
		} else if flag == "Tcp:" {
			if strings.Contains(line, "active connections openings") {
				r.TcpActiveConnectionsOpenings, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "passive connection openings") {
				r.TcpPassiveConnectionOpenings, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "failed connection attempts") {
				r.TcpFailedConnectionAttempts, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "connection resets received") {
				r.TcpConnectionResetsReceived, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "connections established") {
				r.TcpConnectionsEstablished, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "segments received") && !strings.Contains(line, "bad") {
				r.TcpSegmentsReceived, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "segments send out") {
				r.TcpSegmentsSendOut, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "segments retransmited") {
				r.TcpSegmentsRetransmitted, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "bad segments received") {
				r.TcpBadSegmentsReceived, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "resets sent") {
				r.TcpResetsSent, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			}
		} else if flag == "Udp:" {
			if strings.Contains(line, "packets received") {
				r.UdpPacketsReceived, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "packets to unknown port received") {
				r.UdpPacketsToUnknownPortReceived, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "packet receive errors") {
				r.UdpPacketReceiveErrors, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "packets sent") {
				r.UdpPacketsSent, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "receive buffer errors") {
				r.UdpReceiveBufferErrors, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			} else if strings.Contains(line, "send buffer errors") {
				r.UdpSendBufferErrors, _ = ParseUint32(strings.FieldsFunc(line, splitFunc)[0])
			}
		}
	}
}

func RunNetstat(now time.Time) (*NetstatMetric, string, error) {
	if netstatOptions == nil {
		netstatOptions = &netstatOption{
			Path:    viper.GetString("netstat"),
			Timeout: viper.GetDuration("command-timeout"),
			Arg:     viper.GetString("netstat-arg"),
		}
	}

	cmd := cmd.NewCmd(netstatOptions.Path, netstatOptions.Arg)
	ctx, cancel := context.WithTimeout(context.Background(), netstatOptions.Timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, "", errors.Wrap(ctx.Err(), "netstat timeout")
	case st := <-cmd.Start():
		var r NetstatMetric
		r.Type = MetricType_NET
		r.Timestamp = timestamppb.New(now)

		ParseNetstatOutput(&r, st.Stdout)
		return &r, strings.Join(st.Stdout, "\n"), nil
	}
}
