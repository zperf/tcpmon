package parsing

import (
	"os"
	"strings"

	"github.com/zperf/tcpmon/tcpmon/parsing"
	"github.com/zperf/tcpmon/tcpmon/tproto"
)

func (s *ParsingTestSuite) TestParseSnmp() {
	buf, err := os.ReadFile("snmp.txt")
	s.Require().NoError(err)

	snmp := string(buf)
	var m tproto.NetstatMetric
	err = parsing.ParseSnmp(strings.NewReader(snmp), &m)
	s.Require().NoError(err)

	s.Require().Equal(uint64(1), m.IpForwarding)
	s.Require().Equal(uint64(64), m.IpDefaultTtl)
	s.Require().Equal(uint64(338468), m.IpInReceives)
	s.Require().Equal(uint64(0), m.IpInHdrErrors)
	s.Require().Equal(uint64(0), m.IpInAddrErrors)
	s.Require().Equal(uint64(1), m.IpForwDatagrams)
	s.Require().Equal(uint64(0), m.IpInUnknownProtos)
	s.Require().Equal(uint64(0), m.IpInDiscards)
	s.Require().Equal(uint64(338379), m.IpInDelivers)
	s.Require().Equal(uint64(377770), m.IpOutRequests)
	s.Require().Equal(uint64(0), m.IpOutDiscards)
	s.Require().Equal(uint64(40), m.IpOutNoRoutes)
	s.Require().Equal(uint64(0), m.IpReasmTimeout)
	s.Require().Equal(uint64(0), m.IpReasmReqds)
	s.Require().Equal(uint64(0), m.IpReasmOks)
	s.Require().Equal(uint64(0), m.IpReasmFails)
	s.Require().Equal(uint64(0), m.IpFragOks)
	s.Require().Equal(uint64(0), m.IpFragFails)
	s.Require().Equal(uint64(0), m.IpFragCreates)
	s.Require().Equal(uint64(2956), m.IcmpInMsgs)
	s.Require().Equal(uint64(0), m.IcmpInErrors)
	s.Require().Equal(uint64(0), m.IcmpInCsumErrors)
	s.Require().Equal(uint64(2956), m.IcmpInDestUnreachs)
	s.Require().Equal(uint64(0), m.IcmpInTimeExcds)
	s.Require().Equal(uint64(0), m.IcmpInParmProbs)
	s.Require().Equal(uint64(0), m.IcmpInSrcQuenchs)
	s.Require().Equal(uint64(0), m.IcmpInRedirects)
	s.Require().Equal(uint64(0), m.IcmpInEchos)
	s.Require().Equal(uint64(0), m.IcmpInEchoReps)
	s.Require().Equal(uint64(0), m.IcmpInTimestamps)
	s.Require().Equal(uint64(0), m.IcmpInTimestampReps)
	s.Require().Equal(uint64(0), m.IcmpInAddrMasks)
	s.Require().Equal(uint64(0), m.IcmpInAddrMaskReps)
	s.Require().Equal(uint64(30), m.IcmpOutMsgs)
	s.Require().Equal(uint64(0), m.IcmpOutErrors)
	s.Require().Equal(uint64(0), m.IcmpOutRateLimitGlobal)
	s.Require().Equal(uint64(0), m.IcmpOutRateLimitHost)
	s.Require().Equal(uint64(30), m.IcmpOutDestUnreachs)
	s.Require().Equal(uint64(0), m.IcmpOutTimeExcds)
	s.Require().Equal(uint64(0), m.IcmpOutParmProbs)
	s.Require().Equal(uint64(0), m.IcmpOutSrcQuenchs)
	s.Require().Equal(uint64(0), m.IcmpOutRedirects)
	s.Require().Equal(uint64(0), m.IcmpOutEchos)
	s.Require().Equal(uint64(0), m.IcmpOutEchoReps)
	s.Require().Equal(uint64(0), m.IcmpOutTimestamps)
	s.Require().Equal(uint64(0), m.IcmpOutTimestampReps)
	s.Require().Equal(uint64(0), m.IcmpOutAddrMasks)
	s.Require().Equal(uint64(0), m.IcmpOutAddrMaskReps)
	s.Require().Equal(uint64(1), m.TcpRtoAlgorithm)
	s.Require().Equal(uint64(200), m.TcpRtoMin)
	s.Require().Equal(uint64(120000), m.TcpRtoMax)
	s.Require().Equal(int64(-1), m.TcpMaxConn)
	s.Require().Equal(uint64(4181), m.TcpActiveOpens)
	s.Require().Equal(uint64(52), m.TcpPassiveOpens)
	s.Require().Equal(uint64(3694), m.TcpAttemptFails)
	s.Require().Equal(uint64(10), m.TcpEstabResets)
	s.Require().Equal(uint64(22), m.TcpCurrEstab)
	s.Require().Equal(uint64(220096), m.TcpInSegs)
	s.Require().Equal(uint64(256252), m.TcpOutSegs)
	s.Require().Equal(uint64(1232), m.TcpRetransSegs)
	s.Require().Equal(uint64(15), m.TcpInErrs)
	s.Require().Equal(uint64(2426), m.TcpOutRsts)
	s.Require().Equal(uint64(0), m.TcpInCsumErrors)
	s.Require().Equal(uint64(114505), m.UdpInDatagrams)
	s.Require().Equal(uint64(30), m.UdpNoPorts)
	s.Require().Equal(uint64(0), m.UdpInErrors)
	s.Require().Equal(uint64(149416), m.UdpOutDatagrams)
	s.Require().Equal(uint64(0), m.UdpRcvbufErrors)
	s.Require().Equal(uint64(0), m.UdpSndbufErrors)
	s.Require().Equal(uint64(0), m.UdpInCsumErrors)
	s.Require().Equal(uint64(790), m.UdpIgnoredMulti)
	s.Require().Equal(uint64(0), m.UdpMemErrors)
}