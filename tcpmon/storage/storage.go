package storage

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/klauspost/compress/zstd"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/sasha-s/go-deadlock"
	"github.com/spf13/afero"
)

const DataFilePrefix = "tcpmon-dataf-"
const SealFileSuffix = ".zst"
const Version = uint16(0xadde)
const HeaderSize = 6

type DataStore struct {
	baseDir string   // database base directory
	fs      afero.Fs // filesystem interface
	config  Config

	lastFileNum    uint32     // last file num
	writerFile     afero.File // current file
	writerFilePath string     // current file path
	writerCapacity uint32     // current file capacity

	// mutex Multiple goroutines may access this datastore
	// e.g. HTTP server, monitor write
	mutex deadlock.Mutex
}

func ensureDir(dir string, fs afero.Fs) {
	_, err := fs.Open(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = fs.MkdirAll(dir, 0755)
			if err != nil {
				log.Fatal().Err(err).Str("baseDir", dir).
					Msg("Create dir failed")
			}
		} else {
			log.Fatal().Err(err).Str("baseDir", dir).
				Msg("Unable to determine base dir exists or not")
		}
	}
}

func NewDataStore(config *Config) (*DataStore, error) {
	if config == nil {
		log.Fatal().Msg("Config is nil")
		return nil, nil // make linter happy
	}

	s := &DataStore{
		baseDir: config.BaseDir,
		fs:      config.Fs,
		config:  *config,

		lastFileNum:    0,
		writerFile:     nil,
		writerFilePath: "",
		writerCapacity: 0,
	}

	ensureDir(s.baseDir, s.fs)

	err := s.doNextFile()
	if err != nil {
		return nil, errors.Wrap(err, "go to next file failed")
	}

	return s, nil
}

func (ds *DataStore) Close() error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	return ds.doClose()
}

func (ds *DataStore) doClose() error {
	if ds.writerFile != nil {
		err := ds.writerFile.Sync()
		if err != nil {
			return errors.Wrap(err, "file sync failed")
		}

		err = ds.writerFile.Close()
		if err != nil {
			return errors.Wrap(err, "file close failed")
		}

		ds.writerCapacity = 0
		ds.writerFilePath = ""
	}
	return nil
}

func (ds *DataStore) BaseDir() string {
	return ds.baseDir
}

func (ds *DataStore) Put(value []byte) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	var err error
	header := ds.newHeader(len(value))

	_, err = ds.writerFile.Write(header)
	if err != nil {
		return errors.Wrap(err, "write header failed")
	}

	_, err = ds.writerFile.Write(value)
	if err != nil {
		return errors.Wrap(err, "write failed")
	}

	ds.writerCapacity++

	if ds.writerCapacity >= ds.config.MaxEntriesPerFile {
		err = ds.doNextFile()
		if err != nil {
			return errors.Wrap(err, "rotate file failed")
		}
	}

	return nil
}

func getFileNo(fileName string) uint32 {
	fileName = strings.TrimSuffix(fileName, SealFileSuffix)

	p := strings.LastIndex(fileName, "-")
	if p == -1 {
		log.Fatal().Str("file", fileName).Msg("Invalid file name. Maybe it shouldn't be here")
	}

	numS := fileName[p+1:]
	num, err := strconv.ParseUint(numS, 10, 32)
	if err != nil {
		log.Fatal().Err(err).Str("num", numS).Msg("Invalid seq number in the file name")
	}

	return uint32(num)
}

func (ds *DataStore) GetLatestFileNo() uint32 {
	baseDir, err := ds.fs.Open(ds.baseDir)
	if err != nil {
		log.Fatal().Err(err).Msg("Open base dir failed")
	}

	files, err := baseDir.Readdirnames(-1)
	if err != nil {
		log.Fatal().Err(err).Msg("List files in base dir failed")
	}

	if len(files) == 0 {
		return 0
	}

	files = lo.Map(files, func(f string, i int) string {
		return strings.TrimSuffix(f, SealFileSuffix)
	})

	sort.Slice(files, func(i, j int) bool {
		return getFileNo(files[i]) < getFileNo(files[j])
	})
	lastFile := files[len(files)-1]

	return getFileNo(lastFile)
}

func (ds *DataStore) TotalSize() (int64, error) {
	baseDir, err := ds.fs.Open(ds.baseDir)
	if err != nil {
		return -1, errors.Wrap(err, "open base dir failed")
	}

	files, err := baseDir.Readdir(-1)
	if err != nil {
		return -1, errors.Wrap(err, "list files failed")
	}

	files = lo.Filter(files, func(f os.FileInfo, index int) bool {
		return strings.HasPrefix(f.Name(), DataFilePrefix)
	})

	return lo.SumBy(files, func(file os.FileInfo) int64 {
		return file.Size()
	}), nil
}

func (ds *DataStore) NextFile() error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	return ds.doNextFile()
}

func (ds *DataStore) doNextFile() error {
	if ds.lastFileNum == 0 {
		ds.lastFileNum = ds.GetLatestFileNo() + 1
	}

	err := ds.seal()
	if err != nil {
		return errors.Wrap(err, "close current file failed")
	}

	nextFilePath := filepath.Join(ds.baseDir, fmt.Sprintf("%s%d", DataFilePrefix, ds.lastFileNum))
	fh, err := ds.fs.Create(nextFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Create new file failed")
	}
	ds.lastFileNum++
	log.Info().Str("nextFilePath", nextFilePath).Msg("Next file created")

	ds.writerFile = fh
	ds.writerFilePath = nextFilePath
	ds.writerCapacity = 0

	ds.reclaim()

	return nil
}

func (ds *DataStore) seal() error {
	lastFileName := ds.writerFilePath

	err := ds.doClose()
	if err != nil {
		return err
	}

	if lastFileName != "" {
		err = ds.compressFile(lastFileName, lastFileName+SealFileSuffix)
		if err != nil {
			log.Warn().Err(err).Str("file", lastFileName).Msg("seal failed")
		}

		ds.reclaim()
	}

	return nil
}

func (ds *DataStore) getDataFiles(suffix string) ([]string, error) {
	baseDir, err := ds.fs.Open(ds.baseDir)
	if err != nil {
		return nil, errors.Wrap(err, "open base dir failed")
	}

	files, err := baseDir.Readdir(-1)
	if err != nil {
		return nil, errors.Wrap(err, "list files failed")
	}

	files = lo.Filter(files, func(f os.FileInfo, i int) bool {
		return strings.HasPrefix(f.Name(), DataFilePrefix)
	})

	if suffix != "" {
		files = lo.Filter(files, func(f os.FileInfo, i int) bool {
			return strings.HasSuffix(f.Name(), suffix)
		})
	}

	return lo.Map(files, func(f os.FileInfo, i int) string {
		return filepath.Join(ds.baseDir, f.Name())
	}), nil
}

func (ds *DataStore) reclaim() {
	size, err := ds.TotalSize()
	if err != nil {
		log.Warn().Err(err).Msg("Retrieve total size failed")
	}

	if size > ds.config.MaxSize {
		log.Info().Int64("size", size).
			Int64("maxSize", ds.config.MaxSize).Msg("Reclaiming...")

		// compress all data files
		dataFiles, err := ds.getDataFiles("")
		if err != nil {
			log.Fatal().Err(err).Msg("List data files failed")
		}
		rawFiles := lo.Filter(dataFiles, func(f string, i int) bool {
			return !strings.HasSuffix(f, SealFileSuffix)
		})

		for _, f := range rawFiles {
			if f == ds.writerFilePath {
				continue
			}

			err := ds.compressFile(f, f+SealFileSuffix)
			if err != nil {
				log.Warn().Err(err).Msg("Compress all failed")
			}
		}

		files, err := ds.getDataFiles(SealFileSuffix)
		if err != nil {
			log.Fatal().Err(err).Msg("Get files failed")
		}

		sort.Slice(files, func(i, j int) bool {
			return getFileNo(files[i]) < getFileNo(files[j])
		})

		fileSizes := lo.Map(files, func(f string, i int) int64 {
			info, err := os.Stat(f)
			if err != nil {
				log.Fatal().Err(err).Str("file", f).Msg("stat failed")
			}
			return info.Size()
		})

		i := 0
		for i < len(fileSizes) && size > ds.config.MaxSize {
			size -= fileSizes[i]
			i++
		}

		toDelete := files[0:i]
		for _, file := range toDelete {
			err = ds.fs.Remove(file)
			if err != nil {
				log.Warn().Err(err).Str("file", file).Msg("Delete failed")
			}
			log.Info().Str("file", file).Msg("Deleted")
		}
	}
}

func (ds *DataStore) newHeader(size int) []byte {
	buf := make([]byte, HeaderSize)
	binary.LittleEndian.PutUint16(buf[0:2], Version)      // version
	binary.LittleEndian.PutUint32(buf[2:6], uint32(size)) // body size
	return buf
}

func (ds *DataStore) compressFile(src string, dst string) error {
	in, err := ds.fs.Open(src)
	if err != nil {
		return errors.Wrap(err, "open input file failed")
	}

	out, err := ds.fs.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "open output file failed")
	}

	err = ds.compress(in, out)
	if err != nil {
		return errors.Wrap(err, "compress file failed")
	}

	err = out.Close()
	if err != nil {
		return errors.Wrap(err, "close output file failed")
	}

	err = in.Close()
	if err != nil {
		return errors.Wrap(err, "close input file failed")
	}

	err = ds.fs.RemoveAll(src)
	if err != nil {
		return errors.Wrap(err, "delete input file failed")
	}

	log.Info().Str("dst", dst).Msg("Compressed")
	return nil
}

func (ds *DataStore) compress(reader io.Reader, writer io.Writer) error {
	z, err := zstd.NewWriter(writer)
	if err != nil {
		return errors.Wrap(err, "create new zstd writer failed")
	}
	defer z.Close()

	_, err = io.Copy(z, reader)
	if err != nil {
		return errors.Wrap(err, "copy to writer failed")
	}

	return nil
}

type MetricContext struct {
	Value []byte
}
