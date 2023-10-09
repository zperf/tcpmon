package v2

import "github.com/spf13/afero"

type Config struct {
	BaseDir           string
	Fs                afero.Fs
	MaxSize           int64
	MaxEntriesPerFile uint32
}

func NewConfig(baseDir string) *Config {
	return &Config{
		BaseDir:           baseDir,
		Fs:                afero.NewOsFs(),
		MaxSize:           100 * (1 << 20),
		MaxEntriesPerFile: 500,
	}
}

func (c *Config) WithFs(fs afero.Fs) *Config {
	c.Fs = fs
	return c
}

// WithMaxSize set the max size for storage, in bytes
func (c *Config) WithMaxSize(size int64) *Config {
	c.MaxSize = size
	return c
}

// WithMaxEntriesPerFile set the max entries count per file
func (c *Config) WithMaxEntriesPerFile(entries uint32) *Config {
	c.MaxEntriesPerFile = entries
	return c
}
