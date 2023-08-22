package tcpmon

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
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

func ParseUint64(s string) (uint64, error) {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func ParseFloat64(s string) (float64, error) {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
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

type IpAddr struct {
	Address string
	Port    int
}

func ParseIpAddr(s string) *IpAddr {
	p := strings.Index(s, ":")
	if p == -1 {
		return nil
	}

	port, err := strconv.ParseInt(s[p+1:], 10, 32)
	if err != nil {
		return nil
	}

	return &IpAddr{
		Address: s[:p],
		Port:    int(port),
	}
}

func (a *IpAddr) String() string {
	return fmt.Sprintf("%s:%d", a.Address, a.Port)
}
