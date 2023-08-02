package tcpmon

import (
	"net/http"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitRoutes(router *gin.Engine, mon *Monitor) {
	ds := mon.datastore
	router.GET("/last", GetLast(ds))
	router.GET("/metrics", GetMetrics(ds))
	router.GET("/metrics/:type", GetMetrics(ds))
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

func GetMetrics(ds *Datastore) func(c *gin.Context) {
	return func(c *gin.Context) {
		kind := c.Param("type")
		if !ValidPrefix(kind) && kind != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Newf("invalid type: %v", kind)})
			return
		}

		if kind == "" {
			// without prefix, iterate over all
			keys, err := ds.Keys()
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorJSON(err))
			}
			c.JSON(http.StatusOK, gin.H{
				"len":  len(keys),
				"keys": keys,
			})
		} else {
			// with prefix
			pairs, err := ds.GetPrefix([]byte(kind), 10, true)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorJSON(errors.WithStack(err)))
				return
			}

			buf := make([]gin.H, 0)
			for _, p := range pairs {
				buf = append(buf, p.ToJSON())
			}
			c.JSON(http.StatusOK, buf)
		}
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
	InitRoutes(engine, mon)

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
