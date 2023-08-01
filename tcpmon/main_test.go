package tcpmon_test

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestMain(m *testing.M) {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	os.Exit(m.Run())
}
