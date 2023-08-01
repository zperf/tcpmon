package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/zperf/tcpmon/cmd"
)

func main() {
	// TODO(fanyang) log should write to file
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	cmd.Execute()
}
