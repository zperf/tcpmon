package main

import (
	"os"
	"runtime/debug"
	"strings"

	"github.com/zperf/tcpmon/cmd"
)

func main() {
	debug.SetMemoryLimit(100 * 1024 * 1024)
	cmd.Execute(strings.Join(os.Args[1:], " "))
}
