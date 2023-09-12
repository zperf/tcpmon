package tcpmon

import (
	"fmt"
	"strings"
)

func boolToUint32(x bool) uint32 {
	if !x {
		return 0
	} else {
		return 1
	}
}

func replaceStar(s string) string {
	s = strings.Replace(s, ":", "_", -1)
	s = strings.Replace(s, "*", "all", -1)
	return s
}

type TSDBMetricPrinter struct {}

func (tsdb TSDBMetricPrinter) PrintTcpMetric(m *TcpMetric, hostname string) {
	for _, socket := range m.GetSockets() {
		fmt.Printf("State %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetState(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("RecvQ %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetRecvQ(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("SendQ %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSendQ(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))

		for _, process := range socket.GetProcesses() {
			fmt.Printf("Pid %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s ProcessName=%s\n", m.GetTimestamp(), process.GetPid(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()), process.GetName())
			fmt.Printf("Fd %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s ProcessName=%s\n", m.GetTimestamp(), process.GetFd(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()), process.GetName())
		}

		for _, timer := range socket.GetTimers() {
			fmt.Printf("ExpireTimeUs %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s TimerName=%s\n", m.GetTimestamp(), timer.GetExpireTimeUs(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()), timer.GetName())
			fmt.Printf("Retrans %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s TimerName=%s\n", m.GetTimestamp(), timer.GetRetrans(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()), timer.GetName())
		}

		fmt.Printf("RmemAlloc %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSkmem().GetRmemAlloc(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("RcvBuf %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSkmem().GetRcvBuf(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("WmemAlloc %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSkmem().GetWmemAlloc(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("SndBuf %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSkmem().GetSndBuf(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("FwdAlloc %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSkmem().GetFwdAlloc(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("WmemQueued %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSkmem().GetWmemQueued(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("OptMem %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSkmem().GetOptMem(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("BackLog %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSkmem().GetBackLog(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("SockDrop %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSkmem().GetSockDrop(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))

		fmt.Printf("Ts %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), boolToUint32(socket.GetTs()), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Sack %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), boolToUint32(socket.GetSack()), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Cubic %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), boolToUint32(socket.GetCubic()), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("AppLimited %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), boolToUint32(socket.GetAppLimited()), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("PacingRate %d %f type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetPacingRate(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("DeliveryRate %d %f type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetDeliveryRate(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Send %d %f type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSend(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("SndWscale %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSndWscale(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("RcvWscale %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetRcvWscale(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Rto %d %f type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetRto(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Rtt %d %f type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetRtt(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Rttvar %d %f type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetRttvar(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Minrtt %d %f type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetMinrtt(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("RcvRtt %d %f type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetRcvRtt(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("RetransNow %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetRetransNow(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("RetransTotal %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetRetransTotal(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Ato %d %f type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetAto(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Mss %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetMss(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Pmtu %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetPmtu(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Rcvmss %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetRcvmss(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Advmss %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetAdvmss(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Cwnd %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetCwnd(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("SndWnd %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSndWnd(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("BytesSent %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetBytesSent(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("BytesAcked %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetBytesAcked(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("BytesReceived %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetBytesReceived(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("SegsOut %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSegsOut(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("SegsIn %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSegsIn(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Lastsnd %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetLastsnd(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Lastrcv %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetLastrcv(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Lastack %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetLastack(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Delivered %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetDelivered(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("BusyMs %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetBusyMs(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("RcvSpace %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetRcvSpace(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("RcvSsthresh %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetRcvSsthresh(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("DataSegsOut %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetDataSegsOut(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("DataSegsIn %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetDataSegsIn(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("RwndLimited %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetRwndLimited(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("SndbufLimited %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), socket.GetSndbufLimited(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Ecn %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), boolToUint32(socket.GetEcn()), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
		fmt.Printf("Ecnseen %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\n", m.GetTimestamp(), boolToUint32(socket.GetEcnseen()), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))
	}
}

func (tsdb TSDBMetricPrinter) PrintNicMetric(m *NicMetric, hostname string) {
	for _, iface := range m.GetIfaces() {
		fmt.Printf("RxErrors %d %d type=nic hostname=%s name=%s\n", m.GetTimestamp(), iface.GetRxErrors(), hostname, iface.GetName())
		fmt.Printf("RxDropped %d %d type=nic hostname=%s name=%s\n", m.GetTimestamp(), iface.GetRxDropped(), hostname, iface.GetName())
		fmt.Printf("RxOverruns %d %d type=nic hostname=%s name=%s\n", m.GetTimestamp(), iface.GetRxOverruns(), hostname, iface.GetName())
		fmt.Printf("RxFrame %d %d type=nic hostname=%s name=%s\n", m.GetTimestamp(), iface.GetRxFrame(), hostname, iface.GetName())
		fmt.Printf("TxErrors %d %d type=nic hostname=%s name=%s\n", m.GetTimestamp(), iface.GetTxErrors(), hostname, iface.GetName())
		fmt.Printf("TxDropped %d %d type=nic hostname=%s name=%s\n", m.GetTimestamp(), iface.GetTxDropped(), hostname, iface.GetName())
		fmt.Printf("TxOverruns %d %d type=nic hostname=%s name=%s\n", m.GetTimestamp(), iface.GetTxOverruns(), hostname, iface.GetName())
		fmt.Printf("TxCarrier %d %d type=nic hostname=%s name=%s\n", m.GetTimestamp(), iface.GetTxCarrier(), hostname, iface.GetName())
		fmt.Printf("TxCollisions %d %d type=nic hostname=%s name=%s\n", m.GetTimestamp(), iface.GetTxCollisions(), hostname, iface.GetName())
	}
}

func (tsdb TSDBMetricPrinter) PrintNetstatMetric(m *NetstatMetric, hostname string) {
	fmt.Printf("IpTotalPacketsReceived %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetIpTotalPacketsReceived(), hostname)
	fmt.Printf("IpForwarded %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetIpForwarded(), hostname)
	fmt.Printf("IpIncomingPacketsDiscarded %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetIpIncomingPacketsDiscarded(), hostname)
	fmt.Printf("IpIncomingPacketsDelivered %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetIpIncomingPacketsDelivered(), hostname)
	fmt.Printf("IpRequestsSentOut %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetIpRequestsSentOut(), hostname)
	fmt.Printf("IpOutgoingPacketsDropped %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetIpOutgoingPacketsDropped(), hostname)
	fmt.Printf("TcpActiveConnectionsOpenings %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetTcpActiveConnectionsOpenings(), hostname)
	fmt.Printf("TcpPassiveConnectionOpenings %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetTcpPassiveConnectionOpenings(), hostname)
	fmt.Printf("TcpFailedConnectionAttempts %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetTcpFailedConnectionAttempts(), hostname)
	fmt.Printf("TcpConnectionResetsReceived %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetTcpConnectionResetsReceived(), hostname)
	fmt.Printf("TcpConnectionsEstablished %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetTcpConnectionsEstablished(), hostname)
	fmt.Printf("TcpSegmentsReceived %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetTcpSegmentsReceived(), hostname)
	fmt.Printf("TcpSegmentsSendOut %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetTcpSegmentsSendOut(), hostname)
	fmt.Printf("TcpSegmentsRetransmitted %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetTcpSegmentsRetransmitted(), hostname)
	fmt.Printf("TcpBadSegmentsReceived %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetTcpBadSegmentsReceived(), hostname)
	fmt.Printf("TcpResetsSent %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetTcpResetsSent(), hostname)
	fmt.Printf("UdpPacketsReceived %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetUdpPacketsReceived(), hostname)
	fmt.Printf("UdpPacketsToUnknownPortReceived %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetUdpPacketsToUnknownPortReceived(), hostname)
	fmt.Printf("UdpPacketReceiveErrors %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetUdpPacketReceiveErrors(), hostname)
	fmt.Printf("UdpPacketsSent %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetUdpPacketsSent(), hostname)
	fmt.Printf("UdpReceiveBufferErrors %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetUdpReceiveBufferErrors(), hostname)
	fmt.Printf("UdpSendBufferErrors %d %d type=net hostname=%s\n", m.GetTimestamp(), m.GetUdpSendBufferErrors(), hostname)
}
