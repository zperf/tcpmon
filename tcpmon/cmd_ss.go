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

func ParseSSOutput(t *TcpMetric, out []string) {
	for _, line := range out {
		fields := strings.FieldsFunc(line, func(c rune) bool {
			return c == ' '
		})

		s := SocketMetric{}
		s.State = ToPbState(fields[0])
		n, _ := parseUint32(fields[1])
		s.RecvQ = n
		n, _ = parseUint32(fields[2])
		s.SendQ = n
		s.LocalAddr = fields[3]
		s.PeerAddr = fields[4]

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
