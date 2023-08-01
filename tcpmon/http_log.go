package tcpmon

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func httpLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		now := time.Now()
		ctx.Next()
		elapsed := time.Since(now)

		code := ctx.Writer.Status()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		if query != "" {
			path = fmt.Sprintf("%s?%s", path, query)
		}
		msg := ctx.Errors.String()
		if msg == "" {
			msg = ctx.Request.Method
		}

		switch {
		case code >= 400 && code < 500:
			log.Warn().
				Str("path", path).
				Dur("lat", elapsed).
				Int("code", code).
				Msg(msg)
		case code >= 500:
			log.Error().
				Str("path", path).
				Dur("lat", elapsed).
				Int("code", code).
				Msg(msg)
		default:
			log.Info().
				Str("path", path).
				Dur("lat", elapsed).
				Int("code", code).
				Msg(msg)
		}
	}
}
