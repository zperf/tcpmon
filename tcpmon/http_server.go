package tcpmon

import (
	"net/http"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitRoutes(router *gin.Engine, ds *Datastore) {
	router.GET("/last", GetLast(ds))
}

func GetLast(ds *Datastore) func(c *gin.Context) {
	return func(c *gin.Context) {
		// batch size
		var batch int
		const batchDefault = 10
		batchStr := strings.TrimSpace(c.Query("batch"))
		if batchStr == "" {
			batch = batchDefault
		} else {
			batch, _ = ParseInt(batchStr)
			if batch <= 0 {
				batch = batchDefault
			}
		}

		// return values or not
		returnValueStr := strings.TrimSpace(c.Query("value"))
		returnValue := ParseBool(returnValueStr)

		// prefix filter
		prefixFilter := strings.TrimSpace(c.Query("prefix"))

		keys := ds.LastKeys(batch)
		if prefixFilter != "" {
			keys = FilterStringsByPrefix(keys, prefixFilter)
		}

		if !returnValue {
			c.JSON(http.StatusOK, gin.H{
				"window": keys,
			})
			return
		}

		values, err := ds.GetBatch(keys)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		pairs := make([]gin.H, len(keys))
		for i, key := range keys {
			pairs[i] = gin.H{
				"key":   key,
				"value": values[i],
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"window": pairs,
		})
	}
}

func (mon *Monitor) startHttpServer(addr string) {
	verbose := viper.GetBool("verbose")
	if !verbose {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(httpLogger())
	engine.Use(gin.Recovery())
	InitRoutes(engine, mon.datastore)

	mon.httpServer = &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	go func(srv *http.Server) {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to listen and serve")
		}
	}(mon.httpServer)
}
