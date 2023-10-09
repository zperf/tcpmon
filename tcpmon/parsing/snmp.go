package parsing

import (
	"io"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"

	"github.com/zperf/tcpmon/tcpmon/tproto"
)

func ParseSnmp(r io.Reader, m *tproto.NetstatMetric) error {
	return parseSnmpOrNetstat(r, m, ParseSnmpLine)
}

func ParseSnmpLine(fieldStr string, valueStr string, m *tproto.NetstatMetric) error {
	t := parseProcStatType(fieldStr)
	if t != parseProcStatType(valueStr) {
		log.Fatal().Str("title", fieldStr).Str("value", valueStr).
			Msg("mismatched title and value")
	}

	if t == ProcNetStatIcmpMsg || t == ProcNetStatUdpLite {
		// I don't care ICMP_MSG and UDP_LITE
		return nil
	}

	fields := strings.Fields(fieldStr)
	values := strings.Fields(valueStr)

	for i := 1; i < len(fields); i++ {
		field := fields[i]

		if t == ProcNetStatTcp && field == "MaxConn" {
			// TcpMaxConn field is signed, RFC 2012, net/ipv4/proc.c
			v, err := strconv.ParseInt(values[i], 10, 64)
			if err != nil {
				return errors.Wrap(err, "parse failed")
			}
			m.TcpMaxConn = v
			continue
		}

		value, err := strconv.ParseUint(values[i], 10, 64)
		if err != nil {
			return errors.Wrap(err, "parse failed")
		}

		switch field {
		case "InCsumErrors":
			switch t {
			case ProcNetStatIcmp:
				m.IcmpInCsumErrors = value
			case ProcNetStatUdp:
				m.UdpInCsumErrors = value
			case ProcNetStatTcp:
				m.TcpInCsumErrors = value
			default:
				// ignore UDP lite
			}

		case "InErrors":
			switch t {
			case ProcNetStatIcmp:
				m.IcmpInErrors = value
			case ProcNetStatUdp:
				m.UdpInErrors = value
			default:
				// ignore UDP lite
			}

		case "InMsgs":
			m.IcmpInMsgs = value
		case "InDestUnreachs":
			m.IcmpInDestUnreachs = value
		case "InTimeExcds":
			m.IcmpInTimeExcds = value
		case "InParmProbs":
			m.IcmpInParmProbs = value
		case "InSrcQuenchs":
			m.IcmpInSrcQuenchs = value
		case "InRedirects":
			m.IcmpInRedirects = value
		case "InEchos":
			m.IcmpInEchos = value
		case "InEchoReps":
			m.IcmpInEchoReps = value
		case "InTimestamps":
			m.IcmpInTimestamps = value
		case "InTimestampReps":
			m.IcmpInTimestampReps = value
		case "InAddrMasks":
			m.IcmpInAddrMasks = value
		case "InAddrMaskReps":
			m.IcmpInAddrMaskReps = value
		case "OutMsgs":
			m.IcmpOutMsgs = value
		case "OutErrors":
			m.IcmpOutErrors = value
		case "OutRateLimitGlobal":
			m.IcmpOutRateLimitGlobal = value
		case "OutRateLimitHost":
			m.IcmpOutRateLimitHost = value
		case "OutDestUnreachs":
			m.IcmpOutDestUnreachs = value
		case "OutTimeExcds":
			m.IcmpOutTimeExcds = value
		case "OutParmProbs":
			m.IcmpOutParmProbs = value
		case "OutSrcQuenchs":
			m.IcmpOutSrcQuenchs = value
		case "OutRedirects":
			m.IcmpOutRedirects = value
		case "OutEchos":
			m.IcmpOutEchos = value
		case "OutEchoReps":
			m.IcmpOutEchoReps = value
		case "OutTimestamps":
			m.IcmpOutTimestamps = value
		case "OutTimestampReps":
			m.IcmpOutTimestampReps = value
		case "OutAddrMasks":
			m.IcmpOutAddrMasks = value
		case "OutAddrMaskReps":
			m.IcmpOutAddrMaskReps = value

		// ip
		case "Forwarding":
			m.IpForwarding = value
		case "DefaultTTL":
			m.IpDefaultTtl = value
		case "InReceives":
			m.IpInReceives = value
		case "InHdrErrors":
			m.IpInHdrErrors = value
		case "InAddrErrors":
			m.IpInAddrErrors = value
		case "ForwDatagrams":
			m.IpForwDatagrams = value
		case "InUnknownProtos":
			m.IpInUnknownProtos = value
		case "InDiscards":
			m.IpInDiscards = value
		case "InDelivers":
			m.IpInDelivers = value
		case "OutRequests":
			m.IpOutRequests = value
		case "OutDiscards":
			m.IpOutDiscards = value
		case "OutNoRoutes":
			m.IpOutNoRoutes = value
		case "ReasmTimeout":
			m.IpReasmTimeout = value
		case "ReasmReqds":
			m.IpReasmReqds = value
		case "ReasmOKs":
			m.IpReasmOks = value
		case "ReasmFails":
			m.IpReasmFails = value
		case "FragOKs":
			m.IpFragOks = value
		case "FragFails":
			m.IpFragFails = value
		case "FragCreates":
			m.IpFragCreates = value

		// tcp
		case "RtoAlgorithm":
			m.TcpRtoAlgorithm = value
		case "RtoMin":
			m.TcpRtoMin = value
		case "RtoMax":
			m.TcpRtoMax = value
		case "ActiveOpens":
			m.TcpActiveOpens = value
		case "PassiveOpens":
			m.TcpPassiveOpens = value
		case "AttemptFails":
			m.TcpAttemptFails = value
		case "EstabResets":
			m.TcpEstabResets = value
		case "CurrEstab":
			m.TcpCurrEstab = value
		case "InSegs":
			m.TcpInSegs = value
		case "OutSegs":
			m.TcpOutSegs = value
		case "RetransSegs":
			m.TcpRetransSegs = value
		case "InErrs":
			m.TcpInErrs = value
		case "OutRsts":
			m.TcpOutRsts = value

		// udp
		case "InDatagrams":
			m.UdpInDatagrams = value
		case "NoPorts":
			m.UdpNoPorts = value
		case "OutDatagrams":
			m.UdpOutDatagrams = value
		case "RcvbufErrors":
			m.UdpRcvbufErrors = value
		case "SndbufErrors":
			m.UdpSndbufErrors = value
		case "IgnoredMulti":
			m.UdpIgnoredMulti = value
		case "MemErrors":
			m.UdpMemErrors = value

		default:
			log.Fatal().Str("field", fields[i]).Msg("unrecognizable snmp field")
		}
	}

	return nil
}
