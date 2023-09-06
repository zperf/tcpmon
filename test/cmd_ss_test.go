package test

import (
	"strings"
	"time"

	. "github.com/zperf/tcpmon/tcpmon"
)

func (s *CommandParserTestSuite) TestParseSS() {
	assert := s.Assert()

	out := `State      Recv-Q Send-Q Local Address:Port               Peer Address:Port
LISTEN  0       4096     127.0.0.1:53000      0.0.0.0:*
skmem:(r0,rb131072,t0,tb16384,f0,w0,o0,bl0,d0) cubic cwnd:10
LISTEN  0       4096     127.0.0.1:631        0.0.0.0:*
skmem:(r0,rb131072,t0,tb16384,f0,w0,o0,bl0,d0) cubic cwnd:10
ESTAB   0       0        127.0.0.1:41317    127.0.0.1:40534
timer:(keepalive,11sec,0) skmem:(r0,rb131072,t0,tb2626560,f0,w0,o0,bl0,d276) ts sack cubic wscale:7,7 rto:206.666 rtt:4.413/8.337 ato:40 mss:32768 pmtu:65535 rcvmss:536 advmss:65483 cwnd:10 bytes_sent:1257 bytes_acked:1257 bytes_received:3442 segs_out:1112 segs_in:1118 data_segs_out:7 data_segs_in:7 send 594026739bps lastsnd:12721397 lastrcv:12721397 lastack:3530 pacing_rate 1187918880bps delivery_rate 23831272720bps delivered:8 app_limited busy:53ms rcv_space:65483 rcv_ssthresh:65483 minrtt:0.01 snd_wnd:65536
ESTAB   0       0        127.0.0.1:41317    127.0.0.1:40554
timer:(keepalive,11sec,0) skmem:(r0,rb131072,t0,tb2626560,f0,w0,o0,bl0,d276) ts sack cubic wscale:7,7 rto:206.666 rtt:6.562/12.05 ato:40 mss:32768 pmtu:65535 rcvmss:536 advmss:65483 cwnd:10 bytes_sent:312 bytes_acked:312 bytes_received:2099 segs_out:1109 segs_in:1112 data_segs_out:4 data_segs_in:4 send 399487961bps lastsnd:12722864 lastrcv:12722877 lastack:3530 pacing_rate 798975920bps delivery_rate 29127111104bps delivered:5 app_limited busy:53ms rcv_space:65483 rcv_ssthresh:65483 minrtt:0.009 snd_wnd:65536
`
	lines := strings.FieldsFunc(out, SplitNewline)

	var t TcpMetric
	t.Timestamp = time.Now().UnixMilli()
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

	assert.Nil(t.Sockets[0].Timers)
	assert.Nil(t.Sockets[1].Timers)

	assert.Equal(uint32(0), t.Sockets[0].Skmem.RmemAlloc)
	assert.Equal(uint32(0), t.Sockets[1].Skmem.RmemAlloc)
	assert.Equal(uint32(0), t.Sockets[2].Skmem.RmemAlloc)
	assert.Equal(uint32(0), t.Sockets[3].Skmem.RmemAlloc)

	assert.Equal(uint32(131072), t.Sockets[0].Skmem.RcvBuf)
	assert.Equal(uint32(131072), t.Sockets[1].Skmem.RcvBuf)
	assert.Equal(uint32(131072), t.Sockets[2].Skmem.RcvBuf)
	assert.Equal(uint32(131072), t.Sockets[3].Skmem.RcvBuf)

	assert.Equal(uint32(0), t.Sockets[0].Skmem.WmemAlloc)
	assert.Equal(uint32(0), t.Sockets[1].Skmem.WmemAlloc)
	assert.Equal(uint32(0), t.Sockets[2].Skmem.WmemAlloc)
	assert.Equal(uint32(0), t.Sockets[3].Skmem.WmemAlloc)

	assert.Equal(uint32(16384), t.Sockets[0].Skmem.SndBuf)
	assert.Equal(uint32(16384), t.Sockets[1].Skmem.SndBuf)
	assert.Equal(uint32(2626560), t.Sockets[2].Skmem.SndBuf)
	assert.Equal(uint32(2626560), t.Sockets[3].Skmem.SndBuf)

	assert.Equal(uint32(0), t.Sockets[0].Skmem.FwdAlloc)
	assert.Equal(uint32(0), t.Sockets[1].Skmem.FwdAlloc)
	assert.Equal(uint32(0), t.Sockets[2].Skmem.FwdAlloc)
	assert.Equal(uint32(0), t.Sockets[3].Skmem.FwdAlloc)

	assert.Equal(uint32(0), t.Sockets[0].Skmem.WmemAlloc)
	assert.Equal(uint32(0), t.Sockets[1].Skmem.WmemAlloc)
	assert.Equal(uint32(0), t.Sockets[2].Skmem.WmemAlloc)
	assert.Equal(uint32(0), t.Sockets[3].Skmem.WmemAlloc)

	assert.Equal(uint32(0), t.Sockets[0].Skmem.WmemQueued)
	assert.Equal(uint32(0), t.Sockets[1].Skmem.WmemQueued)
	assert.Equal(uint32(0), t.Sockets[2].Skmem.WmemQueued)
	assert.Equal(uint32(0), t.Sockets[3].Skmem.WmemQueued)

	assert.Equal(uint32(0), t.Sockets[0].Skmem.OptMem)
	assert.Equal(uint32(0), t.Sockets[1].Skmem.OptMem)
	assert.Equal(uint32(0), t.Sockets[2].Skmem.OptMem)
	assert.Equal(uint32(0), t.Sockets[3].Skmem.OptMem)

	assert.Equal(uint32(0), t.Sockets[0].Skmem.BackLog)
	assert.Equal(uint32(0), t.Sockets[1].Skmem.BackLog)
	assert.Equal(uint32(0), t.Sockets[2].Skmem.BackLog)
	assert.Equal(uint32(0), t.Sockets[3].Skmem.BackLog)

	assert.Equal(uint32(0), t.Sockets[0].Skmem.SockDrop)
	assert.Equal(uint32(0), t.Sockets[1].Skmem.SockDrop)
	assert.Equal(uint32(276), t.Sockets[2].Skmem.SockDrop)
	assert.Equal(uint32(276), t.Sockets[3].Skmem.SockDrop)

	assert.True(t.Sockets[0].Cubic)
	assert.True(t.Sockets[1].Cubic)
	assert.True(t.Sockets[2].Cubic)
	assert.True(t.Sockets[3].Cubic)

	assert.True(t.Sockets[2].AppLimited)
	assert.True(t.Sockets[3].AppLimited)

	assert.Equal(uint32(10), t.Sockets[0].Cwnd)
	assert.Equal(uint32(10), t.Sockets[1].Cwnd)
	assert.Equal(uint32(10), t.Sockets[2].Cwnd)
	assert.Equal(uint32(10), t.Sockets[3].Cwnd)

	assert.Equal(uint32(7), t.Sockets[2].SndWscale)
	assert.Equal(uint32(7), t.Sockets[2].RcvWscale)
	assert.Equal(uint32(7), t.Sockets[3].SndWscale)
	assert.Equal(uint32(7), t.Sockets[3].RcvWscale)

	assert.Equal(206.666, t.Sockets[2].Rto)
	assert.Equal(206.666, t.Sockets[3].Rto)

	assert.Equal(4.413, t.Sockets[2].Rtt)
	assert.Equal(8.337, t.Sockets[2].Rttvar)
	assert.Equal(6.562, t.Sockets[3].Rtt)
	assert.Equal(12.05, t.Sockets[3].Rttvar)

	assert.Equal(float64(40), t.Sockets[2].Ato)
	assert.Equal(float64(40), t.Sockets[3].Ato)

	assert.Equal(uint32(32768), t.Sockets[2].Mss)
	assert.Equal(uint32(32768), t.Sockets[3].Mss)

	assert.Equal(uint32(65535), t.Sockets[2].Pmtu)
	assert.Equal(uint32(65535), t.Sockets[3].Pmtu)

	assert.Equal(uint32(536), t.Sockets[2].Rcvmss)
	assert.Equal(uint32(536), t.Sockets[3].Rcvmss)

	assert.Equal(uint32(65483), t.Sockets[2].Advmss)
	assert.Equal(uint32(65483), t.Sockets[3].Advmss)

	assert.Equal(uint32(1257), t.Sockets[2].BytesSent)
	assert.Equal(uint32(312), t.Sockets[3].BytesSent)

	assert.Equal(uint64(1257), t.Sockets[2].BytesAcked)
	assert.Equal(uint64(312), t.Sockets[3].BytesAcked)

	assert.Equal(uint64(3442), t.Sockets[2].BytesReceived)
	assert.Equal(uint64(2099), t.Sockets[3].BytesReceived)

	assert.Equal(uint32(1112), t.Sockets[2].SegsOut)
	assert.Equal(uint32(1109), t.Sockets[3].SegsOut)

	assert.Equal(uint32(1118), t.Sockets[2].SegsIn)
	assert.Equal(uint32(1112), t.Sockets[3].SegsIn)

	assert.Equal(uint32(7), t.Sockets[2].DataSegsIn)
	assert.Equal(uint32(7), t.Sockets[2].DataSegsOut)

	assert.Equal(uint32(4), t.Sockets[3].DataSegsIn)
	assert.Equal(uint32(4), t.Sockets[3].DataSegsOut)

	assert.Equal(uint32(12721397), t.Sockets[2].Lastsnd)
	assert.Equal(uint32(12721397), t.Sockets[2].Lastrcv)
	assert.Equal(uint32(3530), t.Sockets[2].Lastack)
	assert.Equal(uint32(12722864), t.Sockets[3].Lastsnd)
	assert.Equal(uint32(12722877), t.Sockets[3].Lastrcv)
	assert.Equal(uint32(3530), t.Sockets[3].Lastack)

	assert.Equal(594026.739, t.Sockets[2].Send)
	assert.Equal(1187918.880, t.Sockets[2].PacingRate)
	assert.Equal(23831272.720, t.Sockets[2].DeliveryRate)
	assert.Equal(399487.961, t.Sockets[3].Send)
	assert.Equal(798975.920, t.Sockets[3].PacingRate)
	assert.Equal(29127111.104, t.Sockets[3].DeliveryRate)

	assert.Equal(uint32(8), t.Sockets[2].Delivered)
	assert.Equal(uint32(53), t.Sockets[2].BusyMs)
	assert.Equal(uint32(65483), t.Sockets[2].RcvSpace)
	assert.Equal(uint32(65483), t.Sockets[2].RcvSsthresh)
	assert.Equal(0.01, t.Sockets[2].Minrtt)
	assert.Equal(uint32(65536), t.Sockets[2].SndWnd)

	assert.Equal(uint32(5), t.Sockets[3].Delivered)
	assert.Equal(uint32(53), t.Sockets[3].BusyMs)
	assert.Equal(uint32(65483), t.Sockets[3].RcvSpace)
	assert.Equal(uint32(65483), t.Sockets[3].RcvSsthresh)
	assert.Equal(0.009, t.Sockets[3].Minrtt)
	assert.Equal(uint32(65536), t.Sockets[3].SndWnd)
}
