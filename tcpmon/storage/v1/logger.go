package v1

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type BadgerDbLogger struct {
	log zerolog.Logger
}

func NewBadgerLogger() *BadgerDbLogger {
	return &BadgerDbLogger{
		log: log.With().Str("mod", "badger").Logger().Level(zerolog.WarnLevel),
	}
}

func (b *BadgerDbLogger) Errorf(format string, args ...interface{}) {
	b.log.Error().Msgf(strings.TrimSpace(format), args...)
}

func (b *BadgerDbLogger) Warningf(format string, args ...interface{}) {
	b.log.Warn().Msgf(strings.TrimSpace(format), args...)
}

func (b *BadgerDbLogger) Infof(format string, args ...interface{}) {
	b.log.Info().Msgf(strings.TrimSpace(format), args...)
}

func (b *BadgerDbLogger) Debugf(format string, args ...interface{}) {
	b.log.Debug().Msgf(strings.TrimSpace(format), args...)
}
