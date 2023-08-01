package tcpmon_test

import (
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	. "github.com/zperf/tcpmon/tcpmon"
)

func (s *CommandParserTestSuite) TestParseNetstat() {
	assert := s.Assert()

	out := `Ip:
    1066283 total packets received
    0 forwarded
    0 incoming packets discarded
    1061338 incoming packets delivered
    917653 requests sent out
    16 outgoing packets dropped
    2 reassemblies required
    1 packets reassembled ok
Icmp:
    36 ICMP messages received
    0 input ICMP message failed.
    ICMP input histogram:
        destination unreachable: 34
        echo requests: 2
    39 ICMP messages sent
    0 ICMP messages failed
    ICMP output histogram:
        destination unreachable: 37
        echo replies: 2
IcmpMsg:
        InType3: 34
        InType8: 2
        OutType0: 2
        OutType3: 37
Tcp:
    566 active connections openings
    81 passive connection openings
    120 failed connection attempts
    9 connection resets received
    5 connections established
    1059991 segments received
    944977 segments send out
    2328 segments retransmited
    0 bad segments received.
    256 resets sent
Udp:
    1277 packets received
    34 packets to unknown port received.
    0 packet receive errors
    1352 packets sent
    0 receive buffer errors
    0 send buffer errors
UdpLite:
TcpExt:
    282 TCP sockets finished time wait in fast timer
    74 packets rejects in established connections because of timestamp
    33878 delayed acks sent
    6 delayed acks further delayed because of locked socket
    Quick ack mode was activated 892 times
    26 packets directly queued to recvmsg prequeue.
    2404 bytes directly received in process context from prequeue
    406412 packet headers predicted
    2 packets header predicted and directly queued to user
    91197 acknowledgments not containing data payload received
    238951 predicted acknowledgments
    227 times recovered from packet loss by selective acknowledgements
    Detected reordering 1 times using SACK
    Detected reordering 3 times using time stamp
    3 congestion windows partially recovered using Hoe heuristic
    131 congestion windows recovered without slow start by DSACK
    138 congestion windows recovered without slow start after partial ack
    202 timeouts after SACK recovery
    2 timeouts in loss state
    324 fast retransmits
    15 forward retransmits
    48 retransmits in slow start
    209 other TCP timeouts
    TCPLossProbes: 837
    TCPLossProbeRecovery: 558
    16 SACK retransmits failed
    950 DSACKs sent for old packets
    12 DSACKs sent for out of order packets
    627 DSACKs received
    29 DSACKs for out of order packets received
    12 connections reset due to unexpected data
    4 connections reset due to early user close
    128 connections aborted due to timeout
    TCPDSACKIgnoredOld: 1
    TCPDSACKIgnoredNoUndo: 177
    TCPSpuriousRTOs: 62
    TCPSackShifted: 24
    TCPSackMerged: 232
    TCPSackShiftFallback: 1304
    IPReversePathFilter: 1
    TCPRcvCoalesce: 83625
    TCPOFOQueue: 73474
    TCPOFOMerge: 18
    TCPAutoCorking: 2273
    TCPSynRetrans: 827
    TCPOrigDataSent: 484717
    TCPHystartDelayDetect: 4
    TCPHystartDelayCwnd: 608
IpExt:
    InNoRoutes: 2
    InMcastPkts: 2
    InBcastPkts: 4355
    InOctets: 587672631
    OutOctets: 228218297
    InMcastOctets: 72
    InBcastOctets: 673050
    InNoECTPkts: 1021973
    InECT0Pkts: 44310`

	lines := strings.FieldsFunc(out, func(r rune) bool {
		return r == '\n'
	})

	var r NetstatMetric
	r.Timestamp = timestamppb.New(time.Now())
	r.Type = MetricType_NET
	ParseNetstatOutput(&r, lines)

	assert.Equal(uint32(1066283), r.IpTotalPacketsReceived)
	assert.Equal(uint32(0), r.IpForwarded)
	assert.Equal(uint32(0), r.IpIncomingPacketsDiscarded)
	assert.Equal(uint32(1061338), r.IpIncomingPacketsDelivered)
	assert.Equal(uint32(917653), r.IpRequestsSentOut)
	assert.Equal(uint32(16), r.IpOutgoingPacketsDropped)

	assert.Equal(uint32(566), r.TcpActiveConnectionsOpenings)
	assert.Equal(uint32(81), r.TcpPassiveConnectionOpenings)
	assert.Equal(uint32(120), r.TcpFailedConnectionAttempts)
	assert.Equal(uint32(9), r.TcpConnectionResetsReceived)
	assert.Equal(uint32(5), r.TcpConnectionsEstablished)
	assert.Equal(uint32(1059991), r.TcpSegmentsReceived)
	assert.Equal(uint32(944977), r.TcpSegmentsSendOut)
	assert.Equal(uint32(2328), r.TcpSegmentsRetransmited)
	assert.Equal(uint32(0), r.TcpBadSegmentsReceived)
	assert.Equal(uint32(256), r.TcpResetsSent)

	assert.Equal(uint32(1277), r.UdpPacketsReceived)
	assert.Equal(uint32(34), r.UdpPacketsToUnknownPortReceived)
	assert.Equal(uint32(0), r.UdpPacketReceiveErrors)
	assert.Equal(uint32(1352), r.UdpPacketsSent)
	assert.Equal(uint32(0), r.UdpReceiveBufferErrors)
	assert.Equal(uint32(0), r.UdpSendBufferErrors)
}
