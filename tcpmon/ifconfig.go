package tcpmon

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
)

type IFCONFIGRecord struct {
	NicName    string
	Errors     int
	Dropped    int
	Overruns   int
	Frame      int
	Carrier    int
	Collisions int
}

func Ifconfig() (*[]IFCONFIGRecord, string, error) {
	cmd := cmd.NewCmd("/usr/sbin/ifconfig")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, "", errors.Wrap(ctx.Err(), "ifconfig timeout")
	case st := <-cmd.Start():
		// rr := []IFCONFIGRecord{}
		var builder strings.Builder
		ii := IFCONFIGRecord{}
		for _, line := range st.Stdout {
			builder.WriteString(line)
			if line == "" {
				ii = IFCONFIGRecord{}
			}
			if strings.Contains(line, ": flags=") {
				fields := strings.Split(line, ":")
				ii.NicName = fields[0]
				fmt.Println(fields)
				fmt.Println(len(fields))
			} else if strings.Contains(line, "RX errors ") {
				fields := strings.Split(line, " ")
				fmt.Println(fields)
				fmt.Println(len(fields))
			} else if strings.Contains(line, "TX errors ") {
				fields := strings.Split(line, " ")
				fmt.Println(fields)
				fmt.Println(len(fields))
			}
		}
	}
	return nil, "", errors.Wrap(ctx.Err(), "ifconfig timeout")
}
