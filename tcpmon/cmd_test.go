package tcpmon_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CommandParserTestSuite struct {
	suite.Suite
}

func TestCommandParser(t *testing.T) {
	suite.Run(t, new(CommandParserTestSuite))
}
