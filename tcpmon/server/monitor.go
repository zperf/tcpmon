package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	storagev2 "github.com/zperf/tcpmon/tcpmon/storage/v2"
)

type Monitor struct {
	config MonitorConfig

	socketMonitor *SocketMonitor
	nicMonitor    *NicMonitor
	netMonitor    *NetstatMonitor

	datastore  *storagev2.DataStore
	httpServer *http.Server
	quorum     *Quorum
}

type MonitorConfig struct {
	QuorumPort      int
	CollectInterval time.Duration
	HttpListen      string
	DataStoreConfig storagev2.Config
}

func New(config MonitorConfig) (*Monitor, error) {
	ds, err := storagev2.NewDataStore(&config.DataStoreConfig)
	if err != nil {
		return nil, err
	}

	var quorum *Quorum
	if config.QuorumPort != -1 {
		quorum = NewQuorum(&config)
	}

	return &Monitor{
		config:        config,
		datastore:     ds,
		quorum:        NewQuorum(&config),
		socketMonitor: &SocketMonitor{monitorConfig},
		nicMonitor:    &NicMonitor{monitorConfig},
		netMonitor:    &NetstatMonitor{monitorConfig},
	}, nil
}

func (m *Monitor) Collect(now time.Time, tx chan<- []byte) {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		req, err := m.socketMonitor.Collect(now)
		if err != nil {
			log.Warn().Err(err).Msg("collect socket metrics failed")
			return
		}
		tx <- req
	}()

	go func() {
		defer wg.Done()
		req, err := m.nicMonitor.Collect(now)
		if err != nil {
			log.Warn().Err(err).Msg("collect nic metrics failed")
			return
		}
		tx <- req
	}()

	go func() {
		defer wg.Done()
		req, err := m.netMonitor.Collect(now)
		if err != nil {
			log.Warn().Err(err).Msg("collect net metrics failed")
			return
		}
		tx <- req
	}()

	wg.Wait()
}

func (m *Monitor) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.config.CollectInterval)

	m.startHttpServer(m.config.HttpListen)

	if m.quorum != nil {
		members := viper.GetStringMapString("members")
		if members != nil {
			_, err := m.quorum.TryJoin(members)
			if err != nil {
				log.Warn().Err(err).Msg("Join cluster failed")
			}
		}
	}

	tx := make(chan []byte, 256)

	go func() {
		for {
			select {
			case c := <-tx:
				err := m.datastore.Put(c)
				if err != nil {
					log.Fatal().Err(err).Msg("Write failed")
				}

			case <-ctx.Done():
				log.Info().Msg("Shutting down writer...")
				return
			}
		}
	}()

	for {
		select {
		case now := <-ticker.C:
			m.Collect(now, tx)

		case <-ctx.Done():
			log.Info().Msg("Shutting down monitor...")
			m.Close()
			return nil
		}
	}
}

func (m *Monitor) Close() {
	if m.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := m.httpServer.Shutdown(ctx)
		if err != nil {
			log.Warn().Err(err).Msg("Shutdown HTTP server failed")
		}
	}

	if m.datastore != nil {
		err := m.datastore.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Close datastore failed")
		}
	}

	if m.quorum != nil {
		m.quorum.Close()
	}
}
