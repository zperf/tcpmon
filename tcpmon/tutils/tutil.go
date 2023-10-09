package tutils

import (
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
)

func FatalIf(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal error")
	}
}

func IsDirEmpty(p string) (bool, error) {
	fh, err := os.Open(p)
	if err != nil {
		return false, err
	}
	defer fh.Close()

	_, err = fh.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func Btoi(x bool) uint32 {
	if x {
		return 1
	}
	return 0
}

func ParseUint32(s string) (uint32, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		log.Warn().Err(errors.WithStack(err)).Msg("Parse failed")
		return 0, err
	}
	return uint32(val), nil
}

func ParseInt64(s string) (int64, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Warn().Err(errors.WithStack(err)).Msg("Parse failed")
		return 0, err
	}
	return val, nil
}

func ParseUint64(s string) (uint64, error) {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Warn().Err(errors.WithStack(err)).Msg("Parse failed")
		return 0, err
	}
	return val, nil
}

func ParseFloat64(s string) (float64, error) {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Warn().Err(errors.WithStack(err)).Msg("Parse failed")
		return 0, err
	}
	return val, nil
}

func ErrorJSON(err error) map[string]any {
	return map[string]any{"error": err.Error()}
}

func Hostname() string {
	name, err := os.Hostname()
	if err != nil {
		log.Warn().Err(err).Msg("Get hostname failed")
		return ""
	}
	return name
}

var reFilenameFilter = regexp.MustCompile(`[\\/:*?"<>|]`)

func SafeFilename(filename string) string {
	return reFilenameFilter.ReplaceAllString(filename, "_")
}

func GetIpFromAddress(s string) string {
	p := strings.Index(s, ":")
	if p == -1 {
		return s
	}
	return s[:p]
}

func FileExists(s string) (bool, error) {
	_, err := os.Stat(s)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return false, err
	}
	return true, nil
}

func SplitNewline(c rune) bool {
	return c == '\n'
}

func SplitSpace(c rune) bool {
	return c == '\n' || c == '\r' || c == '\t' || c == ' '
}

func FileFallback(path ...string) string {
	for _, p := range path {
		ok, err := FileExists(p)
		if err != nil {
			log.Fatal().Err(err).Str("file", p).Msg("Stat file failed")
		}

		if ok {
			return p
		}
	}
	return ""
}
