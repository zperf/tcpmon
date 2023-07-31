package tcpmon

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
)

type IFCONFIGRecord struct {
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

func ifconfig() (*[]IFCONFIGRecord, string, error) {
	cmd := cmd.NewCmd("/sbin/ifconfig")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, "", errors.Wrap(ctx.Err(), "ifconfig timeout")
	case st := <-cmd.Start():
		rr := []IFCONFIGRecord{}
		var builder strings.Builder
		ii := IFCONFIGRecord{}
		for _, line := range st.Stdout {
			builder.WriteString(line)
			if strings.Contains(line, ": flags=") {
				fields := strings.FieldsFunc(line, func(c rune) bool {
					return c == ':'
				})
				ii.NicName = fields[0]
			} else if strings.Contains(line, "RX errors ") {
				fields := strings.FieldsFunc(line, func(c rune) bool {
					return c == ' '
				})
				ii.RXErrors, _ = strconv.Atoi(fields[2])
				ii.RXDropped, _ = strconv.Atoi(fields[4])
				ii.RXOverruns, _ = strconv.Atoi(fields[6])
				ii.RXFrame, _ = strconv.Atoi(fields[8])
			} else if strings.Contains(line, "TX errors ") {
				fields := strings.FieldsFunc(line, func(c rune) bool {
					return c == ' '
				})
				ii.TXErrors, _ = strconv.Atoi(fields[2])
				ii.TXDropped, _ = strconv.Atoi(fields[4])
				ii.TXOverruns, _ = strconv.Atoi(fields[6])
				ii.TXCarrier, _ = strconv.Atoi(fields[8])
				ii.TXCollisions, _ = strconv.Atoi(fields[10])
				rr = append(rr, ii)
				ii = IFCONFIGRecord{}
			}
		}
		return &rr, builder.String(), nil
	}
}
