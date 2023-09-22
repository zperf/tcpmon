package tcpmon

import (
	"context"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
	"github.com/rs/zerolog/log"
	"github.com/umisama/go-regexpcache"

	"github.com/zperf/tcpmon/tcpmon/tutils"
)

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

func setRate(m *SocketMetric, field string, value float64) {
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
		m.SndWscale, _ = tutils.ParseUint32(valueStr[:q])
		m.RcvWscale, _ = tutils.ParseUint32(valueStr[q+1:])
		return
	case "rto":
		m.Rto, _ = tutils.ParseFloat64(valueStr)
		return
	case "rtt":
		q := strings.IndexRune(valueStr, '/')
		if q == -1 {
			log.Fatal().Str("field", valueStr).Msg("invalid rtt/rttvar")
		}
		m.Rtt, _ = tutils.ParseFloat64(valueStr[:q])
		m.Rttvar, _ = tutils.ParseFloat64(valueStr[q+1:])
		return
	case "minrtt":
		m.Minrtt, _ = tutils.ParseFloat64(valueStr)
		return
	case "busy":
		m.BusyMs, _ = tutils.ParseUint32(strings.TrimSuffix(valueStr, "ms"))
		return
	case "rcv_rtt":
		m.RcvRtt, _ = tutils.ParseFloat64(valueStr)
		return
	case "retrans":
		q := strings.IndexRune(valueStr, '/')
		if q == -1 {
			log.Fatal().Str("field", valueStr).Msg("invalid retrans")
		}
		m.RetransNow, _ = tutils.ParseUint32(valueStr[:q])
		m.RetransTotal, _ = tutils.ParseUint32(valueStr[q+1:])
		return
	case "rwnd_limited":
		// these two fields(rwnd_limited and sndbuf_limited) must be in ms, check the source of iproute2
		// https://www.mail-archive.com/netdev@vger.kernel.org/msg140890.html
		m.RwndLimited, _ = tutils.ParseUint32(getFirstNumberFromMess(valueStr))
		return
	case "sndbuf_limited":
		m.SndbufLimited, _ = tutils.ParseUint32(getFirstNumberFromMess(valueStr))
		return
	case "ato":
		m.Ato, _ = tutils.ParseFloat64(valueStr)
		return
	case "bytes_acked":
		m.BytesAcked, _ = tutils.ParseUint64(valueStr)
		return
	case "bytes_received":
		m.BytesReceived, _ = tutils.ParseUint64(valueStr)
		return
	}

	value, err := tutils.ParseUint32(valueStr)
	if err != nil {
		log.Warn().Str("value", valueStr).Str("key", key).Err(errors.WithStack(err)).Msg("parse failed")
		return
	}
	switch key {
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

func parseInfos(m *SocketMetric, s string) {
	p := strings.Index(s, ":(")
	if p == -1 {
		log.Fatal().Str("field", s).Msg("parse failed")
	}

	name := s[:p]
	if name == "skmem" {
		fields := strings.FieldsFunc(s[p+2:len(s)-1], func(r rune) bool {
			return ',' == r
		})
		skmem := SocketMemoryUsage{}
		skmem.RmemAlloc, _ = tutils.ParseUint32(strings.TrimPrefix(fields[0], "r"))
		skmem.RcvBuf, _ = tutils.ParseUint32(strings.TrimPrefix(fields[1], "rb"))
		skmem.WmemAlloc, _ = tutils.ParseUint32(strings.TrimPrefix(fields[2], "t"))
		skmem.SndBuf, _ = tutils.ParseUint32(strings.TrimPrefix(fields[3], "tb"))
		skmem.FwdAlloc, _ = tutils.ParseUint32(strings.TrimPrefix(fields[4], "f"))
		skmem.WmemQueued, _ = tutils.ParseUint32(strings.TrimPrefix(fields[5], "w"))
		skmem.OptMem, _ = tutils.ParseUint32(strings.TrimPrefix(fields[6], "o"))
		skmem.BackLog, _ = tutils.ParseUint32(strings.TrimPrefix(fields[7], "bl"))
		if len(fields) > 8 {
			skmem.SockDrop, _ = tutils.ParseUint32(strings.TrimPrefix(fields[8], "d"))
		}
		m.Skmem = &skmem
	} else if name == "timer" {
		fields := strings.FieldsFunc(s[p+2:len(s)-1], func(r rune) bool {
			return ',' == r
		})
		t := &TimerInfo{}
		t.Name = fields[0]
		if len(fields) == 3 {
			if strings.Contains(fields[1], "min") && strings.HasSuffix(fields[1], "sec") {
				ExpireTime := strings.Split(strings.TrimSuffix(fields[1], "sec"), "min")
				ExpireTimeMin, _ := tutils.ParseUint64(ExpireTime[0])
				ExpireTimeSec, _ := tutils.ParseUint64(ExpireTime[1])
				t.ExpireTimeUs = ExpireTimeMin*60000000 + ExpireTimeSec*1000000
			} else if strings.HasSuffix(fields[1], "min") {
				ExpireTimeMin, _ := tutils.ParseUint64(strings.TrimSuffix(fields[1], "min"))
				t.ExpireTimeUs = ExpireTimeMin * 60000000
			} else if strings.HasSuffix(fields[1], "sec") {
				ExpireTimeSec, _ := tutils.ParseUint64(strings.TrimSuffix(fields[1], "sec"))
				t.ExpireTimeUs = ExpireTimeSec * 1000000
			} else if strings.HasSuffix(fields[1], "ms") {
				ExpireTimeMillisecond, _ := tutils.ParseFloat64(strings.TrimSuffix(fields[1], "ms"))
				t.ExpireTimeUs = uint64(ExpireTimeMillisecond * 1000)
			}
			t.Retrans, _ = tutils.ParseUint32(fields[2])
		}
		m.Timers = append(m.Timers, t)
	} else if name == "users" {
		fields := strings.Split(s[p+3:len(s)-2], "),(")
		for _, field := range fields {
			p := &ProcessInfo{}
			f := strings.Split(field, ",")
			p.Name = strings.Trim(f[0], "\"")
			p.Pid, _ = tutils.ParseUint32(strings.TrimPrefix(f[1], "pid="))
			p.Fd, _ = tutils.ParseUint32(strings.TrimPrefix(f[2], "fd="))
			m.Processes = append(m.Processes, p)
		}
	}
}

func ParseSSOutput(t *TcpMetric, out []string) {
	if len(out) == 0 {
		log.Fatal().Msg("Command 'ss' outputs empty")
		return
	}

	header := out[0]
	if regexpcache.MustCompile(`State\s+Recv-Q\s+Send-Q`).Match([]byte(header)) {
		out = out[1:]
	}

	s := &SocketMetric{}
	for _, line := range out {
		fields := strings.FieldsFunc(line, tutils.SplitSpace)

		var exist bool
		if len(fields) == 0 {
			exist = false
		} else {
			_, exist = socketStateMap[fields[0]]
		}

		if exist {
			s = &SocketMetric{}
			s.State = ToPbState(fields[0])
			m, _ := tutils.ParseUint32(fields[1])
			s.RecvQ = m
			n, _ := tutils.ParseInt64(fields[2])
			s.SendQ = n
			s.LocalAddr = fields[3]
			s.PeerAddr = fields[4]
			for _, field := range fields[5:] {
				if strings.Contains(field, ":(") {
					// users and timer
					parseInfos(s, field)
				}
			}
		} else {
			var lastRateName string
			for _, field := range fields {
				switch field {
				case "ts":
					s.Ts = true
				case "sack":
					s.Sack = true
				case "cubic":
					s.Cubic = true
				case "app_limited":
					s.AppLimited = true
				case "ecn":
					s.Ecn = true
				case "ecnseen":
					s.Ecnseen = true
				default:
					// rate handling: pacing_rate, delivery_rate and send
					if isRate(field) {
						lastRateName = field
					} else if lastRateName != "" && strings.HasSuffix(strings.ToLower(field), "bps") {
						field = strings.ToLower(field)
						field = strings.TrimSuffix(field, "bps")

						var rate float64
						carry := 1000.0
						if strings.HasSuffix(strings.ToLower(field), "i") {
							carry = 1024.0
							field = strings.TrimSuffix(field, "i")
						}

						// Base in Kbps or KiBps
						if strings.HasSuffix(field, "g") {
							rateG, _ := tutils.ParseFloat64(strings.TrimSuffix(field, "g"))
							rate = rateG * carry * carry
						} else if strings.HasSuffix(field, "m") {
							rateM, _ := tutils.ParseFloat64(strings.TrimSuffix(field, "m"))
							rate = rateM * carry
						} else if strings.HasSuffix(field, "k") {
							rate, _ = tutils.ParseFloat64(strings.TrimSuffix(field, "k"))
						} else {
							rate, _ = tutils.ParseFloat64(field)
							rate /= carry
						}
						setRate(s, lastRateName, rate)
					} else if strings.Contains(field, ":(") {
						// skmem
						parseInfos(s, field)
					} else {
						// other metrics
						setMetric(s, field)
					}
				}
			}
			t.Sockets = append(t.Sockets, s)
		}
	}
}

func (m *SocketMonitor) RunSS(now time.Time) (*TcpMetric, string, error) {
	c := cmd.NewCmd(m.config.PathSS, m.config.ArgSS)
	ctx, cancel := context.WithTimeout(context.Background(), m.config.Timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, "", errors.Wrap(ctx.Err(), "ss timeout")
	case st := <-c.Start():
		var t TcpMetric
		t.Timestamp = now.Unix()
		t.Type = MetricType_TCP
		ParseSSOutput(&t, st.Stdout)
		return &t, strings.Join(st.Stdout, "\n"), nil
	}
}

func getFirstNumberFromMess(s string) string {
	return regexpcache.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`).FindString(s)
}
