package tcpmon

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
)

type Monitor struct {
	config     MonitorConfig
	sockMon    *SocketMonitor
	ifaceMon   *NicMonitor
	netstatMon *NetstatMonitor
	datastore  *DataStore
	httpServer *http.Server
	quorum     *Quorum
}

type MonitorConfig struct {
	DataStoreConfig

	QuorumPort      int
	CollectInterval time.Duration
	HttpListen      string
}

func New(config MonitorConfig) (*Monitor, error) {
	epoch, err := randUint64()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate initial epoch")
	}

	ds := NewDataStore(epoch, &config.DataStoreConfig)
	cmdConfig := NewCmdConfig()

	return &Monitor{
		config:     config,
		sockMon:    &SocketMonitor{cmdConfig},
		ifaceMon:   &NicMonitor{cmdConfig},
		netstatMon: &NetstatMonitor{cmdConfig},
		datastore:  ds,
		quorum:     NewQuorum(ds, &config),
	}, nil
}

func (m *Monitor) Collect(now time.Time, tx chan<- *KVPair) {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		req, err := m.sockMon.Collect(now)
		if err != nil {
			log.Warn().Err(err).Msg("collect socket metrics failed")
			return
		}
		tx <- req
	}()

	go func() {
		defer wg.Done()
		req, err := m.ifaceMon.Collect(now)
		if err != nil {
			log.Warn().Err(err).Msg("collect nic metrics failed")
			return
		}
		tx <- req
	}()

	go func() {
		defer wg.Done()
		req, err := m.netstatMon.Collect(now)
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
	tx := m.datastore.Tx()

	m.startHttpServer(m.config.HttpListen)

	initialMembers, err := m.datastore.GetMemberAddressList()
	if err != nil {
		log.Info().Err(err).Msg("Get members from db failed")
	} else if len(initialMembers) != 0 {
		m.quorum.Join(initialMembers)
	}

	for {
		select {
		case now := <-ticker.C:
			m.Collect(now, tx)
		case <-ctx.Done():
			log.Info().Err(ctx.Err()).Msg("shutting down...")
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
			log.Warn().Err(err).Msg("failed shutdown HTTP server")
		}
	}
	if m.datastore != nil {
		m.datastore.Close()
	}
	if m.quorum != nil {
		m.quorum.Close()
	}
}
