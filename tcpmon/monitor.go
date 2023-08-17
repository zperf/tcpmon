package tcpmon

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Monitor struct {
	sockMon      *SocketMonitor
	ifaceMon     *NicMonitor
	netstatMon   *NetstatMonitor
	datastore    *Datastore
	httpServer   *http.Server
	gossipServer *GossipServer
}

func New() (*Monitor, error) {
	epoch, err := randUint64()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate initial epoch")
	}

	path := viper.GetString("db")
	periodOptions := &PeriodOption{
		MaxSize:       viper.GetInt("max-size"),
		DeleteSize:    viper.GetInt("delete-size"),
		ReclaimPeriod: viper.GetDuration("reclaim-period"),
		GCPeriod:      viper.GetDuration("gc-period"),
	}

	return &Monitor{
		sockMon:      &SocketMonitor{},
		ifaceMon:     &NicMonitor{},
		netstatMon:   &NetstatMonitor{},
		datastore:    NewDatastore(epoch, path, periodOptions),
		gossipServer: NewGossipServer(),
	}, nil
}

func (mon *Monitor) Collect(now time.Time, tx chan<- *KVPair) {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		req, err := mon.sockMon.Collect(now)
		if err != nil {
			log.Warn().Err(err).Msg("collect socket metrics failed")
			return
		}
		tx <- req
	}()

	go func() {
		defer wg.Done()
		req, err := mon.ifaceMon.Collect(now)
		if err != nil {
			log.Warn().Err(err).Msg("collect nic metrics failed")
			return
		}
		tx <- req
	}()

	go func() {
		defer wg.Done()
		req, err := mon.netstatMon.Collect(now)
		if err != nil {
			log.Warn().Err(err).Msg("collect net metrics failed")
			return
		}
		tx <- req
	}()

	wg.Wait()
}

func (mon *Monitor) Run(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	tx := mon.datastore.Tx()

	mon.startHttpServer(viper.GetString("listen"))

	initialMembers := strings.FieldsFunc(viper.GetString("initial-members"), func(c rune) bool {
		return unicode.IsSpace(c) || c == ','
	})
	mon.gossipServer.Join(initialMembers)

	for {
		select {
		case now := <-ticker.C:
			mon.Collect(now, tx)
		case <-ctx.Done():
			log.Info().Err(ctx.Err()).Msg("shutting down...")
			mon.Close()
			return nil
		}
	}
}

func (mon *Monitor) Close() {
	if mon.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := mon.httpServer.Shutdown(ctx)
		if err != nil {
			log.Warn().Err(err).Msg("failed shutdown HTTP server")
		}
	}
	if mon.datastore != nil {
		mon.datastore.Close()
	}
}
