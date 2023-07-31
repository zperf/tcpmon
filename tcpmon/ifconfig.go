package tcpmon

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
)

type IfaceRecord struct {
	NicName string

	RXErrors   int
	RXDropped  int
	RXOverruns int
	RXFrame    int

	TXErrors     int
	TXDropped    int
	TXOverruns   int
	TXCarrier    int
	TXCollisions int
}

func ifconfig() (*[]IfaceRecord, string, error) {
	cmd := cmd.NewCmd("/sbin/ifconfig")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, "", errors.Wrap(ctx.Err(), "ifconfig timeout")
	case st := <-cmd.Start():
		var records []IfaceRecord
		var builder strings.Builder

		r := IfaceRecord{}
		for _, line := range st.Stdout {
			builder.WriteString(line)
			if strings.Contains(line, ": flags=") {
				fields := strings.FieldsFunc(line, func(c rune) bool {
					return c == ':'
				})
				r.NicName = fields[0]
			} else if strings.Contains(line, "RX errors ") {
				fields := strings.FieldsFunc(line, func(c rune) bool {
					return c == ' '
				})
				r.RXErrors, _ = strconv.Atoi(fields[2])
				r.RXDropped, _ = strconv.Atoi(fields[4])
				r.RXOverruns, _ = strconv.Atoi(fields[6])
				r.RXFrame, _ = strconv.Atoi(fields[8])
			} else if strings.Contains(line, "TX errors ") {
				fields := strings.FieldsFunc(line, func(c rune) bool {
					return c == ' '
				})
				r.TXErrors, _ = strconv.Atoi(fields[2])
				r.TXDropped, _ = strconv.Atoi(fields[4])
				r.TXOverruns, _ = strconv.Atoi(fields[6])
				r.TXCarrier, _ = strconv.Atoi(fields[8])
				r.TXCollisions, _ = strconv.Atoi(fields[10])
				records = append(records, r)
				r = IfaceRecord{}
			}
		}
		return &records, builder.String(), nil
	}
}
