package tcpmon

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Monitor struct {
	sockMon    *SocketMonitor
	ifaceMon   *NicMonitor
	netstatMon *NetstatMonitor
	datastore  *Datastore
	httpServer *http.Server
	quorum     *Quorum
}

func New() (*Monitor, error) {
	epoch, err := randUint64()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate initial epoch")
	}

	path := viper.GetString("db")
	reclaimConfig := &DataStoreConfig{
		ExpectedKeyCount: viper.GetInt("max-size"),
		ReclaimBatch:     viper.GetInt("delete-size"),
		ReclaimInterval:  viper.GetDuration("reclaim-period"),
		GcInterval:       viper.GetDuration("gc-period"),
	}

	ds := NewDatastore(epoch, path, reclaimConfig)
	return &Monitor{
		sockMon:    &SocketMonitor{},
		ifaceMon:   &NicMonitor{},
		netstatMon: &NetstatMonitor{},
		datastore:  ds,
		quorum:     NewQuorum(ds),
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

func (m *Monitor) Run(ctx context.Context, interval time.Duration, addr string) error {
	ticker := time.NewTicker(interval)
	tx := m.datastore.Tx()

	m.startHttpServer(addr)

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
