package tcpmon

import (
	"time"

	"github.com/rs/zerolog/log"
)

type Monitor struct {
	sockMon    *SocketMonitor
	ifaceMon   *IfaceMonitor
	netstatMon *NetstatMonitor
}

func New() *Monitor {
	return &Monitor{
		sockMon:    &SocketMonitor{},
		ifaceMon:   &IfaceMonitor{},
		netstatMon: &NetstatMonitor{},
	}
}

func (m *Monitor) Collect(now time.Time) {
	// metric, err := m.sockMon.Collect(now)
	// metric, err := m.ifaceMon.Collect(now)
	metric, err := m.netstatMon.Collect(now)
	if err != nil {
		log.Warn().Err(err).Msg("collect socket metrics failed")
	}

	log.Info().Str("type", metric.Type).Msgf("metric: %+v", metric.Record)
}

func (m *Monitor) Run() error {
	ticker := time.NewTicker(1000 * time.Millisecond)
	for {
		m.Collect(<-ticker.C)
	}
}
