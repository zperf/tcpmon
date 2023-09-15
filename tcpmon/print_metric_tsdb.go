package tcpmon

import (
	"fmt"
)

func boolToUint32(x bool) uint32 {
	if !x {
		return 0
	} else {
		return 1
	}
}

type TSDBMetricPrinter struct {}

func (tsdb TSDBMetricPrinter) PrintTcpMetric(m *TcpMetric, hostname string) {
	for _, socket := range m.GetSockets() {
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s State=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetState(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s RecvQ=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetRecvQ(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s SendQ=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSendQ(), m.GetTimestamp())

		for _, process := range socket.GetProcesses() {
			fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,ProcessName=%s,hostname=%s Pid=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), process.GetName(), hostname, process.GetPid(), m.GetTimestamp())
			fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,ProcessName=%s,hostname=%s Fd=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), process.GetName(), hostname, process.GetFd(), m.GetTimestamp())
		}

		for _, timer := range socket.GetTimers() {
			fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,TimerName=%s,hostname=%s ExpireTimeUs=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), timer.GetName(), hostname, timer.GetExpireTimeUs(), m.GetTimestamp())
			fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,TimerName=%s,hostname=%s Retrans=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), timer.GetName(), hostname, timer.GetRetrans(), m.GetTimestamp())
		}

		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s RmemAlloc=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSkmem().GetRmemAlloc(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s RcvBuf=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSkmem().GetRcvBuf(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s WmemAlloc=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSkmem().GetWmemAlloc(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s SndBuf=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSkmem().GetSndBuf(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s FwdAlloc=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSkmem().GetFwdAlloc(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s WmemQueued=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSkmem().GetWmemQueued(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s OptMem=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSkmem().GetOptMem(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s BackLog=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSkmem().GetBackLog(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s SockDrop=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSkmem().GetSockDrop(), m.GetTimestamp())

		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Ts=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, boolToUint32(socket.GetTs()), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Sack=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, boolToUint32(socket.GetSack()), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Cubic=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, boolToUint32(socket.GetCubic()), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s AppLimited=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, boolToUint32(socket.GetAppLimited()), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s PacingRate=%f %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetPacingRate(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s DeliveryRate=%f %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetDeliveryRate(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Send=%f %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSend(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s SndWscale=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSndWscale(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s RcvWscale=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetRcvWscale(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Rto=%f %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetRto(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Rtt=%f %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetRtt(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Rttvar=%f %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetRttvar(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Minrtt=%f %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetMinrtt(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s RcvRtt=%f %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetRcvRtt(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s RetransNow=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetRetransNow(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s RetransTotal=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetRetransTotal(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Ato=%f %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetAto(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Mss=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetMss(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Pmtu=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetPmtu(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Rcvmss=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetRcvmss(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Advmss=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetAdvmss(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Cwnd=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetCwnd(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s SndWnd=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSndWnd(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s BytesSent=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetBytesSent(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s BytesAcked=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetBytesAcked(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s BytesReceived=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetBytesReceived(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s SegsOut=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSegsOut(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s SegsIn=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSegsIn(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Lastsnd=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetLastsnd(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Lastrcv=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetLastrcv(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Lastack=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetLastack(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Delivered=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetDelivered(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s BusyMs=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetBusyMs(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s RcvSpace=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetRcvSpace(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s RcvSsthresh=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetRcvSsthresh(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s DataSegsOut=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetDataSegsOut(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s DataSegsIn=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetDataSegsIn(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s RwndLimited=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetRwndLimited(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s SndbufLimited=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSndbufLimited(), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Ecn=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, boolToUint32(socket.GetEcn()), m.GetTimestamp())
		fmt.Printf("tcp,LocalAddr=%s,PeerAddr=%s,hostname=%s Ecnseen=%d %d\n", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, boolToUint32(socket.GetEcnseen()), m.GetTimestamp())
	}
}

func (tsdb TSDBMetricPrinter) PrintNicMetric(m *NicMetric, hostname string) {
	for _, iface := range m.GetIfaces() {
		fmt.Printf("nic,name=%s,hostname=%s RxErrors=%d %d\n", iface.GetName(), hostname, iface.GetRxErrors(), m.GetTimestamp())
		fmt.Printf("nic,name=%s,hostname=%s RxDropped=%d %d\n", iface.GetName(), hostname, iface.GetRxDropped(), m.GetTimestamp())
		fmt.Printf("nic,name=%s,hostname=%s RxOverruns=%d %d\n", iface.GetName(), hostname, iface.GetRxOverruns(), m.GetTimestamp())
		fmt.Printf("nic,name=%s,hostname=%s RxFrame=%d %d\n", iface.GetName(), hostname, iface.GetRxFrame(), m.GetTimestamp())
		fmt.Printf("nic,name=%s,hostname=%s TxErrors=%d %d\n", iface.GetName(), hostname, iface.GetTxErrors(), m.GetTimestamp())
		fmt.Printf("nic,name=%s,hostname=%s TxDropped=%d %d\n", iface.GetName(), hostname, iface.GetTxDropped(), m.GetTimestamp())
		fmt.Printf("nic,name=%s,hostname=%s TxOverruns=%d %d\n", iface.GetName(), hostname, iface.GetTxOverruns(), m.GetTimestamp())
		fmt.Printf("nic,name=%s,hostname=%s TxCarrier=%d %d\n", iface.GetName(), hostname, iface.GetTxCarrier(), m.GetTimestamp())
		fmt.Printf("nic,name=%s,hostname=%s TxCollisions=%d %d\n", iface.GetName(), hostname, iface.GetTxCollisions(), m.GetTimestamp())
	}
}

func (tsdb TSDBMetricPrinter) PrintNetstatMetric(m *NetstatMetric, hostname string) {
	fmt.Printf("net,hostname=%s IpTotalPacketsReceived=%d %d\n", hostname, m.GetIpTotalPacketsReceived(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s IpForwarded=%d %d\n", hostname, m.GetIpForwarded(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s IpIncomingPacketsDiscarded=%d %d\n", hostname, m.GetIpIncomingPacketsDiscarded(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s IpIncomingPacketsDelivered=%d %d\n", hostname, m.GetIpIncomingPacketsDelivered(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s IpRequestsSentOut=%d %d\n", hostname, m.GetIpRequestsSentOut(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s IpOutgoingPacketsDropped=%d %d\n", hostname, m.GetIpOutgoingPacketsDropped(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s TcpActiveConnectionsOpenings=%d %d\n", hostname, m.GetTcpActiveConnectionsOpenings(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s TcpPassiveConnectionOpenings=%d %d\n", hostname, m.GetTcpPassiveConnectionOpenings(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s TcpFailedConnectionAttempts=%d %d\n", hostname, m.GetTcpFailedConnectionAttempts(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s TcpConnectionResetsReceived=%d %d\n", hostname, m.GetTcpConnectionResetsReceived(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s TcpConnectionsEstablished=%d %d\n", hostname, m.GetTcpConnectionsEstablished(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s TcpSegmentsReceived=%d %d\n", hostname, m.GetTcpSegmentsReceived(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s TcpSegmentsSendOut=%d %d\n", hostname, m.GetTcpSegmentsSendOut(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s TcpSegmentsRetransmitted=%d %d\n", hostname, m.GetTcpSegmentsRetransmitted(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s TcpBadSegmentsReceived=%d %d\n", hostname, m.GetTcpBadSegmentsReceived(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s TcpResetsSent=%d %d\n", hostname, m.GetTcpResetsSent(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s UdpPacketsReceived=%d %d\n", hostname, m.GetUdpPacketsReceived(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s UdpPacketsToUnknownPortReceived=%d %d\n", hostname, m.GetUdpPacketsToUnknownPortReceived(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s UdpPacketReceiveErrors=%d %d\n", hostname, m.GetUdpPacketReceiveErrors(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s UdpPacketsSent=%d %d\n", hostname, m.GetUdpPacketsSent(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s UdpReceiveBufferErrors=%d %d\n", hostname, m.GetUdpReceiveBufferErrors(), m.GetTimestamp())
	fmt.Printf("net,hostname=%s UdpSendBufferErrors=%d %d\n", hostname, m.GetUdpSendBufferErrors(), m.GetTimestamp())
}
