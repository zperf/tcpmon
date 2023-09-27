package tcpmon

import (
	"context"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"

	"github.com/zperf/tcpmon/tcpmon/tutils"
)

var headSet = map[string]struct{}{
	"Ip:":      {},
	"Icmp:":    {},
	"IcmpMsg:": {},
	"Tcp:":     {},
	"Udp:":     {},
	"UdpLite:": {},
	"TcpExt:":  {},
	"IpExt:":   {},
}

func ParseNetstatOutput(r *NetstatMetric, out []string) {
	flag := ""
	for _, line := range out {
		if _, exist := headSet[line]; exist {
			flag = line
			continue
		}
		if flag == "Ip:" {
			if strings.Contains(line, "total packets received") {
				r.IpTotalPacketsReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "forwarded") {
				r.IpForwarded, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "incoming packets discarded") {
				r.IpIncomingPacketsDiscarded, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "incoming packets delivered") {
				r.IpIncomingPacketsDelivered, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "requests sent out") {
				r.IpRequestsSentOut, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "outgoing packets dropped") {
				r.IpOutgoingPacketsDropped, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "dropped because of missing route") {
				r.IpDroppedBecauseOfMissingRoute, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "fragments dropped after timeout") {
				r.IpFragmentsDroppedAfterTimeout, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "reassemblies required") {
				r.IpReassembliesRequired, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "packets reassembled ok") {
				r.IpPacketsReassembledOk, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "packet reassembles failed") {
				r.IpPacketReassemblesFailed, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			}
		} else if flag == "Tcp:" {
			if strings.Contains(line, "active connections openings") {
				r.TcpActiveConnectionsOpenings, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "passive connection openings") {
				r.TcpPassiveConnectionOpenings, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "failed connection attempts") {
				r.TcpFailedConnectionAttempts, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "connection resets received") {
				r.TcpConnectionResetsReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "connections established") {
				r.TcpConnectionsEstablished, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "segments received") && !strings.Contains(line, "bad") {
				r.TcpSegmentsReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "segments send out") {
				r.TcpSegmentsSendOut, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "segments retransmit") {
				// NOTE: in the newer Linux version (like fc38), `netstat -s | grep retrans` will return retransmitted
				// The older (like el7) will return a typo: retransmited
				// We should take the common prefix `retransmit`
				r.TcpSegmentsRetransmitted, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "bad segments received") {
				r.TcpBadSegmentsReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "resets sent") {
				r.TcpResetsSent, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			}
		} else if flag == "Udp:" {
			if strings.Contains(line, "packets received") {
				r.UdpPacketsReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "packets to unknown port received") {
				r.UdpPacketsToUnknownPortReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "packet receive errors") {
				r.UdpPacketReceiveErrors, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "packets sent") {
				r.UdpPacketsSent, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "receive buffer errors") {
				r.UdpReceiveBufferErrors, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "send buffer errors") {
				r.UdpSendBufferErrors, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "IgnoredMulti:") {
				r.UdpIgnoredMulti, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			}
		} else if flag == "TcpExt:" {
			if strings.Contains(line, "invalid SYN cookies received") {
				r.TcpextInvalidSynCookiesReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "resets received for embryonic SYN_RECV sockets") {
				r.TcpextResetsReceivedForEmbryonicSynRecvSockets, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "ICMP packets dropped because they were out-of-window") {
				r.TcpextIcmpPacketsDroppedBecauseTheyWereOutOfWindow, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "ICMP packets dropped because socket was locked") {
				r.TcpextIcmpPacketsDroppedBecauseSocketWasLocked, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "TCP sockets finished time wait in fast timer") {
				r.TcpextTcpSocketsFinishedTimeWaitInFastTimer, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "packets rejects in established connections because of timestamp") {
				r.TcpextPacketsRejectsInEstablishedConnectionsBecauseOfTimestamp, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "delayed acks sent") {
				r.TcpextDelayedAcksSent, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "delayed acks further delayed because of locked socket") {
				r.TcpextDelayedAcksFurtherDelayedBecauseOfLockedSocket, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "Quick ack mode was activated") {
				r.TcpextQuickAckModeWasActivatedTimes, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[5])
			} else if strings.Contains(line, "times the listen queue of a socket overflowed") {
				r.TcpextTimesTheListenQueueOfASocketOverflowed, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "SYNs to LISTEN sockets dropped") {
				r.TcpextSynsToListenSocketsDropped, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "packet headers predicted") {
				r.TcpextPacketHeadersPredicted, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "acknowledgments not containing data payload received") {
				r.TcpextAcknowledgmentsNotContainingDataPayloadReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "predicted acknowledgments") {
				r.TcpextPredictedAcknowledgments, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "times recovered from packet loss by selective acknowledgements") {
				r.TcpextTimesRecoveredFromPacketLossBySelectiveAcknowledgements, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "Detected reordering") {
				r.TcpextDetectedReorderingTimesUsingSack, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[2])
			} else if strings.Contains(line, "congestion windows fully recovered without slow start") {
				r.TcpextCongestionWindowsFullyRecoveredWithoutSlowStart, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "congestion windows recovered without slow start by DSACK") {
				r.TcpextCongestionWindowsRecoveredWithoutSlowStartByDsack, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "congestion windows recovered without slow start after partial ack") {
				r.TcpextCongestionWindowsRecoveredWithoutSlowStartAfterPartialAck, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "TCPLostRetransmit:") {
				r.TcpextTcpLostRetransmit, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "timeouts after SACK recovery") {
				r.TcpextTimeoutsAfterSackRecovery, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "timeouts in loss state") {
				r.TcpextTimeoutsInLossState, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "fast retransmits") {
				r.TcpextFastRetransmits, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "retransmits in slow start") {
				r.TcpextRetransmitsInSlowStart, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "other TCP timeouts") {
				r.TcpextOtherTcpTimeouts, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "TCPLossProbes:") {
				r.TcpextTcpLossProbes, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPLossProbeRecovery:") {
				r.TcpextTcpLossProbeRecovery, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "SACK retransmits failed") {
				r.TcpextSackRetransmitsFailed, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "TCPBacklogCoalesce:") {
				r.TcpextTcpBacklogCoalesce, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "DSACKs sent for old packets") {
				r.TcpextDsacksSentForOldPackets, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "DSACKs sent for out of order packets") {
				r.TcpextDsacksSentForOutOfOrderPackets, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "DSACKs received") {
				r.TcpextDsacksReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "DSACKs for out of order packets received") {
				r.TcpextDsacksForOutOfOrderPacketsReceived, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "connections reset due to unexpected data") {
				r.TcpextConnectionsResetDueToUnexpectedData, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "connections reset due to early user close") {
				r.TcpextConnectionsResetDueToEarlyUserClose, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "connections aborted due to timeout") {
				r.TcpextConnectionsAbortedDueToTimeout, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[0])
			} else if strings.Contains(line, "TCPDSACKIgnoredOld:") {
				r.TcpextTcpDsackIgnoredOld, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPDSACKIgnoredNoUndo:") {
				r.TcpextTcpDsackIgnoredNoUndo, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPSpuriousRTOs:") {
				r.TcpextTcpSpuriousRtos, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPSackMerged:") {
				r.TcpextTcpSackMerged, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPSackShiftFallback:") {
				r.TcpextTcpSackShiftFallback, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "IPReversePathFilter:") {
				r.TcpextIpReversePathFilter, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPRetransFail:") {
				r.TcpextTcpRetransFail, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPRcvCoalesce:") {
				r.TcpextTcpRcvCoalesce, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPOFOQueue:") {
				r.TcpextTcpOfoQueue, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPOFOMerge:") {
				r.TcpextTcpOfoMerge, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPChallengeACK:") {
				r.TcpextTcpChallengeAck, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPSYNChallenge:") {
				r.TcpextTcpSynChallenge, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPFastOpenActiveFail:") {
				r.TcpextTcpFastOpenActiveFail, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPFastOpenCookieReqd:") {
				r.TcpextTcpFastOpenCookieReqd, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPFastOpenBlackhole:") {
				r.TcpextTcpFastOpenBlackhole, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPSpuriousRtxHostQueues:") {
				r.TcpextTcpSpuriousRtxHostQueues, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPAutoCorking:") {
				r.TcpextTcpAutoCorking, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPSynRetrans:") {
				r.TcpextTcpSynRetrans, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPOrigDataSent:") {
				r.TcpextTcpOrigDataSent, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPHystartTrainDetect:") {
				r.TcpextTcpHystartTrainDetect, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPHystartTrainCwnd:") {
				r.TcpextTcpHystartTrainCwnd, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPHystartDelayDetect:") {
				r.TcpextTcpHystartDelayDetect, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPHystartDelayCwnd:") {
				r.TcpextTcpHystartDelayCwnd, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPACKSkippedSeq:") {
				r.TcpextTcpAckSkippedSeq, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPACKSkippedChallenge:") {
				r.TcpextTcpAckSkippedChallenge, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPKeepAlive:") {
				r.TcpextTcpKeepAlive, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPDelivered:") {
				r.TcpextTcpDelivered, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "TCPAckCompressed:") {
				r.TcpextTcpAckCompressed, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			}
		} else if flag == "IpExt:" {
			if strings.Contains(line, "InBcastPkts:") {
				r.IpextInBcastPkts, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "InOctets:") {
				r.IpextInOctets, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "OutOctets:") {
				r.IpextOutOctets, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "InBcastOctets:") {
				r.IpextInBcastOctets, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "InNoECTPkts:") {
				r.IpextInNoEctPkts, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			} else if strings.Contains(line, "InECT0Pkts:") {
				r.IpextInEct0Pkts, _ = tutils.ParseUint64(strings.FieldsFunc(line, tutils.SplitSpace)[1])
			}
		}
	}
}

func (m *NetstatMonitor) RunNetstat(now time.Time) (*NetstatMetric, string, error) {
	c := cmd.NewCmd(m.config.PathNetstat, m.config.ArgNetstat)
	ctx, cancel := context.WithTimeout(context.Background(), m.config.Timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, "", errors.Wrap(ctx.Err(), "netstat timeout")
	case st := <-c.Start():
		var r NetstatMetric
		r.Type = MetricType_NET
		r.Timestamp = now.Unix()

		ParseNetstatOutput(&r, st.Stdout)
		return &r, strings.Join(st.Stdout, "\n"), nil
	}
}
