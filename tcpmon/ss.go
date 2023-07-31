package tcpmon

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func ss(now time.Time) (*TcpMetric, string, error) {
	c := cmd.NewCmd("/usr/bin/ss", "-4ntipmoHOna")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, "", errors.Wrap(ctx.Err(), "ss timeout")
	case st := <-c.Start():
		var t TcpMetric
		t.Timestamp = timestamppb.New(now)
		t.Type = MetricType_Tcp

		var builder strings.Builder
		for _, line := range st.Stdout {
			builder.WriteString(line)
			fields := strings.FieldsFunc(line, func(c rune) bool {
				return c == ' '
			})

			s := SocketMetric{}
			s.State = ToPbState(fields[0])
			n, _ := strconv.ParseUint(fields[1], 10, 32)
			s.RecvQ = uint32(n)
			n, _ = strconv.ParseUint(fields[2], 10, 32)
			s.SendQ = uint32(n)
			s.LocalAddr = fields[3]
			s.PeerAddr = fields[4]

			t.Sockets = append(t.Sockets, &s)
		}

		return &t, builder.String(), nil
	}
}
