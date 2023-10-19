package parsing

import (
	"os"
	"strings"
	"time"

	"github.com/samber/lo"

	. "github.com/zperf/tcpmon/tcpmon/parsing"
	. "github.com/zperf/tcpmon/tcpmon/tproto"
	. "github.com/zperf/tcpmon/tcpmon/tutils"
)

func (s *ParsingTestSuite) TestParseSS() {
	buf, err := os.ReadFile("ss.txt")
	s.Require().NoError(err)

	lines := strings.FieldsFunc(string(buf), SplitNewline)

	var t TcpMetric
	t.Timestamp = time.Now().Unix()
	t.Type = MetricType_TCP
	ParseSS(&t, lines)

	// check result count
	c := lo.CountBy(t.Sockets, func(item *SocketMetric) bool {
		return item.GetState() == SocketState_TCP_LISTEN
	})
	s.Assert().Equal(1905, c)

	c = lo.CountBy(t.Sockets, func(item *SocketMetric) bool {
		return item.GetState() == SocketState_TCP_CLOSE_WAIT
	})
	s.Assert().Equal(59, c)

	c = lo.CountBy(t.Sockets, func(item *SocketMetric) bool {
		return item.GetState() == SocketState_TCP_ESTABLISHED
	})
	s.Assert().Equal(2141, c)

	c = lo.CountBy(t.Sockets, func(item *SocketMetric) bool {
		return item.GetState() == SocketState_TCP_FIN_WAIT2
	})
	s.Assert().Equal(1, c)

	c = lo.CountBy(t.Sockets, func(item *SocketMetric) bool {
		return item.GetState() == SocketState_TCP_TIME_WAIT
	})
	s.Assert().Equal(1434, c)

	zoos := lo.Filter(t.Sockets, func(m *SocketMetric, i int) bool {
		_ = i
		_, ok := lo.Find(m.Processes, func(p *ProcessInfo) bool {
			return p.Name == "java"
		})
		return ok
	})
	s.Assert().Equal(118, len(zoos))

	zoo, ok := lo.Find(t.Sockets, func(m *SocketMetric) bool {
		_, ok := lo.Find(m.Processes, func(p *ProcessInfo) bool {
			return p.Fd == 154 && p.Name == "java"
		})
		return ok
	})
	s.Assert().True(ok)

	s.Assert().Equal(SocketState_TCP_ESTABLISHED, zoo.State)
	s.Assert().Equal(int64(0), zoo.SendQ)
	s.Assert().Equal(uint32(0), zoo.RecvQ)
	s.Assert().Equal("::ffff:10.255.0.141:2181", zoo.LocalAddr)
	s.Assert().Equal("::ffff:10.255.0.102:35648", zoo.PeerAddr)
	s.Assert().Equal(uint64(169202297), zoo.BytesAcked)
}
