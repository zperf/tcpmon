package tcpmon_test

import (
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	. "github.com/zperf/tcpmon/tcpmon"
)

func (s *CommandParserTestSuite) TestParseSS() {
	assert := s.Assert()

	out := `
LISTEN  0       4096     127.0.0.1:53000      0.0.0.0:*       skmem:(r0,rb131072,t0,tb16384,f0,w0,o0,bl0,d0) cubic cwnd:10
LISTEN  0       4096     127.0.0.1:631        0.0.0.0:*       skmem:(r0,rb131072,t0,tb16384,f0,w0,o0,bl0,d0) cubic cwnd:10
ESTAB   0       0        127.0.0.1:41317    127.0.0.1:40534   timer:(keepalive,11sec,0) skmem:(r0,rb131072,t0,tb2626560,f0,w0,o0,bl0,d276) ts sack cubic wscale:7,7 rto:206.666 rtt:4.413/8.337 ato:40 mss:32768 pmtu:65535 rcvmss:536 advmss:65483 cwnd:10 bytes_sent:1257 bytes_acked:1257 bytes_received:3442 segs_out:1112 segs_in:1118 data_segs_out:7 data_segs_in:7 send 594026739bps lastsnd:12721397 lastrcv:12721397 lastack:3530 pacing_rate 1187918880bps delivery_rate 23831272720bps delivered:8 app_limited busy:53ms rcv_space:65483 rcv_ssthresh:65483 minrtt:0.01 snd_wnd:65536
ESTAB   0       0        127.0.0.1:41317    127.0.0.1:40554   timer:(keepalive,11sec,0) skmem:(r0,rb131072,t0,tb2626560,f0,w0,o0,bl0,d276) ts sack cubic wscale:7,7 rto:206.666 rtt:6.562/12.05 ato:40 mss:32768 pmtu:65535 rcvmss:536 advmss:65483 cwnd:10 bytes_sent:312 bytes_acked:312 bytes_received:2099 segs_out:1109 segs_in:1112 data_segs_out:4 data_segs_in:4 send 399487961bps lastsnd:12722864 lastrcv:12722877 lastack:3530 pacing_rate 798975920bps delivery_rate 29127111104bps delivered:5 app_limited busy:53ms rcv_space:65483 rcv_ssthresh:65483 minrtt:0.009 snd_wnd:65536
`

	lines := strings.FieldsFunc(out, func(r rune) bool {
		return r == '\n'
	})

	var t TcpMetric
	t.Timestamp = timestamppb.New(time.Now())
	t.Type = MetricType_TCP
	ParseSSOutput(&t, lines)

	assert.Equal(SocketState_TCP_LISTEN, t.Sockets[0].State)
	assert.Equal(SocketState_TCP_LISTEN, t.Sockets[1].State)
	assert.Equal(SocketState_TCP_ESTABLISHED, t.Sockets[2].State)
	assert.Equal(SocketState_TCP_ESTABLISHED, t.Sockets[3].State)

	assert.Equal(uint32(0), t.Sockets[0].RecvQ)
	assert.Equal(uint32(0), t.Sockets[1].RecvQ)
	assert.Equal(uint32(0), t.Sockets[2].RecvQ)
	assert.Equal(uint32(0), t.Sockets[3].RecvQ)

	assert.Equal(uint32(4096), t.Sockets[0].SendQ)
	assert.Equal(uint32(4096), t.Sockets[1].SendQ)
	assert.Equal(uint32(0), t.Sockets[2].RecvQ)
	assert.Equal(uint32(0), t.Sockets[3].RecvQ)

	assert.Equal("127.0.0.1:53000", t.Sockets[0].LocalAddr)
	assert.Equal("127.0.0.1:631", t.Sockets[1].LocalAddr)
	assert.Equal("127.0.0.1:41317", t.Sockets[2].LocalAddr)
	assert.Equal("127.0.0.1:41317", t.Sockets[3].LocalAddr)

	assert.Equal("0.0.0.0:*", t.Sockets[0].PeerAddr)
	assert.Equal("0.0.0.0:*", t.Sockets[1].PeerAddr)
	assert.Equal("127.0.0.1:40534", t.Sockets[2].PeerAddr)
	assert.Equal("127.0.0.1:40554", t.Sockets[3].PeerAddr)
}
