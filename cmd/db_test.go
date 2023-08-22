package cmd

import (
	"testing"

	"github.com/dgraph-io/badger/v4"

	"github.com/stretchr/testify/suite"
)

type CmdDbTestSuite struct {
	suite.Suite
}

func TestCommandParser(t *testing.T) {
	suite.Run(t, &CmdDbTestSuite{})
}

func (s *CmdDbTestSuite) TestPrint() {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
	s.Assert().NoError(err)
	defer db.Close()

	_ = db.Update(func(txn *badger.Txn) error {
		_ = txn.Set([]byte("bar1"), []byte(""))
		_ = txn.Set([]byte("bar2"), []byte(""))
		_ = txn.Set([]byte("bar3"), []byte(""))
		_ = txn.Set([]byte("foo1"), []byte(""))
		_ = txn.Set([]byte("foo2"), []byte(""))
		_ = txn.Set([]byte("foo3"), []byte(""))
		return nil
	})

	// DoPrint(db, "", false, nil)
	// DoPrint(db, "", true, nil)
	// DoPrint(db, "foo", false, nil)
	// DoPrint(db, "foo", true, nil)
	// DoPrint(db, "foo", true, nil)
}
