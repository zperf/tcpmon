package influxdb

import (
	"fmt"
	"strings"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"github.com/zperf/tcpmon/tcpmon/gproto"
)

type MetricConv struct {
	Hostname string
}

func NewMetricConv(hostname string) *MetricConv {
	return &MetricConv{Hostname: hostname}
}

func (c *MetricConv) Metric(metric *gproto.Metric) (int64, []*write.Point) {
	switch m := metric.Body.(type) {
	case *gproto.Metric_Tcp:
		return m.Tcp.Timestamp, c.Tcp(m.Tcp)
	case *gproto.Metric_Net:
		return m.Net.Timestamp, c.Net(m.Net)
	case *gproto.Metric_Nic:
		return m.Nic.Timestamp, c.Nic(m.Nic)
	default:
		log.Fatal().Msg("Unknown metric type")
	}
	return 0, nil
}

func (c *MetricConv) Nic(metric *gproto.NicMetric) []*write.Point {
	ts := time.Unix(metric.GetTimestamp(), 0)
	points := make([]*write.Point, 0)

	for _, iface := range metric.Ifaces {
		p := write.NewPoint("nic",
			map[string]string{"Hostname": c.Hostname, "Name": iface.Name},
			map[string]interface{}{"RxErrors": iface.RxErrors},
			ts)
		points = append(points, p)
		p = write.NewPoint("nic",
			map[string]string{"Hostname": c.Hostname, "Name": iface.Name},
			map[string]interface{}{"RxDropped": iface.RxDropped},
			ts)
		points = append(points, p)
		p = write.NewPoint("nic",
			map[string]string{"Hostname": c.Hostname, "Name": iface.Name},
			map[string]interface{}{"RxOverruns": iface.RxOverruns},
			ts)
		points = append(points, p)
		p = write.NewPoint("nic",
			map[string]string{"Hostname": c.Hostname, "Name": iface.Name},
			map[string]interface{}{"RxFrame": iface.RxFrame},
			ts)
		points = append(points, p)
		p = write.NewPoint("nic",
			map[string]string{"Hostname": c.Hostname, "Name": iface.Name},
			map[string]interface{}{"TxErrors": iface.TxErrors},
			ts)
		points = append(points, p)
		p = write.NewPoint("nic",
			map[string]string{"Hostname": c.Hostname, "Name": iface.Name},
			map[string]interface{}{"TxDropped": iface.TxDropped},
			ts)
		points = append(points, p)
		p = write.NewPoint("nic",
			map[string]string{"Hostname": c.Hostname, "Name": iface.Name},
			map[string]interface{}{"TxOverruns": iface.TxOverruns},
			ts)
		points = append(points, p)
		p = write.NewPoint("nic",
			map[string]string{"Hostname": c.Hostname, "Name": iface.Name},
			map[string]interface{}{"TxCarrier": iface.TxCarrier},
			ts)
		points = append(points, p)
		p = write.NewPoint("nic",
			map[string]string{"Hostname": c.Hostname, "Name": iface.Name},
			map[string]interface{}{"TxCollisions": iface.TxCollisions},
			ts)
		points = append(points, p)

		points = append(points, p)
	}

	return points
}

func getPort(addr string) string {
	p := strings.LastIndex(addr, ":")
	if p == -1 {
		return ""
	}
	return addr[p+1:]
}

// b2i: bool to int
func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (c *MetricConv) Tcp(metric *gproto.TcpMetric) []*write.Point {
	ts := time.Unix(metric.GetTimestamp(), 0)
	points := make([]*write.Point, 0)

	for _, s := range metric.GetSockets() {
		tags := map[string]string{
			"LocalAddr": s.LocalAddr,
			"LocalPort": getPort(s.LocalAddr),
			"PeerAddr":  s.PeerAddr,
			"PeerPort":  getPort(s.PeerAddr),
			"Hostname":  c.Hostname,
		}

		processText := strings.Join(lo.Map(s.GetProcesses(), func(p *gproto.ProcessInfo, _ int) string {
			return fmt.Sprintf("cmd:%s;pid:%v;fd:%v", p.GetName(), p.GetPid(), p.GetFd())
		}), ";")
		if processText != "" {
			tags["Process"] = processText
			tags["ProcessName"] = s.GetProcesses()[0].GetName()
		}

		for _, timer := range s.GetTimers() {
			p := write.NewPoint("tcp", tags,
				map[string]interface{}{
					"Timer":        timer.GetName(),
					"ExpireTimeUs": timer.GetExpireTimeUs(),
					"Retrans":      timer.GetRetrans(),
				}, ts)
			points = append(points, p)
		}

		p := write.NewPoint("tcp", tags,
			map[string]interface{}{"State": s.State.String()},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"RecvQ": s.RecvQ},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"SendQ": s.SendQ},
			ts)
		points = append(points, p)

		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"RmemAlloc": s.GetSkmem().GetRmemAlloc()},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"RcvBuf": s.GetSkmem().GetRcvBuf()},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"WmemAlloc": s.GetSkmem().GetWmemAlloc()},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"SndBuf": s.GetSkmem().GetSndBuf()},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"FwdAlloc": s.GetSkmem().GetFwdAlloc()},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"WmemQueued": s.GetSkmem().GetWmemQueued()},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"OptMem": s.GetSkmem().GetOptMem()},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"BackLog": s.GetSkmem().GetBackLog()},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"SockDrop": s.GetSkmem().GetSockDrop()},
			ts)
		points = append(points, p)

		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Ts": b2i(s.Ts)},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Sack": b2i(s.Sack)},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Cubic": b2i(s.Cubic)},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"AppLimited": b2i(s.AppLimited)},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"PacingRate": s.PacingRate},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"DeliveryRate": s.DeliveryRate},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Send": s.Send},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"SndWscale": s.SndWscale},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"RcvWscale": s.RcvWscale},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Rto": s.Rto},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Rtt": s.Rtt},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Rttvar": s.Rttvar},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Minrtt": s.Minrtt},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"RcvRtt": s.RcvRtt},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"RetransNow": s.RetransNow},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"RetransTotal": s.RetransTotal},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Ato": s.Ato},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Mss": s.Mss},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Pmtu": s.Pmtu},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Rcvmss": s.Rcvmss},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Advmss": s.Advmss},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Cwnd": s.Cwnd},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"SndWnd": s.SndWnd},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"BytesSent": s.BytesSent},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"BytesAcked": s.BytesAcked},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"BytesReceived": s.BytesReceived},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"SegsOut": s.SegsOut},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"SegsIn": s.SegsIn},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Lastsnd": s.Lastsnd},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Lastrcv": s.Lastrcv},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Lastack": s.Lastack},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Delivered": s.Delivered},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"BusyMs": s.BusyMs},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"RcvSpace": s.RcvSpace},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"RcvSsthresh": s.RcvSsthresh},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"DataSegsOut": s.DataSegsOut},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"DataSegsIn": s.DataSegsIn},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"RwndLimited": s.RwndLimited},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"SndbufLimited": s.SndbufLimited},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Ecn": b2i(s.Ecn)},
			ts)
		points = append(points, p)
		p = write.NewPoint("tcp", tags,
			map[string]interface{}{"Ecnseen": b2i(s.Ecnseen)},
			ts)
		points = append(points, p)
	}

	return points
}

func (c *MetricConv) Net(metric *gproto.NetstatMetric) []*write.Point {
	ts := time.Unix(metric.GetTimestamp(), 0)
	points := make([]*write.Point, 0)

	p := write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpForwarding": metric.IpForwarding},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpDefaultTtl": metric.IpDefaultTtl},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInReceives": metric.IpInReceives},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInHdrErrors": metric.IpInHdrErrors},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInAddrErrors": metric.IpInAddrErrors},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpForwDatagrams": metric.IpForwDatagrams},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInUnknownProtos": metric.IpInUnknownProtos},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInDiscards": metric.IpInDiscards},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInDelivers": metric.IpInDelivers},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpOutRequests": metric.IpOutRequests},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpOutDiscards": metric.IpOutDiscards},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpOutNoRoutes": metric.IpOutNoRoutes},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpReasmTimeout": metric.IpReasmTimeout},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpReasmReqds": metric.IpReasmReqds},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpReasmOks": metric.IpReasmOks},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpReasmFails": metric.IpReasmFails},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpFragOks": metric.IpFragOks},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpFragFails": metric.IpFragFails},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpFragCreates": metric.IpFragCreates},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInNoRoutes": metric.IpInNoRoutes},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInTruncatedPkts": metric.IpInTruncatedPkts},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInMcastPkts": metric.IpInMcastPkts},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpOutMcastPkts": metric.IpOutMcastPkts},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInBcastPkts": metric.IpInBcastPkts},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpOutBcastPkts": metric.IpOutBcastPkts},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInOctets": metric.IpInOctets},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpOutOctets": metric.IpOutOctets},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInMcastOctets": metric.IpInMcastOctets},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpOutMcastOctets": metric.IpOutMcastOctets},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInBcastOctets": metric.IpInBcastOctets},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpOutBcastOctets": metric.IpOutBcastOctets},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInCsumErrors": metric.IpInCsumErrors},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInNoEctPkts": metric.IpInNoEctPkts},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInEct1Pkts": metric.IpInEct1Pkts},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInEct0Pkts": metric.IpInEct0Pkts},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpInCePkts": metric.IpInCePkts},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IpReasmOverlaps": metric.IpReasmOverlaps},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"UdpInDatagrams": metric.UdpInDatagrams},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"UdpNoPorts": metric.UdpNoPorts},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"UdpInErrors": metric.UdpInErrors},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"UdpOutDatagrams": metric.UdpOutDatagrams},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"UdpRcvbufErrors": metric.UdpRcvbufErrors},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"UdpSndbufErrors": metric.UdpSndbufErrors},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"UdpInCsumErrors": metric.UdpInCsumErrors},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"UdpIgnoredMulti": metric.UdpIgnoredMulti},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"UdpMemErrors": metric.UdpMemErrors},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpRtoAlgorithm": metric.TcpRtoAlgorithm},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpRtoMin": metric.TcpRtoMin},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpRtoMax": metric.TcpRtoMax},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpMaxConn": metric.TcpMaxConn},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpActiveOpens": metric.TcpActiveOpens},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpPassiveOpens": metric.TcpPassiveOpens},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpAttemptFails": metric.TcpAttemptFails},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpEstabResets": metric.TcpEstabResets},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpCurrEstab": metric.TcpCurrEstab},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpInSegs": metric.TcpInSegs},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpOutSegs": metric.TcpOutSegs},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpRetransSegs": metric.TcpRetransSegs},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpInErrs": metric.TcpInErrs},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpOutRsts": metric.TcpOutRsts},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpInCsumErrors": metric.TcpInCsumErrors},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSyncookiesSent": metric.TcpSyncookiesSent},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSyncookiesRecv": metric.TcpSyncookiesRecv},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSyncookiesFailed": metric.TcpSyncookiesFailed},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpEmbryonicRsts": metric.TcpEmbryonicRsts},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpPruneCalled": metric.TcpPruneCalled},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpRcvPruned": metric.TcpRcvPruned},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpOfoPruned": metric.TcpOfoPruned},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpOutOfWindowIcmps": metric.TcpOutOfWindowIcmps},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpLockDroppedIcmps": metric.TcpLockDroppedIcmps},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpArpFilter": metric.TcpArpFilter},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpTw": metric.TcpTw},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpTwRecycled": metric.TcpTwRecycled},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpTwKilled": metric.TcpTwKilled},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpPawsActive": metric.TcpPawsActive},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpPawsEstab": metric.TcpPawsEstab},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDelayedAcks": metric.TcpDelayedAcks},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDelayedAckLocked": metric.TcpDelayedAckLocked},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDelayedAckLost": metric.TcpDelayedAckLost},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpListenOverflows": metric.TcpListenOverflows},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpListenDrops": metric.TcpListenDrops},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpHpHits": metric.TcpHpHits},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpPureAcks": metric.TcpPureAcks},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpHpAcks": metric.TcpHpAcks},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpRenoRecovery": metric.TcpRenoRecovery},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSackRecovery": metric.TcpSackRecovery},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSackReneging": metric.TcpSackReneging},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSackReorder": metric.TcpSackReorder},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpRenoReorder": metric.TcpRenoReorder},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpTsReorder": metric.TcpTsReorder},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpFullUndo": metric.TcpFullUndo},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpPartialUndo": metric.TcpPartialUndo},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDsackUndo": metric.TcpDsackUndo},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpLossUndo": metric.TcpLossUndo},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpLostRetransmit": metric.TcpLostRetransmit},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpRenoFailures": metric.TcpRenoFailures},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSackFailures": metric.TcpSackFailures},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpLossFailures": metric.TcpLossFailures},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpFastRetrans": metric.TcpFastRetrans},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSlowStartRetrans": metric.TcpSlowStartRetrans},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpTimeouts": metric.TcpTimeouts},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpLossProbes": metric.TcpLossProbes},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpLossProbeRecovery": metric.TcpLossProbeRecovery},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpRenoRecoveryFail": metric.TcpRenoRecoveryFail},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSackRecoveryFail": metric.TcpSackRecoveryFail},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpRcvCollapsed": metric.TcpRcvCollapsed},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpBacklogCoalesce": metric.TcpBacklogCoalesce},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDsackOldSent": metric.TcpDsackOldSent},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDsackOfoSent": metric.TcpDsackOfoSent},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDsackRecv": metric.TcpDsackRecv},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDsackOfoRecv": metric.TcpDsackOfoRecv},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpAbortOnData": metric.TcpAbortOnData},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpAbortOnClose": metric.TcpAbortOnClose},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpAbortOnMemory": metric.TcpAbortOnMemory},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpAbortOnTimeout": metric.TcpAbortOnTimeout},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpAbortOnLinger": metric.TcpAbortOnLinger},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpAbortFailed": metric.TcpAbortFailed},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpMemoryPressures": metric.TcpMemoryPressures},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpMemoryPressuresChrono": metric.TcpMemoryPressuresChrono},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSackDiscard": metric.TcpSackDiscard},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDsackIgnoredOld": metric.TcpDsackIgnoredOld},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDsackIgnoredNoUndo": metric.TcpDsackIgnoredNoUndo},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSpuriousRtos": metric.TcpSpuriousRtos},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpMd5NotFound": metric.TcpMd5NotFound},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpMd5Unexpected": metric.TcpMd5Unexpected},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpMd5Failure": metric.TcpMd5Failure},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSackShifted": metric.TcpSackShifted},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSackMerged": metric.TcpSackMerged},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSackShiftFallback": metric.TcpSackShiftFallback},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpBacklogDrop": metric.TcpBacklogDrop},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpPfMemallocDrop": metric.TcpPfMemallocDrop},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpMinTtlDrop": metric.TcpMinTtlDrop},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDeferAcceptDrop": metric.TcpDeferAcceptDrop},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpIpReversePathFilter": metric.TcpIpReversePathFilter},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpTimeWaitOverflow": metric.TcpTimeWaitOverflow},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpReqQFullDoCookies": metric.TcpReqQFullDoCookies},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpReqQFullDrop": metric.TcpReqQFullDrop},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpRetransFail": metric.TcpRetransFail},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpRcvCoalesce": metric.TcpRcvCoalesce},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpOfoQueue": metric.TcpOfoQueue},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpOfoDrop": metric.TcpOfoDrop},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpOfoMerge": metric.TcpOfoMerge},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpChallengeAck": metric.TcpChallengeAck},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSynChallenge": metric.TcpSynChallenge},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpFastOpenActive": metric.TcpFastOpenActive},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpFastOpenActiveFail": metric.TcpFastOpenActiveFail},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpFastOpenPassive": metric.TcpFastOpenPassive},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpFastOpenPassiveFail": metric.TcpFastOpenPassiveFail},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpFastOpenListenOverflow": metric.TcpFastOpenListenOverflow},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpFastOpenCookieReqd": metric.TcpFastOpenCookieReqd},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpFastOpenBlackhole": metric.TcpFastOpenBlackhole},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSpuriousRtxHostQueues": metric.TcpSpuriousRtxHostQueues},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpBusyPollRxPackets": metric.TcpBusyPollRxPackets},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpAutoCorking": metric.TcpAutoCorking},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpFromZeroWindowAdv": metric.TcpFromZeroWindowAdv},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpToZeroWindowAdv": metric.TcpToZeroWindowAdv},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpWantZeroWindowAdv": metric.TcpWantZeroWindowAdv},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpSynRetrans": metric.TcpSynRetrans},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpOrigDataSent": metric.TcpOrigDataSent},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpHystartTrainDetect": metric.TcpHystartTrainDetect},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpHystartTrainCwnd": metric.TcpHystartTrainCwnd},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpHystartDelayDetect": metric.TcpHystartDelayDetect},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpHystartDelayCwnd": metric.TcpHystartDelayCwnd},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpAckSkippedSynRecv": metric.TcpAckSkippedSynRecv},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpAckSkippedPaws": metric.TcpAckSkippedPaws},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpAckSkippedSeq": metric.TcpAckSkippedSeq},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpAckSkippedFinWait2": metric.TcpAckSkippedFinWait2},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpAckSkippedTimeWait": metric.TcpAckSkippedTimeWait},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpAckSkippedChallenge": metric.TcpAckSkippedChallenge},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpWinProbe": metric.TcpWinProbe},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpKeepAlive": metric.TcpKeepAlive},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpMtupFail": metric.TcpMtupFail},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpMtupSuccess": metric.TcpMtupSuccess},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDelivered": metric.TcpDelivered},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDeliveredCe": metric.TcpDeliveredCe},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpAckCompressed": metric.TcpAckCompressed},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpZeroWindowDrop": metric.TcpZeroWindowDrop},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpRcvQDrop": metric.TcpRcvQDrop},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpWqueueTooBig": metric.TcpWqueueTooBig},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpFastOpenPassiveAltKey": metric.TcpFastOpenPassiveAltKey},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpTimeoutRehash": metric.TcpTimeoutRehash},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDuplicateDataRehash": metric.TcpDuplicateDataRehash},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDsackRecvSegs": metric.TcpDsackRecvSegs},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpDsackIgnoredDubious": metric.TcpDsackIgnoredDubious},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpMigrateReqSuccess": metric.TcpMigrateReqSuccess},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpMigrateReqFailure": metric.TcpMigrateReqFailure},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"TcpPlbRehash": metric.TcpPlbRehash},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpInMsgs": metric.IcmpInMsgs},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpInErrors": metric.IcmpInErrors},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpInCsumErrors": metric.IcmpInCsumErrors},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpInDestUnreachs": metric.IcmpInDestUnreachs},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpInTimeExcds": metric.IcmpInTimeExcds},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpInParmProbs": metric.IcmpInParmProbs},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpInSrcQuenchs": metric.IcmpInSrcQuenchs},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpInRedirects": metric.IcmpInRedirects},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpInEchos": metric.IcmpInEchos},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpInEchoReps": metric.IcmpInEchoReps},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpInTimestamps": metric.IcmpInTimestamps},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpInTimestampReps": metric.IcmpInTimestampReps},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpInAddrMasks": metric.IcmpInAddrMasks},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpInAddrMaskReps": metric.IcmpInAddrMaskReps},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpOutMsgs": metric.IcmpOutMsgs},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpOutErrors": metric.IcmpOutErrors},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpOutRateLimitGlobal": metric.IcmpOutRateLimitGlobal},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpOutRateLimitHost": metric.IcmpOutRateLimitHost},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpOutDestUnreachs": metric.IcmpOutDestUnreachs},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpOutTimeExcds": metric.IcmpOutTimeExcds},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpOutParmProbs": metric.IcmpOutParmProbs},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpOutSrcQuenchs": metric.IcmpOutSrcQuenchs},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpOutRedirects": metric.IcmpOutRedirects},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpOutEchos": metric.IcmpOutEchos},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpOutEchoReps": metric.IcmpOutEchoReps},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpOutTimestamps": metric.IcmpOutTimestamps},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpOutTimestampReps": metric.IcmpOutTimestampReps},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpOutAddrMasks": metric.IcmpOutAddrMasks},
		ts)
	points = append(points, p)
	p = write.NewPoint("net",
		map[string]string{"Hostname": c.Hostname},
		map[string]interface{}{"IcmpOutAddrMaskReps": metric.IcmpOutAddrMaskReps},
		ts)
	points = append(points, p)

	return points
}
