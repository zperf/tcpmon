package tcpmon

import (
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig for logging
type LogConfig struct {
	// Log level
	Level zerolog.Level

	// Enable console logging
	ConsoleLoggingEnabled bool

	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false
	FileLoggingEnabled bool

	// Directory to log to when file logging is enabled
	Directory string

	// Filename is the name of the logfile which will be placed inside the directory
	Filename string

	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int

	// MaxBackups the max number of rolled files to keep
	MaxBackups int

	// MaxAge the max age in days to keep a logfile
	MaxAge int
}

func InitLogger(config *LogConfig) {
	var writers []io.Writer
	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339Nano,
		})
	}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}

	log.Logger = zerolog.New(io.MultiWriter(writers...)).
		Level(config.Level).
		With().
		Timestamp().
		Caller().
		Logger()
}

func newRollingFile(config *LogConfig) io.Writer {
	logger := &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
		LocalTime:  true,              // use local time, default UTC time
	}
	_ = logger.Rotate()
	return logger
}

type BadgerLogger struct{}

func (b *BadgerLogger) Errorf(format string, args ...interface{}) {
	log.Error().Str("mod", "badger").Msgf(strings.TrimSuffix(format, "\n"), args...)
}

func (b *BadgerLogger) Warningf(format string, args ...interface{}) {
	log.Warn().Str("mod", "badger").Msgf(strings.TrimSuffix(format, "\n"), args...)
}

func (b *BadgerLogger) Infof(format string, args ...interface{}) {
	log.Info().Str("mod", "badger").Msgf(strings.TrimSuffix(format, "\n"), args...)
}

func (b *BadgerLogger) Debugf(format string, args ...interface{}) {
	log.Debug().Str("mod", "badger").Msgf(strings.TrimSuffix(format, "\n"), args...)
}

type MemberlistLogger struct{}

func (m *MemberlistLogger) Write(p []byte) (int, error) {
	log.Info().Str("mod", "memberlist").Msg(strings.TrimSpace(string(p)))
	return len(p), nil
}
