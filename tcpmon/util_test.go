package tcpmon

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilTestSuite struct {
	suite.Suite
}

func TestUtil(t *testing.T) {
	suite.Run(t, new(UtilTestSuite))
}

func (s *UtilTestSuite) TestParseUint32() {
	assert := s.Assert()
	val, err := ParseUint32("123")
	assert.NoError(err)
	assert.Equal(uint32(123), val)

	val, err = ParseUint32("")
	assert.Error(err)
	assert.Equal(uint32(0), val)
}

func (s *UtilTestSuite) TestParseInt() {
	assert := s.Assert()
	val, err := ParseInt("123")
	assert.NoError(err)
	assert.Equal(123, val)

	val, err = ParseInt("")
	assert.Error(err)
	assert.Equal(0, val)
}
