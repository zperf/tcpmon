package tcpmon_test

import (
	"os"
	"testing"

	"github.com/zperf/tcpmon/tcpmon"
)

func TestMain(m *testing.M) {
	tcpmon.InitLogger()
	os.Exit(m.Run())
}
