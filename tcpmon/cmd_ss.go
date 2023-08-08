package tcpmon

import (
	"context"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ssOption struct {
	Path    string
	Timeout time.Duration
	Arg     string
}

var ssOptions *ssOption
var socketStateMap map[string]SocketState

func init() {
	socketStateMap = map[string]SocketState{
		"ESTAB":      SocketState_TCP_ESTABLISHED,
		"SYN-SENT":   SocketState_TCP_SYN_SENT,
		"SYN-RECV":   SocketState_TCP_SYN_RECV,
		"FIN-WAIT-1": SocketState_TCP_FIN_WAIT1,
		"FIN-WAIT-2": SocketState_TCP_FIN_WAIT2,
		"TIME-WAIT":  SocketState_TCP_TIME_WAIT,
		"UNCONN":     SocketState_TCP_CLOSE,
		"CLOSE-WAIT": SocketState_TCP_CLOSE_WAIT,
		"LAST-ACK":   SocketState_TCP_LAST_ACK,
		"LISTEN":     SocketState_TCP_LISTEN,
		"CLOSING":    SocketState_TCP_CLOSING,
	}
}

// ToPbState converts string to pb enum
// From https://sourcegraph.com/github.com/shemminger/iproute2/-/blob/misc/ss.c?L1397
func ToPbState(s string) SocketState {
	st, ok := socketStateMap[s]
	if !ok {
		log.Fatal().Err(errors.Newf("unknown socket state %v", s)).Msg("failed to convert str to socket state")
	}
	return st
}

func isRate(s string) bool {
	switch s {
	case "pacing_rate":
		return true
	case "delivery_rate":
		return true
	case "send":
		return true
	default:
		return false
	}
}

func setRate(m *SocketMetric, field string, value uint64) {
	switch field {
	case "pacing_rate":
		m.PacingRate = value
	case "delivery_rate":
		m.DeliveryRate = value
	case "send":
		m.Send = value
	default:
		log.Fatal().Str("field", field).Msg("invalid field")
	}
}

func setMetric(m *SocketMetric, field string) {
	p := strings.IndexRune(field, ':')
	if p == -1 {
		log.Fatal().Str("field", field).Msg("invalid field")
	}
	key := field[:p]
	valueStr := field[p+1:]

	switch key {
	case "wscale":
		q := strings.IndexRune(valueStr, ',')
		if q == -1 {
			log.Fatal().Str("field", valueStr).Msg("invalid wscale")
		}
		m.SndWscale, _ = ParseUint32(valueStr[:q])
		m.RcvWscale, _ = ParseUint32(valueStr[q+1:])
		return
	case "rto":
		m.Rto, _ = ParseFloat64(valueStr)
		return
	case "rtt":
		q := strings.IndexRune(valueStr, '/')
		if q == -1 {
			log.Fatal().Str("field", valueStr).Msg("invalid rtt/rttvar")
		}
		m.Rtt, _ = ParseFloat64(valueStr[:q])
		m.Rttvar, _ = ParseFloat64(valueStr[q+1:])
		return
	case "minrtt":
		m.Minrtt, _ = ParseFloat64(valueStr)
		return
	case "busy":
		m.BusyMs, _ = ParseUint32(strings.TrimSuffix(valueStr, "ms"))
		return
	case "rcv_rtt":
		m.RcvRtt, _ = ParseFloat64(valueStr)
		return
	case "retrans":
		q := strings.IndexRune(valueStr, '/')
		if q == -1 {
			log.Fatal().Str("field", valueStr).Msg("invalid retrans")
		}
		m.RetransNow, _ = ParseUint32(valueStr[:q])
		m.RetransTotal, _ = ParseUint32(valueStr[q+1:])
		return
	}

	value, err := ParseUint32(valueStr)
	if err != nil {
		log.Fatal().Str("value", valueStr).Str("key", key).Err(errors.WithStack(err)).Msg("parse failed")
	}
	switch key {
	case "ato":
		m.Ato = value
	case "mss":
		m.Mss = value
	case "pmtu":
		m.Pmtu = value
	case "rcvmss":
		m.Rcvmss = value
	case "advmss":
		m.Advmss = value
	case "cwnd":
		m.Cwnd = value
	case "bytes_sent":
		m.BytesSent = value
	case "bytes_acked":
		m.BytesAcked = value
	case "bytes_received":
		m.BytesReceived = value
	case "data_segs_out":
		m.DataSegsOut = value
	case "data_segs_in":
		m.DataSegsIn = value
	case "segs_out":
		m.SegsOut = value
	case "segs_in":
		m.SegsIn = value
	case "lastsnd":
		m.Lastsnd = value
	case "lastrcv":
		m.Lastrcv = value
	case "lastack":
		m.Lastack = value
	case "delivered":
		m.Delivered = value
	case "rcv_space":
		m.RcvSpace = value
	case "rcv_ssthresh":
		m.RcvSsthresh = value
	case "snd_wnd":
		m.SndWnd = value
	}
}

func parseTimerOrMeminfo(m *SocketMetric, s string) {
	p := strings.Index(s, ":(")
	if p == -1 {
		log.Fatal().Str("field", s).Msg("parse failed")
	}

	name := s[:p]
	fields := strings.FieldsFunc(s[p+2:len(s)-1], func(r rune) bool {
		return ',' == r
	})
	if name == "skmem" {
		skmem := SocketMemoryUsage{}
		skmem.RmemAlloc, _ = ParseUint32(strings.TrimPrefix(fields[0], "r"))
		skmem.RcvBuf, _ = ParseUint32(strings.TrimPrefix(fields[1], "rb"))
		skmem.WmemAlloc, _ = ParseUint32(strings.TrimPrefix(fields[2], "t"))
		skmem.SndBuf, _ = ParseUint32(strings.TrimPrefix(fields[3], "tb"))
		skmem.FwdAlloc, _ = ParseUint32(strings.TrimPrefix(fields[4], "f"))
		skmem.WmemQueued, _ = ParseUint32(strings.TrimPrefix(fields[5], "w"))
		skmem.OptMem, _ = ParseUint32(strings.TrimPrefix(fields[6], "o"))
		skmem.BackLog, _ = ParseUint32(strings.TrimPrefix(fields[7], "bl"))
		skmem.SockDrop, _ = ParseUint32(strings.TrimPrefix(fields[8], "d"))
		m.Skmem = &skmem
	} else if name == "timer" {
		t := TimerInfo{}
		t.Name = fields[0]
		if len(fields) == 3 {
			t.ExpireTimeSec, _ = ParseUint32(strings.TrimSuffix(fields[1], "sec"))
			t.Retrans, _ = ParseUint32(fields[2])
		}
		m.Timers = append(m.Timers, &t)
	}
}

func ParseSSOutput(t *TcpMetric, out []string) {
	for _, line := range out {
		fields := strings.FieldsFunc(line, func(c rune) bool {
			return c == ' '
		})

		s := SocketMetric{}
		s.State = ToPbState(fields[0])
		n, _ := ParseUint32(fields[1])
		s.RecvQ = n
		n, _ = ParseUint32(fields[2])
		s.SendQ = n
		s.LocalAddr = fields[3]
		s.PeerAddr = fields[4]

		var lastRateName string
		for _, field := range fields[5:] {
			switch field {
			case "ts":
				s.Ts = true
			case "sack":
				s.Sack = true
			case "cubic":
				s.Cubic = true
			case "app_limited":
				s.AppLimited = true
			default:
				// rate handling: pacing_rate, delivery_rate and send
				if isRate(field) {
					lastRateName = field
				} else if lastRateName != "" && strings.HasSuffix(field, "bps") {
					rate, _ := ParseUint64(strings.TrimSuffix(field, "bps"))
					setRate(&s, lastRateName, rate)
				} else if strings.Contains(field, ":(") {
					// timers
					parseTimerOrMeminfo(&s, field)
				} else {
					// other metrics
					setMetric(&s, field)
				}
			}
		}

		t.Sockets = append(t.Sockets, &s)
	}
}

func RunSS(now time.Time) (*TcpMetric, string, error) {
	if ssOptions == nil {
		ssOptions = &ssOption{
			Path:    viper.GetString("ss"),
			Timeout: viper.GetDuration("command-timeout"),
			Arg:     viper.GetString("ss-arg"),
		}
	}

	c := cmd.NewCmd(ssOptions.Path, ssOptions.Arg)
	ctx, cancel := context.WithTimeout(context.Background(), ssOptions.Timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, "", errors.Wrap(ctx.Err(), "ss timeout")
	case st := <-c.Start():
		var t TcpMetric
		t.Timestamp = timestamppb.New(now)
		t.Type = MetricType_TCP
		ParseSSOutput(&t, st.Stdout)
		return &t, strings.Join(st.Stdout, "\n"), nil
	}
}
