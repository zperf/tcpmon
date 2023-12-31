package influxdb

import (
	"bufio"
	"io"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/afero"
	proto2 "google.golang.org/protobuf/proto"

	"github.com/zperf/tcpmon/tcpmon/gproto"
	"github.com/zperf/tcpmon/tcpmon/storage"
)

var ErrTimePointNotIncluded = errors.New("time point not included in this data file")

type FastExporter struct {
	fh     afero.File
	ranges []RecordRange
}

type Range struct {
	Offset int64
	Len    uint32
}

type RecordRange struct {
	Header Range
	Body   Range
}

// NewFastExporter creates a new exporter
// f: file path
// fs: A mock for unit test. Pass nil for a real fs.
func NewFastExporter(f string, fs afero.Fs) (*FastExporter, error) {
	if fs == nil {
		fs = afero.NewOsFs()
	}

	fh, err := fs.Open(f)
	if err != nil {
		return nil, errors.Wrap(err, "open file failed")
	}

	return &FastExporter{
		fh:     fh,
		ranges: nil,
	}, nil
}

func (r *FastExporter) Close() {
	if r.fh != nil {
		err := r.fh.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Close file failed")
		}
	}
}

func (r *FastExporter) ReadAt(offset int64, len uint32) ([]byte, error) {
	buf := make([]byte, len)
	_, err := r.fh.ReadAt(buf, offset)
	if err != nil {
		return nil, errors.Wrap(err, "read at failed")
	}

	return buf, nil
}

func (r *FastExporter) ReadRange(ra Range) ([]byte, error) {
	return r.ReadAt(ra.Offset, ra.Len)
}

func (r *FastExporter) UnmarshalMetric(buf []byte) (*gproto.Metric, error) {
	var m gproto.Metric
	err := proto2.Unmarshal(buf, &m)
	if err != nil {
		return nil, errors.Wrap(err, "parse failed")
	}
	return &m, nil
}

func getTimestamp(m *gproto.Metric) (time.Time, error) {
	switch m := m.Body.(type) {
	case *gproto.Metric_Tcp:
		return time.Unix(m.Tcp.GetTimestamp(), 0), nil
	case *gproto.Metric_Nic:
		return time.Unix(m.Nic.GetTimestamp(), 0), nil
	case *gproto.Metric_Net:
		return time.Unix(m.Net.GetTimestamp(), 0), nil
	default:
		return time.Time{}, errors.New("unknown metric type")
	}
}

type ExportOptions struct {
	Hostname string
	Target   time.Time
	ShowOnly bool

	WriteDb   bool
	Org       string
	Bucket    string
	Token     string
	DbAddress string

	Bar *progressbar.ProgressBar
}

func (r *FastExporter) Export(w io.Writer, option *ExportOptions) error {
	ra, err := r.Scan()
	if err != nil {
		return err
	}

	getRecordTime := func(rr RecordRange) (time.Time, error) {
		buf, err := r.ReadRange(rr.Body)
		if err != nil {
			return time.Time{}, err
		}

		metric, err := r.UnmarshalMetric(buf)
		if err != nil {
			return time.Time{}, err
		}

		start, err := getTimestamp(metric)
		if err != nil {
			return time.Time{}, err
		}
		return start, nil
	}

	var start, end time.Time
	start, err = getRecordTime(ra[0])
	if err != nil {
		return err
	}

	end, err = getRecordTime(ra[len(ra)-1])
	if err != nil {
		return err
	}

	if start.After(end) {
		start, end = end, start
	}

	if option.ShowOnly {
		log.Info().Time("start", start).Time("end", end).Str("file", r.fh.Name()).Send()
		return nil
	}

	needExport := true
	if !option.Target.IsZero() {
		if !start.Before(option.Target) || !option.Target.Before(end) {
			needExport = false
		}
	}

	if !needExport {
		if option.Bar == nil {
			log.Info().Time("start", start).Time("end", end).Msg("ignored")
		}
		return ErrTimePointNotIncluded
	}

	var exit sync.WaitGroup
	var writerMutex sync.Mutex
	jobs := make(chan RecordRange, 256)
	defer close(jobs)
	stop := make(chan struct{})

	workerNum := runtime.NumCPU()
	exit.Add(workerNum)
	for i := 0; i < workerNum; i++ {
		go r.exportWorker(w, &writerMutex, jobs, &exit, stop, *option)
	}

	for _, rr := range ra {
		if option.Bar != nil {
			err = option.Bar.Add(1)
			if err != nil {
				log.Warn().Err(err).Msg("Add progress bar failed")
			}
		}
		jobs <- rr
	}

	close(stop)
	exit.Wait()
	return nil
}

func (r *FastExporter) exportWorker(w io.Writer, m *sync.Mutex, jobs <-chan RecordRange, exit *sync.WaitGroup,
	stop <-chan struct{}, option ExportOptions) {
work:
	for {
		select {
		case <-stop:
			break work
		default:
		}

		select {
		case job := <-jobs:
			buf, err := r.ReadRange(job.Body)
			if err != nil {
				log.Fatal().Err(err).Msg("Read failed")
			}

			metric, err := r.UnmarshalMetric(buf)
			if err != nil {
				log.Fatal().Err(err).Msg("Unmarshal failed")
			}

			if option.WriteDb {
				// use influxdb client write metrics to db
				conn := NewImporter(&ImportOption{
					Bucket:   option.Bucket,
					Org:      option.Org,
					Token:    option.Token,
					Address:  option.DbAddress,
					Hostname: option.Hostname,
				})

				err = conn.Submit(metric)
				if err != nil {
					conn.Close()
					log.Fatal().Err(err).Msg("Submit failed")
				}
				conn.Close()

			} else {
				// export metrics to txt file with line protocol
				var builder strings.Builder
				exporter := New(option.Hostname, &builder)
				exporter.ExportMetric(metric)

				m.Lock()
				writer := bufio.NewWriter(w)
				_, err = writer.WriteString(builder.String())
				if err != nil {
					log.Fatal().Err(err).Msg("Write failed")
				}

				err = writer.Flush()
				if err != nil {
					log.Fatal().Err(err).Msg("Flush failed")
				}
				m.Unlock()
			}

		case <-time.After(time.Second):
			log.Info().Msg("Get job timeout")
		}
	}
	exit.Done()
}

func (r *FastExporter) Scan() ([]RecordRange, error) {
	if r.ranges == nil {
		return r.doScan()
	}

	return r.ranges, nil
}

func (r *FastExporter) doScan() ([]RecordRange, error) {
	ranges := make([]RecordRange, 0)

	offset := int64(0)
	var err error
	for {
		var ra RecordRange
		var size uint32

		size, err = storage.ReadHeader(r.fh)
		if err != nil {
			break
		}
		ra.Header.Offset = offset
		ra.Header.Len = storage.HeaderSize
		ra.Body.Offset = offset + storage.HeaderSize
		ra.Body.Len = size

		offset, err = r.fh.Seek(int64(size), io.SeekCurrent)
		if err != nil {
			break
		}

		ranges = append(ranges, ra)
	}
	if err != nil && err != io.EOF {
		return nil, errors.Wrapf(err, "scan data file failed")
	}

	r.ranges = ranges
	return ranges, nil
}
