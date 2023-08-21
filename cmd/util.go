package cmd

import (
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

func fatalIf(err error) {
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
