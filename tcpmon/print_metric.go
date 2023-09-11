package tcpmon

type MetricPrinter interface {
	PrintNetstatMetric(*NetstatMetric, string)
	PrintNicMetric(*NicMetric, string)
	PrintTcpMetric(*TcpMetric, string)
}
