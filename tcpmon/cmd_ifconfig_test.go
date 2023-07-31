package tcpmon_test

import (
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	. "github.com/zperf/tcpmon/tcpmon"
)

func (s *CommandParserTestSuite) TestParseIfconfig() {
	assert := s.Assert()

	out := `docker0: flags=4099<UP,BROADCAST,MULTICAST>  mtu 1500
        inet 172.17.0.1  netmask 255.255.0.0  broadcast 172.17.255.255
        ether 02:42:25:39:c8:a8  txqueuelen 0  (Ethernet)
        RX packets 0  bytes 0 (0.0 B)
        RX errors 1  dropped 2  overruns 3  frame 4
        TX packets 0  bytes 0 (0.0 B)
        TX errors 5  dropped 6 overruns 7  carrier 8  collisions 9

ens192: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 192.168.126.128  netmask 255.255.255.0  broadcast 192.168.126.255
        inet6 fe80::4491:9dfc:c5cb:df80  prefixlen 64  scopeid 0x20<link>
        ether 00:0c:29:60:55:22  txqueuelen 1000  (Ethernet)
        RX packets 340574  bytes 482919210 (460.5 MiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 51720  bytes 3840693 (3.6 MiB)
        TX errors 0  dropped 99999 overruns 0  carrier 0  collisions 0

lo: flags=73<UP,LOOPBACK,RUNNING>  mtu 65536
        inet 127.0.0.1  netmask 255.0.0.0
        inet6 ::1  prefixlen 128  scopeid 0x10<host>
        loop  txqueuelen 1000  (Local Loopback)
        RX packets 5135  bytes 653744 (638.4 KiB)
        RX errors 0  dropped 0  overruns 10000  frame 0
        TX packets 5135  bytes 653744 (638.4 KiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0`
	lines := strings.FieldsFunc(out, func(r rune) bool {
		return r == '\n'
	})

	var nic NicMetric
	nic.Type = MetricType_NIC
	nic.Timestamp = timestamppb.New(time.Now())
	ParseIfconfigOutput(&nic, lines)

	assert.Equal(uint32(7), nic.Ifaces[0].TXOverruns)
	assert.Equal(uint32(99999), nic.Ifaces[1].TXDropped)
	assert.Equal(uint32(10000), nic.Ifaces[2].RXOverruns)
}
