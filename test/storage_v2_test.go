package test

import (
	"crypto/rand"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"

	"github.com/zperf/tcpmon/tcpmon/storage"
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
func (s *StorageV2TestSuite) SetupTest() {
	err := s.fs.RemoveAll(s.baseDir)
	if err != nil {
		log.Fatal().Err(err).Msg("Delete dir failed")
	}

	err = s.fs.MkdirAll(s.baseDir, 0755)
	if err != nil {
		log.Fatal().Err(err).Msg("Create dir failed")
	}
}

// TestBasic perform basic functional tests
func (s *StorageV2TestSuite) TestBasic() {
	ds, err := storage.NewDataStore(storage.NewConfig(s.baseDir).WithFs(s.fs))
	s.Require().NoError(err)
	defer ds.Close()

	err = ds.Put(randBuf(1 << 10))
	s.Require().NoError(err)

	err = ds.Put(randBuf(1 << 10))
	s.Require().NoError(err)

	err = ds.Put(randBuf(1 << 10))
	s.Require().NoError(err)
}

func (s *StorageV2TestSuite) TestRotateFile() {
	cfg := storage.NewConfig(s.baseDir).
		WithFs(s.fs).
		WithMaxSize(10 * (1 << 20)).
		WithMaxEntriesPerFile(3)

	ds, err := storage.NewDataStore(cfg)
	s.Require().NoError(err)

	const toWrite = 10
	const bufSize = 1 << 10
	buf := randBuf(bufSize)
	for i := 0; i < toWrite; i++ {
		err := ds.Put(buf)
		s.Require().NoError(err)
	}
	err = ds.Close()
	s.Require().NoError(err)

	r, err := storage.NewDataStoreReader(storage.NewReaderConfig(s.baseDir).WithFs(s.fs))
	s.Require().NoError(err)

	count := 0
	err = r.Iterate(func(buf []byte) error {
		s.Require().Equal(bufSize, len(buf))
		count++
		return nil
	})
	s.Require().NoError(err)
	s.Require().Equal(toWrite, count)
}

func (s *StorageV2TestSuite) TestGetLastestFileNo() {
	_, err := s.fs.Create(s.baseDir + "/tcpmon-dataf-1")
	s.Require().NoError(err)
	_, err = s.fs.Create(s.baseDir + "/tcpmon-dataf-2.zst")
	s.Require().NoError(err)

	cfg := storage.NewConfig(s.baseDir).
		WithFs(s.fs).
		WithMaxSize(10 * (1 << 20)).
		WithMaxEntriesPerFile(3)

	ds, err := storage.NewDataStore(cfg)
	s.Require().NoError(err)
	defer ds.Close()

	last := ds.GetLatestFileNo()
	s.Require().Equal(uint32(3), last)
}

func (s *StorageV2TestSuite) TestGetLatestFileNo2() {
	fileNames := []string{
		"tcpmon-dataf-1",
		"tcpmon-dataf-2",
		"tcpmon-dataf-3.zst",
		"tcpmon-dataf-4.zst",
		"tcpmon-dataf-5.zst",
		"tcpmon-dataf-6.zst",
		"tcpmon-dataf-7.zst",
		"tcpmon-dataf-8",
		"tcpmon-dataf-9.zst",
		"tcpmon-dataf-10.zst",
		"tcpmon-dataf-11.zst",
		"tcpmon-dataf-12",
		"tcpmon-dataf-13.zst",
		"tcpmon-dataf-14.zst",
		"tcpmon-dataf-15.zst",
		"tcpmon-dataf-16.zst",
		"tcpmon-dataf-17.zst",
		"tcpmon-dataf-18.zst",
		"tcpmon-dataf-19.zst",
		"tcpmon-dataf-20",
	}
	for _, name := range fileNames {
		_, err := s.fs.Create(filepath.Join(s.baseDir, name))
		s.Require().NoError(err)
	}

	cfg := storage.NewConfig(s.baseDir).
		WithFs(s.fs).
		WithMaxSize(10 * (1 << 20)).
		WithMaxEntriesPerFile(3)

	ds, err := storage.NewDataStore(cfg)
	s.Require().NoError(err)
	defer ds.Close()

	last := ds.GetLatestFileNo()
	s.Require().Equal(uint32(20+1), last) // when create the datastore, a new file is created
}

func randBuf(size int) []byte {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate new random buffer")
	}

	return buf
}
