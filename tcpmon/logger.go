package tcpmon

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

var LevelMap = map[string]zerolog.Level{
	"TRACE":    zerolog.TraceLevel,
	"DEBUG":    zerolog.DebugLevel,
	"INFO":     zerolog.InfoLevel,
	"WARN":     zerolog.WarnLevel,
	"ERROR":    zerolog.ErrorLevel,
	"FATAL":    zerolog.FatalLevel,
	"PANIC":    zerolog.PanicLevel,
	"NO":       zerolog.NoLevel,
	"DISABLED": zerolog.Disabled,
}

// Configuration for logging
type Config struct {
	// Enable console logging
	ConsoleLoggingEnabled bool
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false
	FileLoggingEnabled bool
	// Directory to log to to when filelogging is enabled
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

func InitLogger(config Config) {
	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano})
	}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}
	mw := io.MultiWriter(writers...)

	logLevel, exist := LevelMap[strings.ToUpper(viper.GetString("log-level"))]
	if !exist {
		fmt.Println("log level not exist, please check")
		os.Exit(1)
	}

	logger := zerolog.New(mw).
		Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	logger.Info().
		Bool("fileLogging", config.FileLoggingEnabled).
		Str("logDirectory", config.Directory).
		Str("fileName", config.Filename).
		Int("maxSizeMB", config.MaxSize).
		Int("maxBackups", config.MaxBackups).
		Int("maxAgeInDays", config.MaxAge).
		Msg("logging configured")

	log.Logger = logger
}

func newRollingFile(config Config) io.Writer {
	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
		LocalTime:  true,              // use local time, default UTC time
	}
}
