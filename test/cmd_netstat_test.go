package test

import (
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	. "github.com/zperf/tcpmon/tcpmon"
)

func (s *CommandParserTestSuite) TestParseNetstat() {
	assert := s.Assert()

	out := `Ip:
    12984693843 total packets received
    2 with invalid headers
    13 with invalid addresses
    0 forwarded
    0 incoming packets discarded
    12975893623 incoming packets delivered
    11393863778 requests sent out
    2960 outgoing packets dropped
    776 dropped because of missing route
    12 reassemblies required
    6 packets reassembled ok
    184 fragments received ok
    368 fragments created
Icmp:
    140154534 ICMP messages received
    14186 input ICMP message failed.
    ICMP input histogram:
        destination unreachable: 48182
        redirects: 77
        echo requests: 70633563
        echo replies: 69472712
    140181726 ICMP messages sent
    0 ICMP messages failed
    ICMP output histogram:
        destination unreachable: 48265
        echo request: 69499898
        echo replies: 70633563
IcmpMsg:
        InType0: 69472712
        InType3: 48182
        InType5: 77
        InType8: 70633563
        OutType0: 70633563
        OutType3: 48265
        OutType8: 69499898
Tcp:
    49845605 active connections openings
    46253346 passive connection openings
    6129349 failed connection attempts
    2184394 connection resets received
    892 connections established
    12596717510 segments received
    31113950007 segments send out
    48107538 segments retransmited
    27 bad segments received.
    11229045 resets sent
Udp:
    196786471 packets received
    6480 packets to unknown port received.
    68957158 packet receive errors
    265752378 packets sent
    68957158 receive buffer errors
    0 send buffer errors
    IgnoredMulti: 4363047
UdpLite:
TcpExt:
    70 invalid SYN cookies received
    20 resets received for embryonic SYN_RECV sockets
    27 ICMP packets dropped because they were out-of-window
    52061656 TCP sockets finished time wait in fast timer
    7 time wait sockets recycled by time stamp
    125745 packets rejects in established connections because of timestamp
    112076710 delayed acks sent
    158300 delayed acks further delayed because of locked socket
    Quick ack mode was activated 9998003 times
    981 times the listen queue of a socket overflowed
    981 SYNs to LISTEN sockets dropped
    7168650378 packet headers predicted
    1890902648 acknowledgments not containing data payload received
    3592506430 predicted acknowledgments
    5453091 times recovered from packet loss by selective acknowledgements
    Detected reordering 3544733 times using SACK
    Detected reordering 624427 times using time stamp
    55620 congestion windows fully recovered without slow start
    612208 congestion windows partially recovered using Hoe heuristic
    356499 congestion windows recovered without slow start by DSACK
    7545 congestion windows recovered without slow start after partial ack
    TCPLostRetransmit: 883886
    1203 timeouts after SACK recovery
    878 timeouts in loss state
    44563971 fast retransmits
    798742 retransmits in slow start
    27969 other TCP timeouts
    TCPLossProbes: 5213051
    TCPLossProbeRecovery: 41798
    105391 SACK retransmits failed
    TCPBacklogCoalesce: 104127568
    10483485 DSACKs sent for old packets
    82041 DSACKs sent for out of order packets
    4482837 DSACKs received
    3902 DSACKs for out of order packets received
    1010427 connections reset due to unexpected data
    2444880 connections reset due to early user close
    3215 connections aborted due to timeout
    TCPDSACKIgnoredOld: 241397
    TCPDSACKIgnoredNoUndo: 2161962
    TCPSpuriousRTOs: 148
    TCPSackShifted: 4552047
    TCPSackMerged: 19855226
    TCPSackShiftFallback: 51614523
    IPReversePathFilter: 33206
    TCPRcvCoalesce: 756764461
    TCPOFOQueue: 171697744
    TCPOFOMerge: 105963
    TCPChallengeACK: 113666
    TCPSYNChallenge: 30
    TCPSpuriousRtxHostQueues: 119227
    TCPAutoCorking: 110895312
    TCPFromZeroWindowAdv: 30453
    TCPToZeroWindowAdv: 30453
    TCPWantZeroWindowAdv: 1192398
    TCPSynRetrans: 23366
    TCPOrigDataSent: 27644619399
    TCPHystartTrainDetect: 700103
    TCPHystartTrainCwnd: 18378617
    TCPHystartDelayDetect: 6798
    TCPHystartDelayCwnd: 660016
    TCPACKSkippedSynRecv: 1303
    TCPACKSkippedPAWS: 103693
    TCPACKSkippedSeq: 38587
    TCPACKSkippedFinWait2: 201
    TCPACKSkippedTimeWait: 184
    TCPACKSkippedChallenge: 65148
    TCPWinProbe: 1085
    TCPKeepAlive: 55886023
    TCPDelivered: 27689076642
    TCPAckCompressed: 101435862
IpExt:
    InBcastPkts: 4363057
    InOctets: 84129285172077
    OutOctets: 67110093964178
    InBcastOctets: 1400581363
    InNoECTPkts: 43749298150
    InECT0Pkts: 13952
`

	lines := strings.FieldsFunc(out, func(r rune) bool {
		return r == '\n'
	})

	var r NetstatMetric
	r.Timestamp = timestamppb.New(time.Now())
	r.Type = MetricType_NET
	ParseNetstatOutput(&r, lines)

	assert.Equal(uint64(12984693843), r.IpTotalPacketsReceived)
	assert.Equal(uint64(0), r.IpForwarded)
	assert.Equal(uint64(0), r.IpIncomingPacketsDiscarded)
	assert.Equal(uint64(12975893623), r.IpIncomingPacketsDelivered)
	assert.Equal(uint64(11393863778), r.IpRequestsSentOut)
	assert.Equal(uint64(2960), r.IpOutgoingPacketsDropped)

	assert.Equal(uint64(49845605), r.TcpActiveConnectionsOpenings)
	assert.Equal(uint64(46253346), r.TcpPassiveConnectionOpenings)
	assert.Equal(uint64(6129349), r.TcpFailedConnectionAttempts)
	assert.Equal(uint64(2184394), r.TcpConnectionResetsReceived)
	assert.Equal(uint64(892), r.TcpConnectionsEstablished)
	assert.Equal(uint64(12596717510), r.TcpSegmentsReceived)
	assert.Equal(uint64(31113950007), r.TcpSegmentsSendOut)
	assert.Equal(uint64(48107538), r.TcpSegmentsRetransmitted)
	assert.Equal(uint64(27), r.TcpBadSegmentsReceived)
	assert.Equal(uint64(11229045), r.TcpResetsSent)

	assert.Equal(uint64(196786471), r.UdpPacketsReceived)
	assert.Equal(uint64(6480), r.UdpPacketsToUnknownPortReceived)
	assert.Equal(uint64(68957158), r.UdpPacketReceiveErrors)
	assert.Equal(uint64(265752378), r.UdpPacketsSent)
	assert.Equal(uint64(68957158), r.UdpReceiveBufferErrors)
	assert.Equal(uint64(0), r.UdpSendBufferErrors)
}
