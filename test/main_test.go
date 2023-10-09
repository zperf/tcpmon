package test

import (
	"os"
	"testing"

	"github.com/rs/zerolog"

	"github.com/zperf/tcpmon/logging"
)

func TestMain(m *testing.M) {
	logging.InitLogger(&logging.LogConfig{
		ConsoleLoggingEnabled: true,
		FileLoggingEnabled:    false,
		Level:                 zerolog.InfoLevel,
	})
	os.Exit(m.Run())
}
