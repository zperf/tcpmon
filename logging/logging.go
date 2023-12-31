package logging

import (
	"bytes"
	"fmt"
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
	// Level the log level
	Level zerolog.Level

	// ConsoleLoggingEnabled console logging enabled or not
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
}

func InitLogger(config *LogConfig) {
	var writers []io.Writer
	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339Nano,
			PartsOrder: []string{
				zerolog.TimestampFieldName,
				zerolog.LevelFieldName,
				zerolog.CallerFieldName,
				zerolog.MessageFieldName,
			},
			FieldsExclude: []string{
				zerolog.ErrorStackFieldName,
			},
			FormatExtra: func(m map[string]interface{}, buffer *bytes.Buffer) error {
				s, ok := m["stack"]
				if ok {
					_, err := buffer.WriteString(s.(string))
					return err
				}
				return nil
			},
		})
	}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}

	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		return fmt.Sprintf("\n%+v", err)
	}
	log.Logger = zerolog.New(io.MultiWriter(writers...)).
		Level(config.Level).
		With().
		Timestamp().
		Caller().
		Stack().
		Logger()
}

func newRollingFile(config *LogConfig) io.Writer {
	logger := &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		LocalTime:  true,              // use local time, default UTC time
	}
	_ = logger.Rotate()
	return logger
}

type MemberlistLogger struct {
	log zerolog.Logger
}

func NewMemberlistLogger() *MemberlistLogger {
	return &MemberlistLogger{
		log: log.With().Str("mod", "memberlist").Logger().Level(zerolog.WarnLevel),
	}
}

func (m *MemberlistLogger) Write(p []byte) (int, error) {
	m.log.Info().Msg(strings.TrimSpace(string(p)))
	return len(p), nil
}
