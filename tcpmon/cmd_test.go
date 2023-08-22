package tcpmon_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zperf/tcpmon/tcpmon"
)

type CommandParserTestSuite struct {
	suite.Suite
}

func TestCommandParser(t *testing.T) {
	suite.Run(t, new(CommandParserTestSuite))
}

func (s *CommandParserTestSuite) TestFileFallback() {
	f := tcpmon.FileFallback("abc", "foo", "/bin/bash")
	s.Assert().Equal("/bin/bash", f)
}
