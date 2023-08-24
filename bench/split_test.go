package main

import (
	"strings"
	"testing"
)

var str = `ESTAB      0      0      10.255.0.95:58014              10.255.0.98:80                  users:(("tuna-rest-serve",pid=7152,fd=48))
	 skmem:(r0,rb367360,t0,tb87040,f8192,w0,o0,bl0) ts sack cubic wscale:9,9 rto:202 rtt:1.527/0.888 ato:40 mss:1448 cwnd:10 ssthresh:7 bytes_acked:258 bytes_received:2407 segs_out:5 segs_in:4 send 75.9Mbps lastsnd:9648 lastrcv:6523 lastack:6523 pacing_rate 91.0Mbps rcv_rtt:3125 rcv_space:29200
ESTAB      0      0      10.255.0.95:55226              10.255.0.96:9912                users:(("timemachine",pid=5081,fd=31)) timer:(keepalive,3.708ms,0)
	 skmem:(r0,rb1873168,t0,tb304640,f0,w0,o0,bl0) ts sack cubic wscale:9,9 rto:204 rtt:3.028/2.626 ato:40 mss:1448 cwnd:10 ssthresh:7 bytes_acked:280836 bytes_received:31869833 segs_out:14878 segs_in:27845 send 38.3Mbps lastsnd:11292 lastrcv:11292 lastack:11292 pacing_rate 45.9Mbps retrans:0/24 reordering:34 rcv_rtt:1.03 rcv_space:230366
TIME-WAIT  0      0      127.0.0.1:18500              127.0.0.1:40754               timer:(timewait,34sec,0)

ESTAB      0      0      127.0.0.1:3261               127.0.0.1:60802               users:(("zbs-chunkd",pid=5194,fd=235))
	 skmem:(r0,rb2226507,t0,tb7091712,f0,w0,o0,bl0) ts sack cubic wscale:4,9 rto:201 rtt:0.022/0.003 ato:40 mss:65483 cwnd:10 ssthresh:121 bytes_acked:1542620 bytes_received:1542780 segs_out:63978 segs_in:32138 send 238120.0Mbps lastsnd:3240 lastrcv:3240 lastack:3240 reordering:26 rcv_rtt:138382 rcv_space:43740
TIME-WAIT  0      0      192.168.17.95:80                 192.168.25.75:38716               timer:(timewait,14sec,0)`

func BenchmarkSplit(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fields := make([]string, 0)
		for _, x := range strings.Split(str, " ") {
			field := strings.TrimSpace(x)
			if len(field) > 0 {
				fields = append(fields, field)
			}
		}
		_ = fields
	}
}

func splitSpace(c rune) bool {
	return c == ' ' || c == '\r' || c == '\n' || c == '\t'
}

func BenchmarkFieldsFunc(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fields := make([]string, 0)
		for _, field := range strings.FieldsFunc(str, splitSpace) {
			fields = append(fields, field)
		}
		_ = fields
	}
}

func splitSpace2(c rune) bool {
	return c == '\n' || c == '\r' || c == '\t' || c == ' '
}

func BenchmarkFieldsFunc2(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fields := make([]string, 0)
		for _, field := range strings.FieldsFunc(str, splitSpace2) {
			fields = append(fields, field)
		}
		_ = fields
	}
}

func splitSpace3(c rune) bool {
	return c == '\t' || c == ' ' || c == '\n' || c == '\r'
}

func BenchmarkFieldsFunc3(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fields := make([]string, 0)
		for _, field := range strings.FieldsFunc(str, splitSpace3) {
			fields = append(fields, field)
		}
		_ = fields
	}
}
