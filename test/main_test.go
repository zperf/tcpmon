package test

import (
	"os"
	"testing"

	"github.com/rs/zerolog"

	"github.com/zperf/tcpmon/tcpmon"
)

func TestMain(m *testing.M) {
	tcpmon.InitLogger(&tcpmon.LogConfig{
		ConsoleLoggingEnabled: true,
		FileLoggingEnabled:    false,
		Level:                 zerolog.InfoLevel,
	})
	os.Exit(m.Run())
}
