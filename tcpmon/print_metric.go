package tcpmon

type PrintMetric interface {
	PrintNetstatMetric(*NetstatMetric, string)
	PrintNicMetric(*NicMetric, string)
	PrintTcpMetric(*TcpMetric, string)
}
