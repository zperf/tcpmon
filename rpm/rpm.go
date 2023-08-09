package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/rpmpack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ReadFile(name string) []byte {
	buf, err := os.ReadFile(name)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read file")
	}
	return buf
}

func ReadFileString(name string) string {
	return strings.TrimSpace(string(ReadFile(name)))
}

func BuildRPM(arch string) {
	rpmMeta := rpmpack.RPMMetaData{
		Name:       "tcpmon",
		Version:    ReadFileString("VERSION"),
		Release:    "1.el7",
		OS:         "linux",
		Arch:       arch,
		Summary:    "Tcpmon - a portable netowrk monitor",
		Licence:    "MIT",
		Vendor:     "SMTX",
		Packager:   "SMTX",
		URL:        "https://github.com/zperf/tcpmon",
		BuildTime:  time.Now(),
		Compressor: "gzip",
		Group:      "Unspecified",
		Description: "tcpmon is a network monitoring tool that provides key insights into the system's network" +
			"interfaces and connections",
	}
	rpm, err := rpmpack.NewRPM(rpmMeta)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create rpm builder")
	}

	rpmFilename := fmt.Sprintf(
		"%s-%s-%v.%s.rpm", rpmMeta.Name, rpmMeta.Version, rpmMeta.Release, rpmMeta.Arch)

	rpm.AddFile(rpmpack.RPMFile{
		Name:  "/usr/bin/tcpmon",
		Body:  ReadFile("../bin/tcpmon-" + arch),
		Type:  rpmpack.GenericFile,
		Mode:  0775,
		Owner: "root",
		Group: "root",
	})
	rpm.AddFile(rpmpack.RPMFile{
		Name:  "/etc/tcpmon/config.yaml",
		Body:  ReadFile("tcpmon.yaml"),
		Type:  rpmpack.ConfigFile,
		Mode:  0644,
		Owner: "root",
		Group: "root",
	})
	rpm.AddFile(rpmpack.RPMFile{
		Name:  "/etc/systemd/system/tcpmon.service",
		Body:  ReadFile("tcpmon.service"),
		Type:  rpmpack.GenericFile,
		Mode:  0644,
		Owner: "root",
		Group: "root",
	})

	fh, err := os.OpenFile(rpmFilename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open file for writing")
	}
	defer func(fh *os.File) {
		err := fh.Close()
		if err != nil {
			log.Warn().Err(err).Msg("close file failed")
		}
	}(fh)

	err = rpm.Write(fh)
	if err != nil {
		log.Fatal().Err(err).Msg("write to RPM failed")
	}
}

func main() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	BuildRPM("x86_64")
	BuildRPM("aarch64")
}
