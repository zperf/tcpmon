package main

import (
	"github.com/zperf/tcpmon/cmd"
	"github.com/zperf/tcpmon/tcpmon"
)

func main() {
	tcpmon.InitLogger()
	cmd.Execute()
}
