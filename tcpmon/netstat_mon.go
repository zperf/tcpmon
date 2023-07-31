package tcpmon

import "time"

type NetstatMonitor struct{}

func (m *NetstatMonitor) Collect(now time.Time) (*Metric, error) {
	r, _, err := netstat()
	if err != nil {
		return nil, err
	}
	metric := &Metric{
		Timestamp: now,
		Type:      "netstat",
		Record:    r,
	}
	return metric, nil
}
