package cmd

import "github.com/rs/zerolog/log"

func fatalIf(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal error")
	}
}
