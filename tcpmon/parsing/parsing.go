package parsing

import (
	"bufio"
	"io"

	"github.com/cockroachdb/errors"

	"github.com/zperf/tcpmon/tcpmon/gproto"
)

func parseSnmpOrNetstat(r io.Reader, m *gproto.NetstatMetric, parseFn func(string, string, *gproto.NetstatMetric) error) error {
	lines := make([]string, 2)
	p := 0

	s := bufio.NewScanner(r)
	for s.Scan() {
		lines[p] = s.Text()
		p++
		if p >= len(lines) {
			p = 0

			err := parseFn(lines[0], lines[1], m)
			if err != nil {
				return err
			}
		}
	}

	err := s.Err()
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
