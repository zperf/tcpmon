package tcpmon

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-cmd/cmd"
)

type SSRecord struct {
	State     string `bson:"state"`
	RecvQ     int    `bson:"recvq"`
	SendQ     int    `bson:"sendq"`
	LocalAddr string `bson:"local_addr"`
	PeerAddr  string `bson:"peer_addr"`
}

func ss() (*[]SSRecord, string, error) {
	cmd := cmd.NewCmd("/usr/bin/ss", "-4ntipmoHOna")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, "", errors.Wrap(ctx.Err(), "ss timeout")
	case st := <-cmd.Start():
		r := []SSRecord{}
		var builder strings.Builder
		for _, line := range st.Stdout {
			builder.WriteString(line)
			s := SSRecord{}
			fields := strings.FieldsFunc(line, func(c rune) bool {
				return c == ' '
			})
			s.State = fields[0]
			n, _ := strconv.Atoi(fields[1])
			s.RecvQ = n
			n, _ = strconv.Atoi(fields[2])
			s.SendQ = n
			s.LocalAddr = fields[3]
			s.PeerAddr = fields[4]
			r = append(r, s)
		}
		return &r, builder.String(), nil
	}
}
