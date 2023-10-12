package storage

import (
	"archive/tar"
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

type Reader struct {
	baseDir afero.File
	fs      afero.Fs
	config  *ReaderConfig
}

type ReaderConfig struct {
	baseDir string
	fs      afero.Fs
	prefix  string
	suffix  string
}

func NewReaderConfig(baseDir string) *ReaderConfig {
	return &ReaderConfig{
		baseDir: baseDir,
		fs:      afero.NewOsFs(),
		prefix:  DataFilePrefix,
		suffix:  "",
	}
}

func (c *ReaderConfig) WithFs(fs afero.Fs) *ReaderConfig {
	c.fs = fs
	return c
}

func (c *ReaderConfig) WithPrefix(s string) *ReaderConfig {
	c.prefix = s
	return c
}

func (c *ReaderConfig) WithSuffix(s string) *ReaderConfig {
	c.suffix = s
	return c
}

func NewDataStoreReader(config *ReaderConfig) (*Reader, error) {
	fh, err := config.fs.Open(config.baseDir)
	if err != nil {
		return nil, errors.Wrap(err, "open base dir failed")
	}

	r := &Reader{
		baseDir: fh,
		fs:      config.fs,
		config:  config,
	}
	return r, nil
}

func (r *Reader) Close() {
	if r.baseDir != nil {
		err := r.baseDir.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Close base dir failed")
		}
	}
}

func (r *Reader) newFilter() func(f string, i int) bool {
	if r.config.prefix != "" && r.config.suffix != "" {
		return func(f string, i int) bool {
			return strings.HasPrefix(f, r.config.prefix) && strings.HasSuffix(f, r.config.suffix)
		}
	}

	if r.config.prefix != "" {
		return func(f string, i int) bool {
			return strings.HasPrefix(f, r.config.prefix)
		}
	}

	if r.config.suffix != "" {
		return func(f string, i int) bool {
			return strings.HasPrefix(f, r.config.prefix)
		}
	}

	return nil
}

func (r *Reader) files() ([]string, error) {
	files, err := r.baseDir.Readdirnames(-1)
	if err != nil {
		return nil, errors.Wrap(err, "list files in base dir failed")
	}

	filter := r.newFilter()
	if filter != nil {
		files = lo.Filter(files, filter)
	}

	sort.Slice(files, func(i, j int) bool {
		return getFileNo(files[i]) < getFileNo(files[j])
	})

	return files, nil
}

func (r *Reader) Iterate(cb func(buf []byte) error) error {
	files, err := r.files()
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(r.baseDir.Name(), file)
		log.Info().Str("file", filePath).Msg("Iterate over file")

		reader, err := NewDataFileReader(filePath, r.fs)
		if err != nil {
			return err
		}

		offset := uint64(0)
		for {
			buf, err := reader.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return err
			}
			offset += uint64(len(buf))

			if len(buf) != 0 {
				err = cb(buf)
				if err != nil {
					log.Warn().Err(err).Uint64("offset", offset).
						Msg("Error occurred, skip to the next file")
					break
				}
			}
		}
	}

	return nil
}

func (r *Reader) Count() (int, error) {
	count := 0
	err := r.Iterate(func(_ []byte) error {
		count++
		return nil
	})
	return count, err
}

func (r *Reader) Package(writer io.Writer) error {
	files, err := r.files()
	if err != nil {
		return err
	}

	t := tar.NewWriter(writer)
	defer t.Close()

	for _, file := range files {
		filePath := filepath.Join(r.baseDir.Name(), file)

		fh, err := r.fs.Open(filePath)
		if err != nil {
			return errors.Wrap(err, "open file failed")
		}

		stat, err := r.fs.Stat(filePath)
		if err != nil {
			return errors.Wrap(err, "stat failed")
		}

		err = t.WriteHeader(&tar.Header{
			Name:    file,
			Mode:    0644,
			Size:    stat.Size(),
			ModTime: stat.ModTime(),
			Format:  tar.FormatPAX,
		})
		if err != nil {
			return errors.Wrap(err, "write header failed")
		}

		_, err = io.Copy(t, fh)
		if err != nil {
			return errors.Wrap(err, "copy file failed")
		}
	}

	return nil
}

type DataFileReader struct {
	fh     afero.File
	reader io.Reader
}

func NewDataFileReader(filePath string, fs afero.Fs) (*DataFileReader, error) {
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
		return 0, errors.Newf("invalid version 0x%x in header", version)
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
