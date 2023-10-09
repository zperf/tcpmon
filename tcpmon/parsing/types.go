package parsing

import (
	"strings"

	"github.com/rs/zerolog/log"
)

type ProcStatType = string

const (
	ProcNetStatIp    ProcStatType = "Ip"
	ProcNetStatIpExt ProcStatType = "IpExt"

	ProcNetStatTcp      ProcStatType = "Tcp"
	ProcNetStatTcpExt   ProcStatType = "TcpExt"
	ProcNetStatMPTcpExt ProcStatType = "MPTcpExt"

	ProcNetStatUdp     ProcStatType = "Udp"
	ProcNetStatUdpLite ProcStatType = "UdpLite"

	ProcNetStatIcmp    ProcStatType = "Icmp"
	ProcNetStatIcmpMsg ProcStatType = "IcmpMsg"
)

var ProcStatTypes = []ProcStatType{
	ProcNetStatIp, ProcNetStatIpExt, ProcNetStatTcp, ProcNetStatTcpExt, ProcNetStatMPTcpExt,
	ProcNetStatUdp, ProcNetStatUdpLite, ProcNetStatIcmp, ProcNetStatIcmpMsg,
}

func parseProcStatType(line string) ProcStatType {
	p := strings.Index(line, ": ")
	if p == -1 {
		log.Fatal().Str("line", line).Msg("Unrecognizable line")
	}
	t := line[:p]
	for _, s := range ProcStatTypes {
		if t == s {
			return s
		}
	}

	log.Fatal().Str("type", t).Msg("Unrecognizable type")
	return ""
}
