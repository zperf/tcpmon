package influxdb

import (
	"fmt"
	"io"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"github.com/zperf/tcpmon/tcpmon/gproto"
	"github.com/zperf/tcpmon/tcpmon/tutils"
)

type LineProtocolExporter struct {
	hostname string
	writer   io.Writer
}

func New(hostname string, writer io.Writer) *LineProtocolExporter {
	return &LineProtocolExporter{
		hostname: hostname,
		writer:   writer,
	}
}

func (e *LineProtocolExporter) ExportMetric(m *gproto.Metric) {
	switch m := m.Body.(type) {
	case *gproto.Metric_Tcp:
		e.exportMetricTcp(m.Tcp)
	case *gproto.Metric_Net:
		e.exportMetricNet(m.Net)
	case *gproto.Metric_Nic:
		e.exportMetricNic(m.Nic)
	default:
		log.Fatal().Msg("Unknown metric type")
	}
}

func (e *LineProtocolExporter) Printf(format string, a ...any) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}

	_, err := fmt.Fprintf(e.writer, format, a...)
	tutils.FatalIf(err)
}

func (e *LineProtocolExporter) exportMetricTcp(m *gproto.TcpMetric) {
	ts := m.GetTimestamp()
	for _, s := range m.GetSockets() {
		processText := strings.Join(lo.Map(s.GetProcesses(), func(p *gproto.ProcessInfo, _ int) string {
			return fmt.Sprintf("cmd:%s;pid:%v;fd:%v", p.GetName(), p.GetPid(), p.GetFd())
		}), ";")
		var prefix string
		if processText == "" {
			prefix = fmt.Sprintf("tcp,LocalAddr=%s,PeerAddr=%s,Hostname=%s", s.GetLocalAddr(), s.GetPeerAddr(), e.hostname)
		} else {
			prefix = fmt.Sprintf("tcp,LocalAddr=%s,PeerAddr=%s,Hostname=%s,Process=%v", s.GetLocalAddr(), s.GetPeerAddr(), e.hostname, processText)
		}

		for _, timer := range s.GetTimers() {
			e.Printf("%s Timer=\"%v\",ExpireTimeUs=%v,Retrans=%v %v", prefix, timer.GetName(), timer.GetExpireTimeUs(), timer.GetRetrans(), ts)
		}

		e.Printf("%s State=\"%v\" %v\n", prefix, s.GetState(), ts)
		e.Printf("%s RecvQ=%v %v\n", prefix, s.GetRecvQ(), ts)
		e.Printf("%s SendQ=%v %v\n", prefix, s.GetSendQ(), ts)

		e.Printf("%s RmemAlloc=%v %v\n", prefix, s.GetSkmem().GetRmemAlloc(), ts)
		e.Printf("%s RcvBuf=%v %v\n", prefix, s.GetSkmem().GetRcvBuf(), ts)
		e.Printf("%s WmemAlloc=%v %v\n", prefix, s.GetSkmem().GetWmemAlloc(), ts)
		e.Printf("%s SndBuf=%v %v\n", prefix, s.GetSkmem().GetSndBuf(), ts)
		e.Printf("%s FwdAlloc=%v %v\n", prefix, s.GetSkmem().GetFwdAlloc(), ts)
		e.Printf("%s WmemQueued=%v %v\n", prefix, s.GetSkmem().GetWmemQueued(), ts)
		e.Printf("%s OptMem=%v %v\n", prefix, s.GetSkmem().GetOptMem(), ts)
		e.Printf("%s BackLog=%v %v\n", prefix, s.GetSkmem().GetBackLog(), ts)
		e.Printf("%s SockDrop=%v %v\n", prefix, s.GetSkmem().GetSockDrop(), ts)

		e.Printf("%s Ts=%v %v\n", prefix, tutils.Btoi(s.GetTs()), ts)
		e.Printf("%s Sack=%v %v\n", prefix, tutils.Btoi(s.GetSack()), ts)
		e.Printf("%s Cubic=%v %v\n", prefix, tutils.Btoi(s.GetCubic()), ts)
		e.Printf("%s AppLimited=%v %v\n", prefix, tutils.Btoi(s.GetAppLimited()), ts)

		e.Printf("%s PacingRate=%f %d\n", prefix, s.GetPacingRate(), ts)
		e.Printf("%s DeliveryRate=%f %d\n", prefix, s.GetDeliveryRate(), ts)
		e.Printf("%s Send=%f %d\n", prefix, s.GetSend(), ts)
		e.Printf("%s SndWscale=%v %v\n", prefix, s.GetSndWscale(), ts)
		e.Printf("%s RcvWscale=%v %v\n", prefix, s.GetRcvWscale(), ts)
		e.Printf("%s Rto=%f %d\n", prefix, s.GetRto(), ts)
		e.Printf("%s Rtt=%f %d\n", prefix, s.GetRtt(), ts)
		e.Printf("%s Rttvar=%f %d\n", prefix, s.GetRttvar(), ts)
		e.Printf("%s Minrtt=%f %d\n", prefix, s.GetMinrtt(), ts)
		e.Printf("%s RcvRtt=%f %d\n", prefix, s.GetRcvRtt(), ts)
		e.Printf("%s RetransNow=%v %v\n", prefix, s.GetRetransNow(), ts)
		e.Printf("%s RetransTotal=%v %v\n", prefix, s.GetRetransTotal(), ts)
		e.Printf("%s Ato=%f %d\n", prefix, s.GetAto(), ts)
		e.Printf("%s Mss=%v %v\n", prefix, s.GetMss(), ts)
		e.Printf("%s Pmtu=%v %v\n", prefix, s.GetPmtu(), ts)
		e.Printf("%s Rcvmss=%v %v\n", prefix, s.GetRcvmss(), ts)
		e.Printf("%s Advmss=%v %v\n", prefix, s.GetAdvmss(), ts)
		e.Printf("%s Cwnd=%v %v\n", prefix, s.GetCwnd(), ts)
		e.Printf("%s SndWnd=%v %v\n", prefix, s.GetSndWnd(), ts)
		e.Printf("%s BytesSent=%v %v\n", prefix, s.GetBytesSent(), ts)
		e.Printf("%s BytesAcked=%v %v\n", prefix, s.GetBytesAcked(), ts)
		e.Printf("%s BytesReceived=%v %v\n", prefix, s.GetBytesReceived(), ts)
		e.Printf("%s SegsOut=%v %v\n", prefix, s.GetSegsOut(), ts)
		e.Printf("%s SegsIn=%v %v\n", prefix, s.GetSegsIn(), ts)
		e.Printf("%s Lastsnd=%v %v\n", prefix, s.GetLastsnd(), ts)
		e.Printf("%s Lastrcv=%v %v\n", prefix, s.GetLastrcv(), ts)
		e.Printf("%s Lastack=%v %v\n", prefix, s.GetLastack(), ts)
		e.Printf("%s Delivered=%v %v\n", prefix, s.GetDelivered(), ts)
		e.Printf("%s BusyMs=%v %v\n", prefix, s.GetBusyMs(), ts)
		e.Printf("%s RcvSpace=%v %v\n", prefix, s.GetRcvSpace(), ts)
		e.Printf("%s RcvSsthresh=%v %v\n", prefix, s.GetRcvSsthresh(), ts)
		e.Printf("%s DataSegsOut=%v %v\n", prefix, s.GetDataSegsOut(), ts)
		e.Printf("%s DataSegsIn=%v %v\n", prefix, s.GetDataSegsIn(), ts)
		e.Printf("%s RwndLimited=%v %v\n", prefix, s.GetRwndLimited(), ts)
		e.Printf("%s SndbufLimited=%v %v\n", prefix, s.GetSndbufLimited(), ts)

		e.Printf("%s Ecn=%v %v\n", prefix, tutils.Btoi(s.GetEcn()), ts)
		e.Printf("%s Ecnseen=%v %v\n", prefix, tutils.Btoi(s.GetEcnseen()), ts)
	}
}

func (e *LineProtocolExporter) exportMetricNic(m *gproto.NicMetric) {
	ts := m.GetTimestamp()

	for _, i := range m.GetIfaces() {
		prefix := fmt.Sprintf("nic,Name=%v,Hostname=%v", i.GetName(), e.hostname)

		// rx
		e.Printf("%s RxErrors=%v %v\n", prefix, i.GetRxErrors(), ts)
		e.Printf("%s RxDropped=%v %v\n", prefix, i.GetRxDropped(), ts)
		e.Printf("%s RxOverruns=%v %v\n", prefix, i.GetRxOverruns(), ts)
		e.Printf("%s RxFrame=%v %v\n", prefix, i.GetRxFrame(), ts)

		// tx
		e.Printf("%s TxErrors=%v %v\n", prefix, i.GetTxErrors(), ts)
		e.Printf("%s TxDropped=%v %v\n", prefix, i.GetTxDropped(), ts)
		e.Printf("%s TxOverruns=%v %v\n", prefix, i.GetTxOverruns(), ts)
		e.Printf("%s TxCarrier=%v %v\n", prefix, i.GetTxCarrier(), ts)
		e.Printf("%s TxCollisions=%v %v\n", prefix, i.GetTxCollisions(), ts)
	}
}

func (e *LineProtocolExporter) exportMetricNet(m *gproto.NetstatMetric) {
	ts := m.GetTimestamp()
	prefix := fmt.Sprintf("net,Hostname=%v", e.hostname)

	e.Printf("%s IpForwarding=%v %v", prefix, m.GetIpForwarding(), ts)
	e.Printf("%s IpDefaultTtl=%v %v", prefix, m.GetIpDefaultTtl(), ts)
	e.Printf("%s IpInReceives=%v %v", prefix, m.GetIpInReceives(), ts)
	e.Printf("%s IpInHdrErrors=%v %v", prefix, m.GetIpInHdrErrors(), ts)
	e.Printf("%s IpInAddrErrors=%v %v", prefix, m.GetIpInAddrErrors(), ts)
	e.Printf("%s IpForwDatagrams=%v %v", prefix, m.GetIpForwDatagrams(), ts)
	e.Printf("%s IpInUnknownProtos=%v %v", prefix, m.GetIpInUnknownProtos(), ts)
	e.Printf("%s IpInDiscards=%v %v", prefix, m.GetIpInDiscards(), ts)
	e.Printf("%s IpInDelivers=%v %v", prefix, m.GetIpInDelivers(), ts)
	e.Printf("%s IpOutRequests=%v %v", prefix, m.GetIpOutRequests(), ts)
	e.Printf("%s IpOutDiscards=%v %v", prefix, m.GetIpOutDiscards(), ts)
	e.Printf("%s IpOutNoRoutes=%v %v", prefix, m.GetIpOutNoRoutes(), ts)
	e.Printf("%s IpReasmTimeout=%v %v", prefix, m.GetIpReasmTimeout(), ts)
	e.Printf("%s IpReasmReqds=%v %v", prefix, m.GetIpReasmReqds(), ts)
	e.Printf("%s IpReasmOks=%v %v", prefix, m.GetIpReasmOks(), ts)
	e.Printf("%s IpReasmFails=%v %v", prefix, m.GetIpReasmFails(), ts)
	e.Printf("%s IpFragOks=%v %v", prefix, m.GetIpFragOks(), ts)
	e.Printf("%s IpFragFails=%v %v", prefix, m.GetIpFragFails(), ts)
	e.Printf("%s IpFragCreates=%v %v", prefix, m.GetIpFragCreates(), ts)
	e.Printf("%s IpInNoRoutes=%v %v", prefix, m.GetIpInNoRoutes(), ts)
	e.Printf("%s IpInTruncatedPkts=%v %v", prefix, m.GetIpInTruncatedPkts(), ts)
	e.Printf("%s IpInMcastPkts=%v %v", prefix, m.GetIpInMcastPkts(), ts)
	e.Printf("%s IpOutMcastPkts=%v %v", prefix, m.GetIpOutMcastPkts(), ts)
	e.Printf("%s IpInBcastPkts=%v %v", prefix, m.GetIpInBcastPkts(), ts)
	e.Printf("%s IpOutBcastPkts=%v %v", prefix, m.GetIpOutBcastPkts(), ts)
	e.Printf("%s IpInOctets=%v %v", prefix, m.GetIpInOctets(), ts)
	e.Printf("%s IpOutOctets=%v %v", prefix, m.GetIpOutOctets(), ts)
	e.Printf("%s IpInMcastOctets=%v %v", prefix, m.GetIpInMcastOctets(), ts)
	e.Printf("%s IpOutMcastOctets=%v %v", prefix, m.GetIpOutMcastOctets(), ts)
	e.Printf("%s IpInBcastOctets=%v %v", prefix, m.GetIpInBcastOctets(), ts)
	e.Printf("%s IpOutBcastOctets=%v %v", prefix, m.GetIpOutBcastOctets(), ts)
	e.Printf("%s IpInCsumErrors=%v %v", prefix, m.GetIpInCsumErrors(), ts)
	e.Printf("%s IpInNoEctPkts=%v %v", prefix, m.GetIpInNoEctPkts(), ts)
	e.Printf("%s IpInEct1Pkts=%v %v", prefix, m.GetIpInEct1Pkts(), ts)
	e.Printf("%s IpInEct0Pkts=%v %v", prefix, m.GetIpInEct0Pkts(), ts)
	e.Printf("%s IpInCePkts=%v %v", prefix, m.GetIpInCePkts(), ts)
	e.Printf("%s IpReasmOverlaps=%v %v", prefix, m.GetIpReasmOverlaps(), ts)

	e.Printf("%s UdpInDatagrams=%v %v", prefix, m.GetUdpInDatagrams(), ts)
	e.Printf("%s UdpNoPorts=%v %v", prefix, m.GetUdpNoPorts(), ts)
	e.Printf("%s UdpInErrors=%v %v", prefix, m.GetUdpInErrors(), ts)
	e.Printf("%s UdpOutDatagrams=%v %v", prefix, m.GetUdpOutDatagrams(), ts)
	e.Printf("%s UdpRcvbufErrors=%v %v", prefix, m.GetUdpRcvbufErrors(), ts)
	e.Printf("%s UdpSndbufErrors=%v %v", prefix, m.GetUdpSndbufErrors(), ts)
	e.Printf("%s UdpInCsumErrors=%v %v", prefix, m.GetUdpInCsumErrors(), ts)
	e.Printf("%s UdpIgnoredMulti=%v %v", prefix, m.GetUdpIgnoredMulti(), ts)
	e.Printf("%s UdpMemErrors=%v %v", prefix, m.GetUdpMemErrors(), ts)

	e.Printf("%s TcpRtoAlgorithm=%v %v", prefix, m.GetTcpRtoAlgorithm(), ts)
	e.Printf("%s TcpRtoMin=%v %v", prefix, m.GetTcpRtoMin(), ts)
	e.Printf("%s TcpRtoMax=%v %v", prefix, m.GetTcpRtoMax(), ts)
	e.Printf("%s TcpMaxConn=%v %v", prefix, m.GetTcpMaxConn(), ts)
	e.Printf("%s TcpActiveOpens=%v %v", prefix, m.GetTcpActiveOpens(), ts)
	e.Printf("%s TcpPassiveOpens=%v %v", prefix, m.GetTcpPassiveOpens(), ts)
	e.Printf("%s TcpAttemptFails=%v %v", prefix, m.GetTcpAttemptFails(), ts)
	e.Printf("%s TcpEstabResets=%v %v", prefix, m.GetTcpEstabResets(), ts)
	e.Printf("%s TcpCurrEstab=%v %v", prefix, m.GetTcpCurrEstab(), ts)
	e.Printf("%s TcpInSegs=%v %v", prefix, m.GetTcpInSegs(), ts)
	e.Printf("%s TcpOutSegs=%v %v", prefix, m.GetTcpOutSegs(), ts)
	e.Printf("%s TcpRetransSegs=%v %v", prefix, m.GetTcpRetransSegs(), ts)
	e.Printf("%s TcpInErrs=%v %v", prefix, m.GetTcpInErrs(), ts)
	e.Printf("%s TcpOutRsts=%v %v", prefix, m.GetTcpOutRsts(), ts)
	e.Printf("%s TcpInCsumErrors=%v %v", prefix, m.GetTcpInCsumErrors(), ts)
	e.Printf("%s TcpSyncookiesSent=%v %v", prefix, m.GetTcpSyncookiesSent(), ts)
	e.Printf("%s TcpSyncookiesRecv=%v %v", prefix, m.GetTcpSyncookiesRecv(), ts)
	e.Printf("%s TcpSyncookiesFailed=%v %v", prefix, m.GetTcpSyncookiesFailed(), ts)
	e.Printf("%s TcpEmbryonicRsts=%v %v", prefix, m.GetTcpEmbryonicRsts(), ts)
	e.Printf("%s TcpPruneCalled=%v %v", prefix, m.GetTcpPruneCalled(), ts)
	e.Printf("%s TcpRcvPruned=%v %v", prefix, m.GetTcpRcvPruned(), ts)
	e.Printf("%s TcpOfoPruned=%v %v", prefix, m.GetTcpOfoPruned(), ts)
	e.Printf("%s TcpOutOfWindowIcmps=%v %v", prefix, m.GetTcpOutOfWindowIcmps(), ts)
	e.Printf("%s TcpLockDroppedIcmps=%v %v", prefix, m.GetTcpLockDroppedIcmps(), ts)
	e.Printf("%s TcpArpFilter=%v %v", prefix, m.GetTcpArpFilter(), ts)
	e.Printf("%s TcpTw=%v %v", prefix, m.GetTcpTw(), ts)
	e.Printf("%s TcpTwRecycled=%v %v", prefix, m.GetTcpTwRecycled(), ts)
	e.Printf("%s TcpTwKilled=%v %v", prefix, m.GetTcpTwKilled(), ts)
	e.Printf("%s TcpPawsActive=%v %v", prefix, m.GetTcpPawsActive(), ts)
	e.Printf("%s TcpPawsEstab=%v %v", prefix, m.GetTcpPawsEstab(), ts)
	e.Printf("%s TcpDelayedAcks=%v %v", prefix, m.GetTcpDelayedAcks(), ts)
	e.Printf("%s TcpDelayedAckLocked=%v %v", prefix, m.GetTcpDelayedAckLocked(), ts)
	e.Printf("%s TcpDelayedAckLost=%v %v", prefix, m.GetTcpDelayedAckLost(), ts)
	e.Printf("%s TcpListenOverflows=%v %v", prefix, m.GetTcpListenOverflows(), ts)
	e.Printf("%s TcpListenDrops=%v %v", prefix, m.GetTcpListenDrops(), ts)
	e.Printf("%s TcpHpHits=%v %v", prefix, m.GetTcpHpHits(), ts)
	e.Printf("%s TcpPureAcks=%v %v", prefix, m.GetTcpPureAcks(), ts)
	e.Printf("%s TcpHpAcks=%v %v", prefix, m.GetTcpHpAcks(), ts)
	e.Printf("%s TcpRenoRecovery=%v %v", prefix, m.GetTcpRenoRecovery(), ts)
	e.Printf("%s TcpSackRecovery=%v %v", prefix, m.GetTcpSackRecovery(), ts)
	e.Printf("%s TcpSackReneging=%v %v", prefix, m.GetTcpSackReneging(), ts)
	e.Printf("%s TcpSackReorder=%v %v", prefix, m.GetTcpSackReorder(), ts)
	e.Printf("%s TcpRenoReorder=%v %v", prefix, m.GetTcpRenoReorder(), ts)
	e.Printf("%s TcpTsReorder=%v %v", prefix, m.GetTcpTsReorder(), ts)
	e.Printf("%s TcpFullUndo=%v %v", prefix, m.GetTcpFullUndo(), ts)
	e.Printf("%s TcpPartialUndo=%v %v", prefix, m.GetTcpPartialUndo(), ts)
	e.Printf("%s TcpDsackUndo=%v %v", prefix, m.GetTcpDsackUndo(), ts)
	e.Printf("%s TcpLossUndo=%v %v", prefix, m.GetTcpLossUndo(), ts)
	e.Printf("%s TcpLostRetransmit=%v %v", prefix, m.GetTcpLostRetransmit(), ts)
	e.Printf("%s TcpRenoFailures=%v %v", prefix, m.GetTcpRenoFailures(), ts)
	e.Printf("%s TcpSackFailures=%v %v", prefix, m.GetTcpSackFailures(), ts)
	e.Printf("%s TcpLossFailures=%v %v", prefix, m.GetTcpLossFailures(), ts)
	e.Printf("%s TcpFastRetrans=%v %v", prefix, m.GetTcpFastRetrans(), ts)
	e.Printf("%s TcpSlowStartRetrans=%v %v", prefix, m.GetTcpSlowStartRetrans(), ts)
	e.Printf("%s TcpTimeouts=%v %v", prefix, m.GetTcpTimeouts(), ts)
	e.Printf("%s TcpLossProbes=%v %v", prefix, m.GetTcpLossProbes(), ts)
	e.Printf("%s TcpLossProbeRecovery=%v %v", prefix, m.GetTcpLossProbeRecovery(), ts)
	e.Printf("%s TcpRenoRecoveryFail=%v %v", prefix, m.GetTcpRenoRecoveryFail(), ts)
	e.Printf("%s TcpSackRecoveryFail=%v %v", prefix, m.GetTcpSackRecoveryFail(), ts)
	e.Printf("%s TcpRcvCollapsed=%v %v", prefix, m.GetTcpRcvCollapsed(), ts)
	e.Printf("%s TcpBacklogCoalesce=%v %v", prefix, m.GetTcpBacklogCoalesce(), ts)
	e.Printf("%s TcpDsackOldSent=%v %v", prefix, m.GetTcpDsackOldSent(), ts)
	e.Printf("%s TcpDsackOfoSent=%v %v", prefix, m.GetTcpDsackOfoSent(), ts)
	e.Printf("%s TcpDsackRecv=%v %v", prefix, m.GetTcpDsackRecv(), ts)
	e.Printf("%s TcpDsackOfoRecv=%v %v", prefix, m.GetTcpDsackOfoRecv(), ts)
	e.Printf("%s TcpAbortOnData=%v %v", prefix, m.GetTcpAbortOnData(), ts)
	e.Printf("%s TcpAbortOnClose=%v %v", prefix, m.GetTcpAbortOnClose(), ts)
	e.Printf("%s TcpAbortOnMemory=%v %v", prefix, m.GetTcpAbortOnMemory(), ts)
	e.Printf("%s TcpAbortOnTimeout=%v %v", prefix, m.GetTcpAbortOnTimeout(), ts)
	e.Printf("%s TcpAbortOnLinger=%v %v", prefix, m.GetTcpAbortOnLinger(), ts)
	e.Printf("%s TcpAbortFailed=%v %v", prefix, m.GetTcpAbortFailed(), ts)
	e.Printf("%s TcpMemoryPressures=%v %v", prefix, m.GetTcpMemoryPressures(), ts)
	e.Printf("%s TcpMemoryPressuresChrono=%v %v", prefix, m.GetTcpMemoryPressuresChrono(), ts)
	e.Printf("%s TcpSackDiscard=%v %v", prefix, m.GetTcpSackDiscard(), ts)
	e.Printf("%s TcpDsackIgnoredOld=%v %v", prefix, m.GetTcpDsackIgnoredOld(), ts)
	e.Printf("%s TcpDsackIgnoredNoUndo=%v %v", prefix, m.GetTcpDsackIgnoredNoUndo(), ts)
	e.Printf("%s TcpSpuriousRtos=%v %v", prefix, m.GetTcpSpuriousRtos(), ts)
	e.Printf("%s TcpMd5NotFound=%v %v", prefix, m.GetTcpMd5NotFound(), ts)
	e.Printf("%s TcpMd5Unexpected=%v %v", prefix, m.GetTcpMd5Unexpected(), ts)
	e.Printf("%s TcpMd5Failure=%v %v", prefix, m.GetTcpMd5Failure(), ts)
	e.Printf("%s TcpSackShifted=%v %v", prefix, m.GetTcpSackShifted(), ts)
	e.Printf("%s TcpSackMerged=%v %v", prefix, m.GetTcpSackMerged(), ts)
	e.Printf("%s TcpSackShiftFallback=%v %v", prefix, m.GetTcpSackShiftFallback(), ts)
	e.Printf("%s TcpBacklogDrop=%v %v", prefix, m.GetTcpBacklogDrop(), ts)
	e.Printf("%s TcpPfMemallocDrop=%v %v", prefix, m.GetTcpPfMemallocDrop(), ts)
	e.Printf("%s TcpMinTtlDrop=%v %v", prefix, m.GetTcpMinTtlDrop(), ts)
	e.Printf("%s TcpDeferAcceptDrop=%v %v", prefix, m.GetTcpDeferAcceptDrop(), ts)
	e.Printf("%s TcpIpReversePathFilter=%v %v", prefix, m.GetTcpIpReversePathFilter(), ts)
	e.Printf("%s TcpTimeWaitOverflow=%v %v", prefix, m.GetTcpTimeWaitOverflow(), ts)
	e.Printf("%s TcpReqQFullDoCookies=%v %v", prefix, m.GetTcpReqQFullDoCookies(), ts)
	e.Printf("%s TcpReqQFullDrop=%v %v", prefix, m.GetTcpReqQFullDrop(), ts)
	e.Printf("%s TcpRetransFail=%v %v", prefix, m.GetTcpRetransFail(), ts)
	e.Printf("%s TcpRcvCoalesce=%v %v", prefix, m.GetTcpRcvCoalesce(), ts)
	e.Printf("%s TcpOfoQueue=%v %v", prefix, m.GetTcpOfoQueue(), ts)
	e.Printf("%s TcpOfoDrop=%v %v", prefix, m.GetTcpOfoDrop(), ts)
	e.Printf("%s TcpOfoMerge=%v %v", prefix, m.GetTcpOfoMerge(), ts)
	e.Printf("%s TcpChallengeAck=%v %v", prefix, m.GetTcpChallengeAck(), ts)
	e.Printf("%s TcpSynChallenge=%v %v", prefix, m.GetTcpSynChallenge(), ts)
	e.Printf("%s TcpFastOpenActive=%v %v", prefix, m.GetTcpFastOpenActive(), ts)
	e.Printf("%s TcpFastOpenActiveFail=%v %v", prefix, m.GetTcpFastOpenActiveFail(), ts)
	e.Printf("%s TcpFastOpenPassive=%v %v", prefix, m.GetTcpFastOpenPassive(), ts)
	e.Printf("%s TcpFastOpenPassiveFail=%v %v", prefix, m.GetTcpFastOpenPassiveFail(), ts)
	e.Printf("%s TcpFastOpenListenOverflow=%v %v", prefix, m.GetTcpFastOpenListenOverflow(), ts)
	e.Printf("%s TcpFastOpenCookieReqd=%v %v", prefix, m.GetTcpFastOpenCookieReqd(), ts)
	e.Printf("%s TcpFastOpenBlackhole=%v %v", prefix, m.GetTcpFastOpenBlackhole(), ts)
	e.Printf("%s TcpSpuriousRtxHostQueues=%v %v", prefix, m.GetTcpSpuriousRtxHostQueues(), ts)
	e.Printf("%s TcpBusyPollRxPackets=%v %v", prefix, m.GetTcpBusyPollRxPackets(), ts)
	e.Printf("%s TcpAutoCorking=%v %v", prefix, m.GetTcpAutoCorking(), ts)
	e.Printf("%s TcpFromZeroWindowAdv=%v %v", prefix, m.GetTcpFromZeroWindowAdv(), ts)
	e.Printf("%s TcpToZeroWindowAdv=%v %v", prefix, m.GetTcpToZeroWindowAdv(), ts)
	e.Printf("%s TcpWantZeroWindowAdv=%v %v", prefix, m.GetTcpWantZeroWindowAdv(), ts)
	e.Printf("%s TcpSynRetrans=%v %v", prefix, m.GetTcpSynRetrans(), ts)
	e.Printf("%s TcpOrigDataSent=%v %v", prefix, m.GetTcpOrigDataSent(), ts)
	e.Printf("%s TcpHystartTrainDetect=%v %v", prefix, m.GetTcpHystartTrainDetect(), ts)
	e.Printf("%s TcpHystartTrainCwnd=%v %v", prefix, m.GetTcpHystartTrainCwnd(), ts)
	e.Printf("%s TcpHystartDelayDetect=%v %v", prefix, m.GetTcpHystartDelayDetect(), ts)
	e.Printf("%s TcpHystartDelayCwnd=%v %v", prefix, m.GetTcpHystartDelayCwnd(), ts)
	e.Printf("%s TcpAckSkippedSynRecv=%v %v", prefix, m.GetTcpAckSkippedSynRecv(), ts)
	e.Printf("%s TcpAckSkippedPaws=%v %v", prefix, m.GetTcpAckSkippedPaws(), ts)
	e.Printf("%s TcpAckSkippedSeq=%v %v", prefix, m.GetTcpAckSkippedSeq(), ts)
	e.Printf("%s TcpAckSkippedFinWait2=%v %v", prefix, m.GetTcpAckSkippedFinWait2(), ts)
	e.Printf("%s TcpAckSkippedTimeWait=%v %v", prefix, m.GetTcpAckSkippedTimeWait(), ts)
	e.Printf("%s TcpAckSkippedChallenge=%v %v", prefix, m.GetTcpAckSkippedChallenge(), ts)
	e.Printf("%s TcpWinProbe=%v %v", prefix, m.GetTcpWinProbe(), ts)
	e.Printf("%s TcpKeepAlive=%v %v", prefix, m.GetTcpKeepAlive(), ts)
	e.Printf("%s TcpMtupFail=%v %v", prefix, m.GetTcpMtupFail(), ts)
	e.Printf("%s TcpMtupSuccess=%v %v", prefix, m.GetTcpMtupSuccess(), ts)
	e.Printf("%s TcpDelivered=%v %v", prefix, m.GetTcpDelivered(), ts)
	e.Printf("%s TcpDeliveredCe=%v %v", prefix, m.GetTcpDeliveredCe(), ts)
	e.Printf("%s TcpAckCompressed=%v %v", prefix, m.GetTcpAckCompressed(), ts)
	e.Printf("%s TcpZeroWindowDrop=%v %v", prefix, m.GetTcpZeroWindowDrop(), ts)
	e.Printf("%s TcpRcvQDrop=%v %v", prefix, m.GetTcpRcvQDrop(), ts)
	e.Printf("%s TcpWqueueTooBig=%v %v", prefix, m.GetTcpWqueueTooBig(), ts)
	e.Printf("%s TcpFastOpenPassiveAltKey=%v %v", prefix, m.GetTcpFastOpenPassiveAltKey(), ts)
	e.Printf("%s TcpTimeoutRehash=%v %v", prefix, m.GetTcpTimeoutRehash(), ts)
	e.Printf("%s TcpDuplicateDataRehash=%v %v", prefix, m.GetTcpDuplicateDataRehash(), ts)
	e.Printf("%s TcpDsackRecvSegs=%v %v", prefix, m.GetTcpDsackRecvSegs(), ts)
	e.Printf("%s TcpDsackIgnoredDubious=%v %v", prefix, m.GetTcpDsackIgnoredDubious(), ts)
	e.Printf("%s TcpMigrateReqSuccess=%v %v", prefix, m.GetTcpMigrateReqSuccess(), ts)
	e.Printf("%s TcpMigrateReqFailure=%v %v", prefix, m.GetTcpMigrateReqFailure(), ts)
	e.Printf("%s TcpPlbRehash=%v %v", prefix, m.GetTcpPlbRehash(), ts)

	e.Printf("%s IcmpInMsgs=%v %v", prefix, m.GetIcmpInMsgs(), ts)
	e.Printf("%s IcmpInErrors=%v %v", prefix, m.GetIcmpInErrors(), ts)
	e.Printf("%s IcmpInCsumErrors=%v %v", prefix, m.GetIcmpInCsumErrors(), ts)
	e.Printf("%s IcmpInDestUnreachs=%v %v", prefix, m.GetIcmpInDestUnreachs(), ts)
	e.Printf("%s IcmpInTimeExcds=%v %v", prefix, m.GetIcmpInTimeExcds(), ts)
	e.Printf("%s IcmpInParmProbs=%v %v", prefix, m.GetIcmpInParmProbs(), ts)
	e.Printf("%s IcmpInSrcQuenchs=%v %v", prefix, m.GetIcmpInSrcQuenchs(), ts)
	e.Printf("%s IcmpInRedirects=%v %v", prefix, m.GetIcmpInRedirects(), ts)
	e.Printf("%s IcmpInEchos=%v %v", prefix, m.GetIcmpInEchos(), ts)
	e.Printf("%s IcmpInEchoReps=%v %v", prefix, m.GetIcmpInEchoReps(), ts)
	e.Printf("%s IcmpInTimestamps=%v %v", prefix, m.GetIcmpInTimestamps(), ts)
	e.Printf("%s IcmpInTimestampReps=%v %v", prefix, m.GetIcmpInTimestampReps(), ts)
	e.Printf("%s IcmpInAddrMasks=%v %v", prefix, m.GetIcmpInAddrMasks(), ts)
	e.Printf("%s IcmpInAddrMaskReps=%v %v", prefix, m.GetIcmpInAddrMaskReps(), ts)
	e.Printf("%s IcmpOutMsgs=%v %v", prefix, m.GetIcmpOutMsgs(), ts)
	e.Printf("%s IcmpOutErrors=%v %v", prefix, m.GetIcmpOutErrors(), ts)
	e.Printf("%s IcmpOutRateLimitGlobal=%v %v", prefix, m.GetIcmpOutRateLimitGlobal(), ts)
	e.Printf("%s IcmpOutRateLimitHost=%v %v", prefix, m.GetIcmpOutRateLimitHost(), ts)
	e.Printf("%s IcmpOutDestUnreachs=%v %v", prefix, m.GetIcmpOutDestUnreachs(), ts)
	e.Printf("%s IcmpOutTimeExcds=%v %v", prefix, m.GetIcmpOutTimeExcds(), ts)
	e.Printf("%s IcmpOutParmProbs=%v %v", prefix, m.GetIcmpOutParmProbs(), ts)
	e.Printf("%s IcmpOutSrcQuenchs=%v %v", prefix, m.GetIcmpOutSrcQuenchs(), ts)
	e.Printf("%s IcmpOutRedirects=%v %v", prefix, m.GetIcmpOutRedirects(), ts)
	e.Printf("%s IcmpOutEchos=%v %v", prefix, m.GetIcmpOutEchos(), ts)
	e.Printf("%s IcmpOutEchoReps=%v %v", prefix, m.GetIcmpOutEchoReps(), ts)
	e.Printf("%s IcmpOutTimestamps=%v %v", prefix, m.GetIcmpOutTimestamps(), ts)
	e.Printf("%s IcmpOutTimestampReps=%v %v", prefix, m.GetIcmpOutTimestampReps(), ts)
	e.Printf("%s IcmpOutAddrMasks=%v %v", prefix, m.GetIcmpOutAddrMasks(), ts)
	e.Printf("%s IcmpOutAddrMaskReps=%v %v", prefix, m.GetIcmpOutAddrMaskReps(), ts)
}
