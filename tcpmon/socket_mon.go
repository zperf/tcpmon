package tcpmon

import (
	"time"
)

type SocketMonitor struct{}

func (m *SocketMonitor) Collect(now time.Time) (*Metric, error) {
	r, _, err := ss()
	if err != nil {
		return nil, err
	}
	metric := &Metric{
		Timestamp: now,
		Type:      "socket",
		Record:    r,
		// Raw:       raw,
	}
	return metric, nil
}
