package parsing

import (
	"io"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"

	"github.com/zperf/tcpmon/tcpmon/tproto"
)

func ParseNetstat(r io.Reader, m *tproto.NetstatMetric) error {
	return parseSnmpOrNetstat(r, m, ParseNetstatLine)
}

func ParseNetstatLine(fieldStr string, valueStr string, m *tproto.NetstatMetric) error {
	t := parseProcStatType(fieldStr)
	if t != parseProcStatType(valueStr) {
		log.Fatal().Str("title", fieldStr).Str("value", valueStr).
			Msg("mismatched title and value")
	}

	if t == ProcNetStatMPTcpExt {
		// I don't care MPTcp
		return nil
	}

	fields := strings.Fields(fieldStr)
	values := strings.Fields(valueStr)

	for i := 1; i < len(fields); i++ {
		field := fields[i]
		value, err := strconv.ParseUint(values[i], 10, 64)
		if err != nil {
			return errors.Wrap(err, "parse failed")
		}

		switch field {
		case "SyncookiesSent":
			m.TcpSyncookiesSent = value
		case "SyncookiesRecv":
			m.TcpSyncookiesRecv = value
		case "SyncookiesFailed":
			m.TcpSyncookiesFailed = value
		case "EmbryonicRsts":
			m.TcpEmbryonicRsts = value
		case "PruneCalled":
			m.TcpPruneCalled = value
		case "RcvPruned":
			m.TcpRcvPruned = value
		case "OfoPruned":
			m.TcpOfoPruned = value
		case "OutOfWindowIcmps":
			m.TcpOutOfWindowIcmps = value
		case "LockDroppedIcmps":
			m.TcpLockDroppedIcmps = value
		case "ArpFilter":
			m.TcpArpFilter = value
		case "TW":
			m.TcpTw = value
		case "TWRecycled":
			m.TcpTwRecycled = value
		case "TWKilled":
			m.TcpTwKilled = value
		case "PAWSActive":
			m.TcpPawsActive = value
		case "PAWSEstab":
			m.TcpPawsEstab = value
		case "DelayedACKs":
			m.TcpDelayedAcks = value
		case "DelayedACKLocked":
			m.TcpDelayedAckLocked = value
		case "DelayedACKLost":
			m.TcpDelayedAckLost = value
		case "ListenOverflows":
			m.TcpListenOverflows = value
		case "ListenDrops":
			m.TcpListenDrops = value
		case "TCPHPHits":
			m.TcpHpHits = value
		case "TCPPureAcks":
			m.TcpPureAcks = value
		case "TCPHPAcks":
			m.TcpHpAcks = value
		case "TCPRenoRecovery":
			m.TcpRenoRecovery = value
		case "TCPSackRecovery":
			m.TcpSackRecovery = value
		case "TCPSACKReneging":
			m.TcpSackReneging = value
		case "TCPSACKReorder":
			m.TcpSackReorder = value
		case "TCPRenoReorder":
			m.TcpRenoReorder = value
		case "TCPTSReorder":
			m.TcpTsReorder = value
		case "TCPFullUndo":
			m.TcpFullUndo = value
		case "TCPPartialUndo":
			m.TcpPartialUndo = value
		case "TCPDSACKUndo":
			m.TcpDsackUndo = value
		case "TCPLossUndo":
			m.TcpLossUndo = value
		case "TCPLostRetransmit":
			m.TcpLostRetransmit = value
		case "TCPRenoFailures":
			m.TcpRenoFailures = value
		case "TCPSackFailures":
			m.TcpSackFailures = value
		case "TCPLossFailures":
			m.TcpLossFailures = value
		case "TCPFastRetrans":
			m.TcpFastRetrans = value
		case "TCPSlowStartRetrans":
			m.TcpSlowStartRetrans = value
		case "TCPTimeouts":
			m.TcpTimeouts = value
		case "TCPLossProbes":
			m.TcpLossProbes = value
		case "TCPLossProbeRecovery":
			m.TcpLossProbeRecovery = value
		case "TCPRenoRecoveryFail":
			m.TcpRenoRecoveryFail = value
		case "TCPSackRecoveryFail":
			m.TcpSackRecoveryFail = value
		case "TCPRcvCollapsed":
			m.TcpRcvCollapsed = value
		case "TCPBacklogCoalesce":
			m.TcpBacklogCoalesce = value
		case "TCPDSACKOldSent":
			m.TcpDsackOldSent = value
		case "TCPDSACKOfoSent":
			m.TcpDsackOfoSent = value
		case "TCPDSACKRecv":
			m.TcpDsackRecv = value
		case "TCPDSACKOfoRecv":
			m.TcpDsackOfoRecv = value
		case "TCPAbortOnData":
			m.TcpAbortOnData = value
		case "TCPAbortOnClose":
			m.TcpAbortOnClose = value
		case "TCPAbortOnMemory":
			m.TcpAbortOnMemory = value
		case "TCPAbortOnTimeout":
			m.TcpAbortOnTimeout = value
		case "TCPAbortOnLinger":
			m.TcpAbortOnLinger = value
		case "TCPAbortFailed":
			m.TcpAbortFailed = value
		case "TCPMemoryPressures":
			m.TcpMemoryPressures = value
		case "TCPMemoryPressuresChrono":
			m.TcpMemoryPressuresChrono = value
		case "TCPSACKDiscard":
			m.TcpSackDiscard = value
		case "TCPDSACKIgnoredOld":
			m.TcpDsackIgnoredOld = value
		case "TCPDSACKIgnoredNoUndo":
			m.TcpDsackIgnoredNoUndo = value
		case "TCPSpuriousRTOs":
			m.TcpSpuriousRtos = value
		case "TCPMD5NotFound":
			m.TcpMd5NotFound = value
		case "TCPMD5Unexpected":
			m.TcpMd5Unexpected = value
		case "TCPMD5Failure":
			m.TcpMd5Failure = value
		case "TCPSackShifted":
			m.TcpSackShifted = value
		case "TCPSackMerged":
			m.TcpSackMerged = value
		case "TCPSackShiftFallback":
			m.TcpSackShiftFallback = value
		case "TCPBacklogDrop":
			m.TcpBacklogDrop = value
		case "PFMemallocDrop":
			m.TcpPfMemallocDrop = value
		case "TCPMinTTLDrop":
			m.TcpMinTtlDrop = value
		case "TCPDeferAcceptDrop":
			m.TcpDeferAcceptDrop = value
		case "IPReversePathFilter":
			m.TcpIpReversePathFilter = value
		case "TCPTimeWaitOverflow":
			m.TcpTimeWaitOverflow = value
		case "TCPReqQFullDoCookies":
			m.TcpReqQFullDoCookies = value
		case "TCPReqQFullDrop":
			m.TcpReqQFullDrop = value
		case "TCPRetransFail":
			m.TcpRetransFail = value
		case "TCPRcvCoalesce":
			m.TcpRcvCoalesce = value
		case "TCPOFOQueue":
			m.TcpOfoQueue = value
		case "TCPOFODrop":
			m.TcpOfoDrop = value
		case "TCPOFOMerge":
			m.TcpOfoMerge = value
		case "TCPChallengeACK":
			m.TcpChallengeAck = value
		case "TCPSYNChallenge":
			m.TcpSynChallenge = value
		case "TCPFastOpenActive":
			m.TcpFastOpenActive = value
		case "TCPFastOpenActiveFail":
			m.TcpFastOpenActiveFail = value
		case "TCPFastOpenPassive":
			m.TcpFastOpenPassive = value
		case "TCPFastOpenPassiveFail":
			m.TcpFastOpenPassiveFail = value
		case "TCPFastOpenListenOverflow":
			m.TcpFastOpenListenOverflow = value
		case "TCPFastOpenCookieReqd":
			m.TcpFastOpenCookieReqd = value
		case "TCPFastOpenBlackhole":
			m.TcpFastOpenBlackhole = value
		case "TCPSpuriousRtxHostQueues":
			m.TcpSpuriousRtxHostQueues = value
		case "BusyPollRxPackets":
			m.TcpBusyPollRxPackets = value
		case "TCPAutoCorking":
			m.TcpAutoCorking = value
		case "TCPFromZeroWindowAdv":
			m.TcpFromZeroWindowAdv = value
		case "TCPToZeroWindowAdv":
			m.TcpToZeroWindowAdv = value
		case "TCPWantZeroWindowAdv":
			m.TcpWantZeroWindowAdv = value
		case "TCPSynRetrans":
			m.TcpSynRetrans = value
		case "TCPOrigDataSent":
			m.TcpOrigDataSent = value
		case "TCPHystartTrainDetect":
			m.TcpHystartTrainDetect = value
		case "TCPHystartTrainCwnd":
			m.TcpHystartTrainCwnd = value
		case "TCPHystartDelayDetect":
			m.TcpHystartDelayDetect = value
		case "TCPHystartDelayCwnd":
			m.TcpHystartDelayCwnd = value
		case "TCPACKSkippedSynRecv":
			m.TcpAckSkippedSynRecv = value
		case "TCPACKSkippedPAWS":
			m.TcpAckSkippedPaws = value
		case "TCPACKSkippedSeq":
			m.TcpAckSkippedSeq = value
		case "TCPACKSkippedFinWait2":
			m.TcpAckSkippedFinWait2 = value
		case "TCPACKSkippedTimeWait":
			m.TcpAckSkippedTimeWait = value
		case "TCPACKSkippedChallenge":
			m.TcpAckSkippedChallenge = value
		case "TCPWinProbe":
			m.TcpWinProbe = value
		case "TCPKeepAlive":
			m.TcpKeepAlive = value
		case "TCPMTUPFail":
			m.TcpMtupFail = value
		case "TCPMTUPSuccess":
			m.TcpMtupSuccess = value
		case "TCPDelivered":
			m.TcpDelivered = value
		case "TCPDeliveredCE":
			m.TcpDeliveredCe = value
		case "TCPAckCompressed":
			m.TcpAckCompressed = value
		case "TCPZeroWindowDrop":
			m.TcpZeroWindowDrop = value
		case "TCPRcvQDrop":
			m.TcpRcvQDrop = value
		case "TCPWqueueTooBig":
			m.TcpWqueueTooBig = value
		case "TCPFastOpenPassiveAltKey":
			m.TcpFastOpenPassiveAltKey = value
		case "TcpTimeoutRehash":
			m.TcpTimeoutRehash = value
		case "TcpDuplicateDataRehash":
			m.TcpDuplicateDataRehash = value
		case "TCPDSACKRecvSegs":
			m.TcpDsackRecvSegs = value
		case "TCPDSACKIgnoredDubious":
			m.TcpDsackIgnoredDubious = value
		case "TCPMigrateReqSuccess":
			m.TcpMigrateReqSuccess = value
		case "TCPMigrateReqFailure":
			m.TcpMigrateReqFailure = value
		case "TCPPLBRehash":
			m.TcpPlbRehash = value
		case "InNoRoutes":
			m.IpInNoRoutes = value
		case "InTruncatedPkts":
			m.IpInTruncatedPkts = value
		case "InMcastPkts":
			m.IpInMcastPkts = value
		case "OutMcastPkts":
			m.IpOutMcastPkts = value
		case "InBcastPkts":
			m.IpInBcastPkts = value
		case "OutBcastPkts":
			m.IpOutBcastPkts = value
		case "InOctets":
			m.IpInOctets = value
		case "OutOctets":
			m.IpOutOctets = value
		case "InMcastOctets":
			m.IpInMcastOctets = value
		case "OutMcastOctets":
			m.IpOutMcastOctets = value
		case "InBcastOctets":
			m.IpInBcastOctets = value
		case "OutBcastOctets":
			m.IpOutBcastOctets = value
		case "InCsumErrors":
			m.IpInCsumErrors = value
		case "InNoECTPkts":
			m.IpInNoEctPkts = value
		case "InECT1Pkts":
			m.IpInEct1Pkts = value
		case "InECT0Pkts":
			m.IpInEct0Pkts = value
		case "InCEPkts":
			m.IpInCePkts = value
		case "ReasmOverlaps":
			m.IpReasmOverlaps = value
		}
	}

	return nil
}
