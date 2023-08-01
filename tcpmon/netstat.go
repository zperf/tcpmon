package tcpmon

import (
	"context"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
)

var head_set = map[string]bool{
	"Ip:":      true,
	"Icmp:":    true,
	"IcmpMsg:": true,
	"Tcp:":     true,
	"Udp:":     true,
	"UdpLite:": true,
	"TcpExt:":  true,
	"IpExt:":   true,
}

var split_func = func(c rune) bool {
	return c == ' '
}

func RunNetstat() (*NetstatMetric, string, error) {
	cmd := cmd.NewCmd("/usr/bin/netstat", "-s")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, "", errors.Wrap(ctx.Err(), "netstat timeout")
	case st := <-cmd.Start():
		var builder strings.Builder

		var r NetstatMetric
		flag := ""
		for _, line := range st.Stdout {
			builder.WriteString(line)
			if _, exist := head_set[line]; exist {
				flag = line
				if flag == "UdpLite:" {
					break
				}
				continue
			}
			if flag == "Ip:" {
				if strings.Contains(line, "total packets received") {
					r.IpTotalPacketsReceived, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "forwarded") {
					r.IpForwarded, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "incoming packets discarded") {
					r.IpIncomingPacketsDiscarded, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "incoming packets delivered") {
					r.IpIncomingPacketsDelivered, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "requests sent out") {
					r.IpRequestsSentOut, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "outgoing packets dropped") {
					r.IpOutgoingPacketsDropped, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				}
			} else if flag == "Tcp:" {
				if strings.Contains(line, "active connections openings") {
					r.TcpActiveConnectionsOpenings, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "passive connection openings") {
					r.TcpPassiveConnectionOpenings, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "failed connection attempts") {
					r.TcpFailedConnectionAttempts, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "connection resets received") {
					r.TcpConnectionResetsReceived, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "connections established") {
					r.TcpConnectionsEstablished, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "segments received") && !strings.Contains(line, "bad") {
					r.TcpSegmentsReceived, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "segments send out") {
					r.TcpSegmentsSendOut, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "segments retransmited") {
					r.TcpSegmentsRetransmited, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "bad segments received") {
					r.TcpBadSegmentsReceived, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "resets sent") {
					r.TcpResetsSent, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				}
			} else if flag == "Udp:" {
				if strings.Contains(line, "packets received") {
					r.UdpPacketsReceived, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "packets to unknown port received") {
					r.UdpPacketsToUnknownPortReceived, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "packet receive errors") {
					r.UdpPacketReceiveErrors, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "packets sent") {
					r.UdpPacketsSent, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "receive buffer errors") {
					r.UdpReceiveBufferErrors, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "send buffer errors") {
					r.UdpSendBufferErrors, _ = parseUint32(strings.FieldsFunc(line, split_func)[0])
				}
			}
		}
		return &r, builder.String(), nil
	}
}
