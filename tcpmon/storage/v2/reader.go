package v2

import (
	"encoding/binary"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/klauspost/compress/zstd"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/afero"
)

type DataStoreReader struct {
	baseDir afero.File
	fs      afero.Fs
}

func NewDataStoreReader(baseDir string, fs afero.Fs) (*DataStoreReader, error) {
	if fs == nil {
		fs = afero.NewOsFs()
	}

	fh, err := fs.Open(baseDir)
	if err != nil {
		return nil, errors.Wrap(err, "open base dir failed")
	}

	r := &DataStoreReader{
		baseDir: fh,
		fs:      fs,
	}
	return r, nil
}

func (r *DataStoreReader) Iterate(cb func(buf []byte)) error {
	files, err := r.baseDir.Readdirnames(-1)
	if err != nil {
		return errors.Wrap(err, "list files in base dir failed")
	}

	files = lo.Filter(files, func(f string, i int) bool {
		return strings.HasPrefix(f, FilePrefix)
	})

	sort.Slice(files, func(i, j int) bool {
		return strings.Compare(files[i], files[j]) < 0
	})

	for _, file := range files {
		filePath := filepath.Join(r.baseDir.Name(), file)
		log.Info().Str("file", filePath).Msg("Iterate over file")

		reader, err := NewReader(filePath, r.fs)
		if err != nil {
			return err
		}

		off := 0
		for {
			buf, err := reader.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return err
			}

			cb(buf)
			off += len(buf)
		}
	}

	return nil
}

func (r *DataStoreReader) Count() (int, error) {
	count := 0
	err := r.Iterate(func(_ []byte) {
		count++
	})
	return count, err
}

func (r *DataStoreReader) Close() {
	if r.baseDir != nil {
		err := r.baseDir.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Close base dir failed")
		}
	}
}

type DataFileReader struct {
	fh     afero.File
	reader io.Reader
}

func NewReader(filePath string, fs afero.Fs) (*DataFileReader, error) {
	if fs == nil {
		fs = afero.NewOsFs()
	}

	fh, err := fs.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "open file failed")
	}

	if strings.HasSuffix(filePath, SealFileSuffix) {
		reader, err := zstd.NewReader(fh)
		if err != nil {
			return nil, errors.Wrap(err, "create new zstd reader failed")
		}
		return &DataFileReader{fh: fh, reader: reader}, nil
	}

	return &DataFileReader{fh: fh, reader: fh}, nil
}

func (r *DataFileReader) Read() ([]byte, error) {
	size, err := r.ReadHeader()
	if err != nil {
		return nil, err
	}
	log.Info().Uint32("size", size).Msg("Read header")

	buf := make([]byte, size)
	_, err = r.reader.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (r *DataFileReader) ReadHeader() (uint32, error) {
	buf := make([]byte, 6)
	_, err := r.reader.Read(buf)
	if err != nil {
		if err == io.EOF {
			return 0, err
		}
		return 0, errors.Wrap(err, "read header failed")
	}

	version := binary.LittleEndian.Uint16(buf[0:2])
	if version != Version {
		return 0, errors.Newf("invalid version `%d` in header", version)
	}

	size := binary.LittleEndian.Uint32(buf[2:6])
	return size, nil
}

func (r *DataFileReader) Close() {
	if r.fh != nil {
		err := r.fh.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Close file failed")
		}
	}

	if c, ok := r.reader.(io.Closer); ok {
		err := c.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Close reader failed")
		}
	}
}
