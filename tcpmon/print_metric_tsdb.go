package tcpmon

import "fmt"

func boolToUint32(x bool) uint32 {
	if !x {
		return 0
	} else {
		return 1
	}
}

type TSDBPrintMetric struct {}

func (tsdb TSDBPrintMetric) PrintTcpMetric(m *TcpMetric, hostname string) {
	for _, socket := range m.GetSockets() {
		fmt.Printf("State type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetState())
		fmt.Printf("RecvQ type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetRecvQ())
		fmt.Printf("SendQ type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSendQ())

		for _, process := range socket.GetProcesses() {
			fmt.Printf("Pid type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s,ProcessName=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), process.GetName(), m.GetTimestamp(), process.GetPid())
			fmt.Printf("Fd type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s,ProcessName=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), process.GetName(), m.GetTimestamp(), process.GetFd())
		}

		for _, timer := range socket.GetTimers() {
			fmt.Printf("ExpireTimeUs type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s,TimerName=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), timer.GetName(), m.GetTimestamp(), timer.GetExpireTimeUs())
			fmt.Printf("Retrans type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s,TimerName=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), timer.GetName(), m.GetTimestamp(), timer.GetRetrans())
		}

		fmt.Printf("RmemAlloc type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetRmemAlloc())
		fmt.Printf("RcvBuf type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetRcvBuf())
		fmt.Printf("WmemAlloc type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetWmemAlloc())
		fmt.Printf("SndBuf type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetSndBuf())
		fmt.Printf("FwdAlloc type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetFwdAlloc())
		fmt.Printf("WmemQueued type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetWmemQueued())
		fmt.Printf("OptMem type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetOptMem())
		fmt.Printf("BackLog type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetBackLog())
		fmt.Printf("SockDrop type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetSockDrop())

		fmt.Printf("Ts type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), boolToUint32(socket.GetTs()))
		fmt.Printf("Sack type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), boolToUint32(socket.GetSack()))
		fmt.Printf("Cubic type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), boolToUint32(socket.GetCubic()))
		fmt.Printf("AppLimited type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), boolToUint32(socket.GetAppLimited()))
		fmt.Printf("PacingRate type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %f\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetPacingRate())
		fmt.Printf("DeliveryRate type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %f\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetDeliveryRate())
		fmt.Printf("Send type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %f\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSend())
		fmt.Printf("SndWscale type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSndWscale())
		fmt.Printf("RcvWscale type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetRcvWscale())
		fmt.Printf("Rto type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %f\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetRto())
		fmt.Printf("Rtt type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %f\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetRtt())
		fmt.Printf("Rttvar type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %f\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetRttvar())
		fmt.Printf("Minrtt type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %f\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetMinrtt())
		fmt.Printf("RcvRtt type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %f\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetRcvRtt())
		fmt.Printf("RetransNow type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetRetransNow())
		fmt.Printf("RetransTotal type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetRetransTotal())
		fmt.Printf("Ato type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %f\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetAto())
		fmt.Printf("Mss type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetMss())
		fmt.Printf("Pmtu type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetPmtu())
		fmt.Printf("Rcvmss type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetRcvmss())
		fmt.Printf("Advmss type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetAdvmss())
		fmt.Printf("Cwnd type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetCwnd())
		fmt.Printf("SndWnd type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSndWnd())
		fmt.Printf("BytesSent type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetBytesSent())
		fmt.Printf("BytesAcked type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetBytesAcked())
		fmt.Printf("BytesReceived type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetBytesReceived())
		fmt.Printf("SegsOut type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSegsOut())
		fmt.Printf("SegsIn type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSegsIn())
		fmt.Printf("Lastsnd type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetLastsnd())
		fmt.Printf("Lastrcv type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetLastrcv())
		fmt.Printf("Lastack type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetLastack())
		fmt.Printf("Delivered type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetDelivered())
		fmt.Printf("BusyMs type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetBusyMs())
		fmt.Printf("RcvSpace type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetRcvSpace())
		fmt.Printf("RcvSsthresh type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetRcvSsthresh())
		fmt.Printf("DataSegsOut type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetDataSegsOut())
		fmt.Printf("DataSegsIn type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetDataSegsIn())
		fmt.Printf("RwndLimited type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetRwndLimited())
		fmt.Printf("SndbufLimited type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSndbufLimited())
		fmt.Printf("Ecn type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), boolToUint32(socket.GetEcn()))
		fmt.Printf("Ecnseen type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\n", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), boolToUint32(socket.GetEcnseen()))
	}
}

func (tsdb TSDBPrintMetric) PrintNicMetric(m *NicMetric, hostname string) {
	for _, iface := range m.GetIfaces() {
		fmt.Printf("RxErrors type=nic,hostname=%s,name=%s %d %d\n", hostname, iface.GetName(), m.GetTimestamp(), iface.GetRxErrors())
		fmt.Printf("RxDropped type=nic,hostname=%s,name=%s %d %d\n", hostname, iface.GetName(), m.GetTimestamp(), iface.GetRxDropped())
		fmt.Printf("RxOverruns type=nic,hostname=%s,name=%s %d %d\n", hostname, iface.GetName(), m.GetTimestamp(), iface.GetRxOverruns())
		fmt.Printf("RxFrame type=nic,hostname=%s,name=%s %d %d\n", hostname, iface.GetName(), m.GetTimestamp(), iface.GetRxFrame())
		fmt.Printf("TxErrors type=nic,hostname=%s,name=%s %d %d\n", hostname, iface.GetName(), m.GetTimestamp(), iface.GetTxErrors())
		fmt.Printf("TxDropped type=nic,hostname=%s,name=%s %d %d\n", hostname, iface.GetName(), m.GetTimestamp(), iface.GetTxDropped())
		fmt.Printf("TxOverruns type=nic,hostname=%s,name=%s %d %d\n", hostname, iface.GetName(), m.GetTimestamp(), iface.GetTxOverruns())
		fmt.Printf("TxCarrier type=nic,hostname=%s,name=%s %d %d\n", hostname, iface.GetName(), m.GetTimestamp(), iface.GetTxCarrier())
		fmt.Printf("TxCollisions type=nic,hostname=%s,name=%s %d %d\n", hostname, iface.GetName(), m.GetTimestamp(), iface.GetTxCollisions())
	}
}

func (tsdb TSDBPrintMetric) PrintNetstatMetric(m *NetstatMetric, hostname string) {
	fmt.Printf("IpTotalPacketsReceived type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetIpTotalPacketsReceived())
	fmt.Printf("IpForwarded type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetIpForwarded())
	fmt.Printf("IpIncomingPacketsDiscarded type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetIpIncomingPacketsDiscarded())
	fmt.Printf("IpIncomingPacketsDelivered type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetIpIncomingPacketsDelivered())
	fmt.Printf("IpRequestsSentOut type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetIpRequestsSentOut())
	fmt.Printf("IpOutgoingPacketsDropped type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetIpOutgoingPacketsDropped())
	fmt.Printf("TcpActiveConnectionsOpenings type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetTcpActiveConnectionsOpenings())
	fmt.Printf("TcpPassiveConnectionOpenings type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetTcpPassiveConnectionOpenings())
	fmt.Printf("TcpFailedConnectionAttempts type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetTcpFailedConnectionAttempts())
	fmt.Printf("TcpConnectionResetsReceived type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetTcpConnectionResetsReceived())
	fmt.Printf("TcpConnectionsEstablished type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetTcpConnectionsEstablished())
	fmt.Printf("TcpSegmentsReceived type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetTcpSegmentsReceived())
	fmt.Printf("TcpSegmentsSendOut type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetTcpSegmentsSendOut())
	fmt.Printf("TcpSegmentsRetransmitted type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetTcpSegmentsRetransmitted())
	fmt.Printf("TcpBadSegmentsReceived type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetTcpBadSegmentsReceived())
	fmt.Printf("TcpResetsSent type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetTcpResetsSent())
	fmt.Printf("UdpPacketsReceived type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetUdpPacketsReceived())
	fmt.Printf("UdpPacketsToUnknownPortReceived type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetUdpPacketsToUnknownPortReceived())
	fmt.Printf("UdpPacketReceiveErrors type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetUdpPacketReceiveErrors())
	fmt.Printf("UdpPacketsSent type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetUdpPacketsSent())
	fmt.Printf("UdpReceiveBufferErrors type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetUdpReceiveBufferErrors())
	fmt.Printf("UdpSendBufferErrors type=net,hostname=%s %d %d\n", hostname, m.GetTimestamp(), m.GetUdpSendBufferErrors())
}
