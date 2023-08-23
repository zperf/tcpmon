package main

import (
	"os"
	"runtime/debug"
	"strings"

	"github.com/zperf/tcpmon/cmd"
)

func main() {
	debug.SetMemoryLimit(100 << 20)
	cmd.Execute(strings.Join(os.Args[1:], " "))
}
