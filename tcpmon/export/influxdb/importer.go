package influxdb

import (
	"time"

	"github.com/cockroachdb/errors"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/zperf/tcpmon/tcpmon/tproto"
)

type Importer struct {
	option ImportOption
	client influxdb2.Client
}

type ImportOption struct {
	Bucket   string
	Org      string
	Token    string
	Address  string
	Hostname string
}

func NewImporter(option *ImportOption) *Importer {
	client := influxdb2.NewClient(option.Address, option.Token)

	return &Importer{
		option: *option,
		client: client,
	}
}

func (im *Importer) Close() {
	if im.client != nil {
		im.client.Close()
	}
}

func (im *Importer) Submit(metric *tproto.Metric) error {
	writeAPI := im.client.WriteAPI(im.option.Org, im.option.Bucket)
	errCh := writeAPI.Errors()
	conv := NewMetricConv(im.option.Hostname)

	rawTs, points := conv.Metric(metric)
	ts := time.Unix(rawTs, 0)

	for _, p := range points {
		writeAPI.WritePoint(p)

		select {
		case err := <-errCh:
			if err != nil {
				return errors.Wrapf(err, "write point failed, ts=%v", ts)
			}
		default:
		}
	}

	writeAPI.Flush()
	return nil
}
