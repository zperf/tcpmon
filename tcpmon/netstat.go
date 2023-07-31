package tcpmon

import (
	"context"
	"strconv"
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

type NetstatRecord struct {
	ip_total_packets_received     int
	ip_forwarded                  int
	ip_incoming_packets_discarded int
	ip_incoming_packets_delivered int
	ip_requests_sent_out          int
	ip_outgoing_packets_dropped   int

	tcp_active_connections_openings int
	tcp_passive_connection_openings int
	tcp_failed_connection_attempts  int
	tcp_connection_resets_received  int
	tcp_connections_established     int
	tcp_segments_received           int
	tcp_segments_send_out           int
	tcp_segments_retransmited       int
	tcp_bad_segments_received       int
	tcp_resets_sent                 int

	udp_packets_received                 int
	udp_packets_to_unknown_port_received int
	udp_packet_receive_errors            int
	udp_packets_sent                     int
	udp_receive_buffer_errors            int
	udp_send_buffer_errors               int
}

func netstat() (*NetstatRecord, string, error) {
	cmd := cmd.NewCmd("/usr/bin/netstat", "-s")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, "", errors.Wrap(ctx.Err(), "netstat timeout")
	case st := <-cmd.Start():
		var builder strings.Builder

		r := NetstatRecord{}
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
					r.ip_total_packets_received, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "forwarded") {
					r.ip_forwarded, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "incoming packets discarded") {
					r.ip_incoming_packets_discarded, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "incoming packets delivered") {
					r.ip_incoming_packets_delivered, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "requests sent out") {
					r.ip_requests_sent_out, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "outgoing packets dropped") {
					r.ip_outgoing_packets_dropped, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				}
			} else if flag == "Tcp:" {
				if strings.Contains(line, "active connections openings") {
					r.tcp_active_connections_openings, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "passive connection openings") {
					r.tcp_passive_connection_openings, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "failed connection attempts") {
					r.tcp_failed_connection_attempts, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "connection resets received") {
					r.tcp_connection_resets_received, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "connections established") {
					r.tcp_connections_established, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "segments received") && !strings.Contains(line, "bad") {
					r.tcp_segments_received, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "segments send out") {
					r.tcp_segments_send_out, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "segments retransmited") {
					r.tcp_segments_retransmited, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "bad segments received") {
					r.tcp_bad_segments_received, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "resets sent") {
					r.tcp_resets_sent, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				}
			} else if flag == "Udp:" {
				if strings.Contains(line, "packets received") {
					r.udp_packets_received, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "packets to unknown port received") {
					r.udp_packets_to_unknown_port_received, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "packet receive errors") {
					r.udp_packet_receive_errors, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "packets sent") {
					r.udp_packets_sent, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "receive buffer errors") {
					r.udp_receive_buffer_errors, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				} else if strings.Contains(line, "send buffer errors") {
					r.udp_send_buffer_errors, _ = strconv.Atoi(strings.FieldsFunc(line, split_func)[0])
				}
			}
		}
		return &r, builder.String(), nil
	}
}
