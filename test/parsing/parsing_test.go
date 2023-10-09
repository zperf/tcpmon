package parsing

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zperf/tcpmon/tcpmon/tutils"
)

func TestParsing(t *testing.T) {
	suite.Run(t, new(ParsingTestSuite))
}

type ParsingTestSuite struct {
	suite.Suite
}

func (s *ParsingTestSuite) TestFileFallback() {
	f := tutils.FileFallback("abc", "foo", "/bin/bash")
	s.Assert().Equal("/bin/bash", f)
}
