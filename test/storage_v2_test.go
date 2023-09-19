package test

import (
	"crypto/rand"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"

	v2 "github.com/zperf/tcpmon/tcpmon/storage/v2"
)

type StorageV2TestSuite struct {
	suite.Suite
	fs      afero.Fs
	baseDir string
}

func TestStorageV2(t *testing.T) {
	s := &StorageV2TestSuite{
		fs:      afero.NewBasePathFs(afero.NewOsFs(), "./tmp"),
		baseDir: "db",
	}

	suite.Run(t, s)
}

// SetupTest run before each test in the suite
func (suite *StorageV2TestSuite) SetupTest() {
	err := suite.fs.RemoveAll(suite.baseDir)
	if err != nil {
		log.Fatal().Err(err).Msg("Delete dir failed")
	}

	err = suite.fs.MkdirAll(suite.baseDir, 0755)
	if err != nil {
		log.Fatal().Err(err).Msg("Create dir failed")
	}
}

// TestBasic perform basic functional tests
func (suite *StorageV2TestSuite) TestBasic() {
	ds, err := v2.NewDataStore(v2.NewConfig(suite.baseDir).WithFs(suite.fs))
	suite.Require().NoError(err)
	defer ds.Close()

	err = ds.Put(randBuf(1 << 10))
	suite.Require().NoError(err)

	err = ds.Put(randBuf(1 << 10))
	suite.Require().NoError(err)

	err = ds.Put(randBuf(1 << 10))
	suite.Require().NoError(err)
}

func (suite *StorageV2TestSuite) TestRotateFile() {
	cfg := v2.NewConfig(suite.baseDir).
		WithFs(suite.fs).
		WithMaxSize(10 * (1 << 20)).
		WithMaxEntriesPerFile(3)

	ds, err := v2.NewDataStore(cfg)
	suite.Require().NoError(err)

	const toWrite = 10
	const bufSize = 1 << 10
	buf := randBuf(bufSize)
	for i := 0; i < toWrite; i++ {
		err := ds.Put(buf)
		suite.Require().NoError(err)
	}
	ds.Close()

	r, err := v2.NewDataStoreReader(v2.NewReaderConfig(suite.baseDir).WithFs(suite.fs))
	suite.Require().NoError(err)

	count := 0
	err = r.Iterate(func(buf []byte) {
		suite.Require().Equal(bufSize, len(buf))
		count++
	})
	suite.Require().NoError(err)
	suite.Require().Equal(toWrite, count)
}

func (suite *StorageV2TestSuite) TestGetLastFileNo() {
	_, err := suite.fs.Create("tcpmon-dataf-1")
	suite.Require().NoError(err)
	_, err = suite.fs.Create("tcpmon-dataf-1.zst")
	suite.Require().NoError(err)

	cfg := v2.NewConfig(suite.baseDir).
		WithFs(suite.fs).
		WithMaxSize(10 * (1 << 20)).
		WithMaxEntriesPerFile(3)

	ds, err := v2.NewDataStore(cfg)
	suite.Require().NoError(err)
	defer ds.Close()

	last := ds.GetLatestFileNo()
	suite.Require().Equal(uint32(1), last)
}

func randBuf(size int) []byte {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate new random buffer")
	}

	return buf
}
