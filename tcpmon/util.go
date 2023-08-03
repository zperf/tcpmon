package tcpmon

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func randUint64() (uint64, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(b[:]), nil
}

func ParseUint32(s string) (uint32, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(val), nil
}

func ParseInt(s string) (int, error) {
	val, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

func ParseBool(s string) bool {
	return s == "1" || s == "true" || s == "on" || s == "True" || s == "On"
}

func HasPrefix(prefix string, buf []byte) bool {
	return len(buf) >= len(prefix) && prefix == string(buf[:len(prefix)])
}

func ToProtojson(m proto.Message) map[string]any {
	buf, err := protojson.Marshal(m)
	if err != nil {
		return map[string]any{"error": err}
	}
	var val map[string]any
	err = json.Unmarshal(buf, &val)
	if err != nil {
		return map[string]any{"error": err}
	}
	return val
}

func ErrorStr(err error) string {
	return fmt.Sprintf("error: %v", err)
}

func ErrorJSON(err error) map[string]any {
	return map[string]any{"error": err.Error()}
}

func FilterStringsByPrefix(s []string, prefix string) []string {
	r := make([]string, 0)
	for _, v := range s {
		if strings.HasPrefix(v, prefix) {
			r = append(r, v)
		}
	}
	return r
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
