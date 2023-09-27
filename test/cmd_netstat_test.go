package test

import (
	"strings"
	"time"

	. "github.com/zperf/tcpmon/tcpmon"
)

func (s *CommandParserTestSuite) TestParseNetstat() {
	assert := s.Assert()

	out := `Ip:
    890512330 total packets received
    0 forwarded
    0 incoming packets discarded
    887690339 incoming packets delivered
    932042022 requests sent out
    56 dropped because of missing route
    3 fragments dropped after timeout
    26 reassemblies required
    9 packets reassembled ok
    3 packet reassembles failed
Icmp:
    55830975 ICMP messages received
    1174778 input ICMP message failed.
    ICMP input histogram:
        destination unreachable: 7990801
        echo requests: 24009093
        echo replies: 23831081
    64307755 ICMP messages sent
    0 ICMP messages failed
    ICMP output histogram:
        destination unreachable: 7993040
        echo request: 32305622
        echo replies: 24009093
IcmpMsg:
        InType0: 23831081
        InType3: 7990801
        InType8: 24009093
        OutType0: 24009093
        OutType3: 7993040
        OutType8: 32305622
Tcp:
    45658411 active connections openings
    14439757 passive connection openings
    23758981 failed connection attempts
    5424888 connection resets received
    831 connections established
    869435153 segments received
    885611852 segments send out
    24128338 segments retransmited
    2514 bad segments received.
    24671546 resets sent
Udp:
    16511838 packets received
    34381 packets to unknown port received.
    0 packet receive errors
    16806205 packets sent
    0 receive buffer errors
    0 send buffer errors
    IgnoredMulti: 368716
UdpLite:
TcpExt:
    3 invalid SYN cookies received
    43575 resets received for embryonic SYN_RECV sockets
    4 ICMP packets dropped because they were out-of-window
    1745301 ICMP packets dropped because socket was locked
    10684967 TCP sockets finished time wait in fast timer
    277 packets rejects in established connections because of timestamp
    22797882 delayed acks sent
    38468 delayed acks further delayed because of locked socket
    Quick ack mode was activated 1589477 times
    379435 times the listen queue of a socket overflowed
    379436 SYNs to LISTEN sockets dropped
    128152589 packet headers predicted
    178793072 acknowledgments not containing data payload received
    187444989 predicted acknowledgments
    39769 times recovered from packet loss by selective acknowledgements
    Detected reordering 49 times using SACK
    2 congestion windows fully recovered without slow start
    73294 congestion windows recovered without slow start by DSACK
    527956 congestion windows recovered without slow start after partial ack
    TCPLostRetransmit: 10426796
    1519 timeouts after SACK recovery
    165 timeouts in loss state
    40029 fast retransmits
    17395 retransmits in slow start
    10879298 other TCP timeouts
    TCPLossProbes: 2796415
    TCPLossProbeRecovery: 489372
    13 SACK retransmits failed
    TCPBacklogCoalesce: 19866
    1590088 DSACKs sent for old packets
    71 DSACKs sent for out of order packets
    994179 DSACKs received
    21 DSACKs for out of order packets received
    1426199 connections reset due to unexpected data
    1976252 connections reset due to early user close
    120916 connections aborted due to timeout
    TCPDSACKIgnoredOld: 2
    TCPDSACKIgnoredNoUndo: 499062
    TCPSpuriousRTOs: 7
    TCPSackMerged: 102
    TCPSackShiftFallback: 494376
    IPReversePathFilter: 2360
    TCPRetransFail: 890057
    TCPRcvCoalesce: 26238647
    TCPOFOQueue: 49049
    TCPOFOMerge: 71
    TCPChallengeACK: 3002
    TCPSYNChallenge: 2990
    TCPFastOpenActiveFail: 34
    TCPFastOpenCookieReqd: 1
    TCPFastOpenBlackhole: 34
    TCPSpuriousRtxHostQueues: 1881224
    TCPAutoCorking: 18930750
    TCPSynRetrans: 19293870
    TCPOrigDataSent: 327761287
    TCPHystartTrainDetect: 588
    TCPHystartTrainCwnd: 10501
    TCPHystartDelayDetect: 3
    TCPHystartDelayCwnd: 58
    TCPACKSkippedSeq: 43
    TCPACKSkippedChallenge: 1
    TCPKeepAlive: 134113931
    TCPDelivered: 340465026
    TCPAckCompressed: 4
IpExt:
    InBcastPkts: 368719
    InOctets: 354960781370
    OutOctets: 356737196235
    InBcastOctets: 90329354
    InNoECTPkts: 892134044
    InECT0Pkts: 9
`

	lines := strings.FieldsFunc(out, func(r rune) bool {
		return r == '\n'
	})

	var r NetstatMetric
	r.Timestamp = time.Now().UnixMilli()
	r.Type = MetricType_NET
	ParseNetstatOutput(&r, lines)

	assert.Equal(uint64(890512330), r.IpTotalPacketsReceived)
	assert.Equal(uint64(0), r.IpForwarded)
	assert.Equal(uint64(0), r.IpIncomingPacketsDiscarded)
	assert.Equal(uint64(887690339), r.IpIncomingPacketsDelivered)
	assert.Equal(uint64(932042022), r.IpRequestsSentOut)
	assert.Equal(uint64(56), r.IpDroppedBecauseOfMissingRoute)
	assert.Equal(uint64(3), r.IpFragmentsDroppedAfterTimeout)
	assert.Equal(uint64(26), r.IpReassembliesRequired)
	assert.Equal(uint64(9), r.IpPacketsReassembledOk)
	assert.Equal(uint64(3), r.IpPacketReassemblesFailed)

	assert.Equal(uint64(45658411), r.TcpActiveConnectionsOpenings)
	assert.Equal(uint64(14439757), r.TcpPassiveConnectionOpenings)
	assert.Equal(uint64(23758981), r.TcpFailedConnectionAttempts)
	assert.Equal(uint64(5424888), r.TcpConnectionResetsReceived)
	assert.Equal(uint64(831), r.TcpConnectionsEstablished)
	assert.Equal(uint64(869435153), r.TcpSegmentsReceived)
	assert.Equal(uint64(885611852), r.TcpSegmentsSendOut)
	assert.Equal(uint64(24128338), r.TcpSegmentsRetransmitted)
	assert.Equal(uint64(2514), r.TcpBadSegmentsReceived)
	assert.Equal(uint64(24671546), r.TcpResetsSent)

	assert.Equal(uint64(16511838), r.UdpPacketsReceived)
	assert.Equal(uint64(34381), r.UdpPacketsToUnknownPortReceived)
	assert.Equal(uint64(0), r.UdpPacketReceiveErrors)
	assert.Equal(uint64(16806205), r.UdpPacketsSent)
	assert.Equal(uint64(0), r.UdpReceiveBufferErrors)
	assert.Equal(uint64(0), r.UdpSendBufferErrors)
	assert.Equal(uint64(368716), r.UdpIgnoredMulti)

	assert.Equal(uint64(3), r.GetTcpextInvalidSynCookiesReceived())
	assert.Equal(uint64(43575), r.GetTcpextResetsReceivedForEmbryonicSynRecvSockets())
	assert.Equal(uint64(4), r.GetTcpextIcmpPacketsDroppedBecauseTheyWereOutOfWindow())
	assert.Equal(uint64(1745301), r.GetTcpextIcmpPacketsDroppedBecauseSocketWasLocked())
	assert.Equal(uint64(10684967), r.GetTcpextTcpSocketsFinishedTimeWaitInFastTimer())
	assert.Equal(uint64(277), r.GetTcpextPacketsRejectsInEstablishedConnectionsBecauseOfTimestamp())
	assert.Equal(uint64(22797882), r.GetTcpextDelayedAcksSent())
	assert.Equal(uint64(38468), r.GetTcpextDelayedAcksFurtherDelayedBecauseOfLockedSocket())
	assert.Equal(uint64(1589477), r.GetTcpextQuickAckModeWasActivatedTimes())
	assert.Equal(uint64(379435), r.GetTcpextTimesTheListenQueueOfASocketOverflowed())
	assert.Equal(uint64(379436), r.GetTcpextSynsToListenSocketsDropped())
	assert.Equal(uint64(128152589), r.GetTcpextPacketHeadersPredicted())
	assert.Equal(uint64(178793072), r.GetTcpextAcknowledgmentsNotContainingDataPayloadReceived())
	assert.Equal(uint64(187444989), r.GetTcpextPredictedAcknowledgments())
	assert.Equal(uint64(39769), r.GetTcpextTimesRecoveredFromPacketLossBySelectiveAcknowledgements())
	assert.Equal(uint64(49), r.GetTcpextDetectedReorderingTimesUsingSack())
	assert.Equal(uint64(2), r.GetTcpextCongestionWindowsFullyRecoveredWithoutSlowStart())
	assert.Equal(uint64(73294), r.GetTcpextCongestionWindowsRecoveredWithoutSlowStartByDsack())
	assert.Equal(uint64(527956), r.GetTcpextCongestionWindowsRecoveredWithoutSlowStartAfterPartialAck())
	assert.Equal(uint64(10426796), r.GetTcpextTcpLostRetransmit())
	assert.Equal(uint64(1519), r.GetTcpextTimeoutsAfterSackRecovery())
	assert.Equal(uint64(165), r.GetTcpextTimeoutsInLossState())
	assert.Equal(uint64(40029), r.GetTcpextFastRetransmits())
	assert.Equal(uint64(17395), r.GetTcpextRetransmitsInSlowStart())
	assert.Equal(uint64(10879298), r.GetTcpextOtherTcpTimeouts())
	assert.Equal(uint64(2796415), r.GetTcpextTcpLossProbes())
	assert.Equal(uint64(489372), r.GetTcpextTcpLossProbeRecovery())
	assert.Equal(uint64(13), r.GetTcpextSackRetransmitsFailed())
	assert.Equal(uint64(19866), r.GetTcpextTcpBacklogCoalesce())
	assert.Equal(uint64(1590088), r.GetTcpextDsacksSentForOldPackets())
	assert.Equal(uint64(71), r.GetTcpextDsacksSentForOutOfOrderPackets())
	assert.Equal(uint64(994179), r.GetTcpextDsacksReceived())
	assert.Equal(uint64(21), r.GetTcpextDsacksForOutOfOrderPacketsReceived())
	assert.Equal(uint64(1426199), r.GetTcpextConnectionsResetDueToUnexpectedData())
	assert.Equal(uint64(1976252), r.GetTcpextConnectionsResetDueToEarlyUserClose())
	assert.Equal(uint64(120916), r.GetTcpextConnectionsAbortedDueToTimeout())
	assert.Equal(uint64(2), r.GetTcpextTcpDsackIgnoredOld())
	assert.Equal(uint64(499062), r.GetTcpextTcpDsackIgnoredNoUndo())
	assert.Equal(uint64(7), r.GetTcpextTcpSpuriousRtos())
	assert.Equal(uint64(102), r.GetTcpextTcpSackMerged())
	assert.Equal(uint64(494376), r.GetTcpextTcpSackShiftFallback())
	assert.Equal(uint64(2360), r.GetTcpextIpReversePathFilter())
	assert.Equal(uint64(890057), r.GetTcpextTcpRetransFail())
	assert.Equal(uint64(26238647), r.GetTcpextTcpRcvCoalesce())
	assert.Equal(uint64(49049), r.GetTcpextTcpOfoQueue())
	assert.Equal(uint64(71), r.GetTcpextTcpOfoMerge())
	assert.Equal(uint64(3002), r.GetTcpextTcpChallengeAck())
	assert.Equal(uint64(2990), r.GetTcpextTcpSynChallenge())
	assert.Equal(uint64(34), r.GetTcpextTcpFastOpenActiveFail())
	assert.Equal(uint64(1), r.GetTcpextTcpFastOpenCookieReqd())
	assert.Equal(uint64(34), r.GetTcpextTcpFastOpenBlackhole())
	assert.Equal(uint64(1881224), r.GetTcpextTcpSpuriousRtxHostQueues())
	assert.Equal(uint64(18930750), r.GetTcpextTcpAutoCorking())
	assert.Equal(uint64(19293870), r.GetTcpextTcpSynRetrans())
	assert.Equal(uint64(327761287), r.GetTcpextTcpOrigDataSent())
	assert.Equal(uint64(588), r.GetTcpextTcpHystartTrainDetect())
	assert.Equal(uint64(10501), r.GetTcpextTcpHystartTrainCwnd())
	assert.Equal(uint64(3), r.GetTcpextTcpHystartDelayDetect())
	assert.Equal(uint64(58), r.GetTcpextTcpHystartDelayCwnd())
	assert.Equal(uint64(43), r.GetTcpextTcpAckSkippedSeq())
	assert.Equal(uint64(1), r.GetTcpextTcpAckSkippedChallenge())
	assert.Equal(uint64(134113931), r.GetTcpextTcpKeepAlive())
	assert.Equal(uint64(340465026), r.GetTcpextTcpDelivered())
	assert.Equal(uint64(4), r.GetTcpextTcpAckCompressed())

	assert.Equal(uint64(368719), r.GetIpextInBcastPkts())
	assert.Equal(uint64(354960781370), r.GetIpextInOctets())
	assert.Equal(uint64(356737196235), r.GetIpextOutOctets())
	assert.Equal(uint64(90329354), r.GetIpextInBcastOctets())
	assert.Equal(uint64(892134044), r.GetIpextInNoEctPkts())
	assert.Equal(uint64(9), r.GetIpextInEct0Pkts())
}
