package tcpmon

import (
	"context"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"

	"github.com/zperf/tcpmon/tcpmon/tutils"
)

var headSet = map[string]struct{}{
	"Ip:":      {},
	"Icmp:":    {},
	"IcmpMsg:": {},
	"Tcp:":     {},
	"Udp:":     {},
	"UdpLite:": {},
	"TcpExt:":  {},
	"IpExt:":   {},
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
				r.IpTotalPacketsReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "forwarded") {
				r.IpForwarded, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "incoming packets discarded") {
				r.IpIncomingPacketsDiscarded, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "incoming packets delivered") {
				r.IpIncomingPacketsDelivered, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "requests sent out") {
				r.IpRequestsSentOut, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "outgoing packets dropped") {
				r.IpOutgoingPacketsDropped, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			}
		} else if flag == "Tcp:" {
			if strings.Contains(line, "active connections openings") {
				r.TcpActiveConnectionsOpenings, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "passive connection openings") {
				r.TcpPassiveConnectionOpenings, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "failed connection attempts") {
				r.TcpFailedConnectionAttempts, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "connection resets received") {
				r.TcpConnectionResetsReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "connections established") {
				r.TcpConnectionsEstablished, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "segments received") && !strings.Contains(line, "bad") {
				r.TcpSegmentsReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "segments send out") {
				r.TcpSegmentsSendOut, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "segments retransmit") {
				// NOTE: in the newer Linux version (like fc38), `netstat -s | grep retrans` will return retransmitted
				// The older (like el7) will return a typo: retransmited
				// We should take the common prefix `retransmit`
				r.TcpSegmentsRetransmitted, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "bad segments received") {
				r.TcpBadSegmentsReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "resets sent") {
				r.TcpResetsSent, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			}
		} else if flag == "Udp:" {
			if strings.Contains(line, "packets received") {
				r.UdpPacketsReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "packets to unknown port received") {
				r.UdpPacketsToUnknownPortReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "packet receive errors") {
				r.UdpPacketReceiveErrors, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "packets sent") {
				r.UdpPacketsSent, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "receive buffer errors") {
				r.UdpReceiveBufferErrors, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "send buffer errors") {
				r.UdpSendBufferErrors, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			}
		}
	}
}

func (m *NetstatMonitor) RunNetstat(now time.Time) (*NetstatMetric, string, error) {
	c := cmd.NewCmd(m.config.PathNetstat, m.config.ArgNetstat)
	ctx, cancel := context.WithTimeout(context.Background(), m.config.Timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, "", errors.Wrap(ctx.Err(), "netstat timeout")
	case st := <-c.Start():
		var r NetstatMetric
		r.Type = MetricType_NET
		r.Timestamp = now.Unix()

		ParseNetstatOutput(&r, st.Stdout)
		return &r, strings.Join(st.Stdout, "\n"), nil
	}
}
