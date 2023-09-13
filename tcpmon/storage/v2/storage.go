package v2

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
	"github.com/spf13/afero"
)

const FilePrefix = "tcpmon-dataf-"
const SealFileSuffix = ".zst"
const Version = uint16(0xadde)

type DataStore struct {
	baseDir string   // database base directory
	fs      afero.Fs // filesystem interface
	config  Config

	lastFileNum    uint32     // last file num
	writerFile     afero.File // current file
	writerFilePath string     // current file path
	writerCapacity uint32     // current file capacity
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

	err := s.NextFile()
	if err != nil {
		return nil, errors.Wrap(err, "go to next file failed")
	}

	return s, nil
}

func (ds *DataStore) Close() error {
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

func (ds *DataStore) Seal() error {
	lastFileName := ds.writerFilePath

	err := ds.Close()
	if err != nil {
		return err
	}

	if lastFileName != "" {
		err = ds.compressFile(lastFileName, lastFileName+SealFileSuffix)
		if err != nil {
			log.Warn().Err(err).Str("file", lastFileName).Msg("Seal failed")
		}

		ds.Reclaim()
	}

	return nil
}

func (ds *DataStore) Put(value []byte) error {
	log.Info().Int("bufLen", len(value)).Msg("Put")

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
		err = ds.NextFile()
		if err != nil {
			return errors.Wrap(err, "rotate file failed")
		}
	}

	return nil
}

func (ds *DataStore) NextFile() error {
	if ds.lastFileNum == 0 {
		ds.lastFileNum = ds.GetLatestFileNo() + 1
	}

	err := ds.Seal()
	if err != nil {
		return errors.Wrap(err, "close current file failed")
	}

	nextFilePath := filepath.Join(ds.baseDir, fmt.Sprintf("%s%d", FilePrefix, ds.lastFileNum))
	fh, err := ds.fs.Create(nextFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Create new file failed")
	}
	ds.lastFileNum++
	log.Info().Str("nextFilePath", nextFilePath).Msg("Next file created")

	ds.writerFile = fh
	ds.writerFilePath = nextFilePath
	ds.writerCapacity = 0

	if ds.lastFileNum%ds.config.ReclaimAt == 0 {
		ds.Reclaim()
	}

	return nil
}

func (ds *DataStore) TotalSize() (int64, []os.FileInfo, error) {
	baseDir, err := ds.fs.Open(ds.baseDir)
	if err != nil {
		return -1, nil, errors.Wrap(err, "open base dir failed")
	}

	files, err := baseDir.Readdir(-1)
	if err != nil {
		return -1, nil, errors.Wrap(err, "list files failed")
	}

	files = lo.Filter(files, func(f os.FileInfo, index int) bool {
		return strings.HasSuffix(f.Name(), SealFileSuffix)
	})

	return lo.SumBy(files, func(file os.FileInfo) int64 {
		return file.Size()
	}), files, nil
}

func (ds *DataStore) Reclaim() {
	size, files, err := ds.TotalSize()
	if err != nil {
		log.Warn().Err(err).Msg("Retrieve total size failed")
	}

	if size > ds.config.MaxSize {
		log.Info().Int64("size", size).
			Int64("maxSize", ds.config.MaxSize).Msg("Reclaiming...")

		sort.Slice(files, func(i, j int) bool {
			return strings.Compare(files[i].Name(), files[j].Name()) < 0
		})

		i := 0
		for size > ds.config.MaxSize {
			i++
			size -= files[i].Size()
		}

		toDelete := files[0:i]
		for _, file := range toDelete {
			filePath := filepath.Join(ds.baseDir, file.Name())
			err = ds.fs.Remove(filePath)
			if err != nil {
				log.Warn().Err(err).Str("file", filePath).Msg("Delete failed")
			}
			log.Info().Str("file", file.Name()).Msg("Deleted")
		}
	}
}

func (ds *DataStore) GetLatestFileNo() uint32 {
	baseDir, err := ds.fs.Open(ds.baseDir)
	if err != nil {
		log.Fatal().Err(err).Msg("Open base dir failed")
	}

	fileNames, err := baseDir.Readdirnames(-1)
	if err != nil {
		log.Fatal().Err(err).Msg("List files in base dir failed")
	}

	if len(fileNames) == 0 {
		return 0
	}

	sort.Strings(fileNames)
	lastFile := fileNames[len(fileNames)-1]

	p := strings.LastIndex(lastFile, "-")
	if p == -1 {
		log.Fatal().Str("file", lastFile).Msg("Invalid file name. Maybe it shouldn't be here")
	}

	numS := lastFile[p+1:]
	num, err := strconv.ParseUint(numS, 10, 32)
	if err != nil {
		log.Fatal().Err(err).Str("num", numS).Msg("Invalid seq number in the file name")
	}

	return uint32(num)
}

func (ds *DataStore) newHeader(size int) []byte {
	buf := make([]byte, 6)
	binary.LittleEndian.PutUint16(buf[0:2], Version)      // version
	binary.LittleEndian.PutUint32(buf[2:6], uint32(size)) // body size
	return buf
}

func (ds *DataStore) compressFile(src string, dst string) error {
	log.Info().Str("src", src).Msg("Compressing")
	defer log.Info().Str("dst", dst).Msg("Compressed")

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
