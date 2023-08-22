package main

import (
	"os"
	"strings"

	"github.com/zperf/tcpmon/cmd"
)

func main() {
	cmd.Execute(strings.Join(os.Args[1:], " ") != "config default")
}
