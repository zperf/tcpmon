package tcpmon

import (
	"context"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
)

type Monitor struct {
	sockMon   *SocketMonitor
	ifaceMon  *IfaceMonitor
	datastore *Datastore
}

func New() (*Monitor, error) {
	epoch, err := randUint64()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate initial epoch")
	}

	return &Monitor{
		sockMon:   &SocketMonitor{},
		ifaceMon:  &IfaceMonitor{},
		datastore: NewDatastore(epoch),
	}, nil
}

func (mon *Monitor) Collect(now time.Time, tx chan<- *StoreRequest) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		req, err := mon.sockMon.Collect(now)
		if err != nil {
			log.Warn().Err(err).Msg("collect socket metrics failed")
		}
		tx <- req
		wg.Done()
	}()

	//go func() {
	//	m, err := mon.ifaceMon.Collect(now)
	//	if err != nil {
	//		log.Warn().Err(err).Msg("collect iface metrics failed")
	//	}
	//	tx <- m
	//	wg.Done()
	//}()

	wg.Wait()
}

func (mon *Monitor) Run(ctx context.Context) error {
	ticker := time.NewTicker(1000 * time.Millisecond)
	tx := mon.datastore.Tx()
	for {
		select {
		case now := <-ticker.C:
			mon.Collect(now, tx)
		case <-ctx.Done():
			mon.Close()
			return nil
		}
	}
}

func (mon *Monitor) Close() {
	if mon.datastore != nil {
		mon.datastore.Close()
	}
}
