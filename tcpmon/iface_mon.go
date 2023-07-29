package tcpmon

import "time"

type IfaceMonitor struct{}

func (m *IfaceMonitor) Collect(now time.Time) (*Metric, error) {
	r, _, err := ifconfig()
	if err != nil {
		return nil, err
	}
	metric := &Metric{
		Timestamp: now,
		Type:      "iface",
		Record:    r,
	}
	return metric, nil
}
