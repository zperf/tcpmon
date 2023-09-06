package tcpmon

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RegisterRoutes(router *gin.Engine, mon *Monitor) {
	// router.GET("/metrics", GetMetrics(mon.datastore))
	// router.GET("/metrics/*name", GetMetric(mon.datastore))
	router.GET("/members", GetMember(mon.quorum))
	router.POST("/members", JoinCluster(mon.quorum))
	router.POST("/members/leave", LeaveCluster(mon.quorum))
	router.GET("/backup", GetBackup(mon))
	pprof.Register(router)
}

func GetBackup(mon *Monitor) func(c *gin.Context) {
	hostname := Hostname()
	filename := SafeFilename(fmt.Sprintf("tcpmon-dump-%s-%s.bak", hostname, mon.quorum.MyIP()))

	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		_, err := mon.datastore.Backup(c.Writer, 0)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorJSON(err))
			return
		}
	}
}

func GetMetric(ds *DataStore) func(c *gin.Context) {
	return func(c *gin.Context) {
		metricName := c.Param("name")
		p, err := ds.Get(strings.TrimPrefix(metricName, "/"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorJSON(err))
			return
		}
		c.JSON(http.StatusOK, p.ToJSON())
	}
}

func GetMetrics(ds *DataStore) func(c *gin.Context) {
	return func(c *gin.Context) {
		kind := c.Query("type")
		if !ValidMetricPrefix(kind) && kind != "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				ErrorJSON(errors.Newf("invalid type: %v", kind)))
			return
		}

		if kind == "" {
			// without prefix, iterate over all
			keys, err := ds.GetMetrics()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorJSON(err))
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"len":  len(keys),
				"keys": keys,
			})
		} else {
			// with prefix
			pairs, err := ds.GetPrefix([]byte(kind), 0, false)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorJSON(err))
				return
			}

			keys := make([]string, 0)
			for _, p := range pairs {
				keys = append(keys, p.Key)
			}
			c.JSON(http.StatusOK, gin.H{
				"len":  len(keys),
				"keys": keys,
			})
		}
	}
}

func (m *Monitor) startHttpServer(addr string) {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(httpLogger())
	engine.Use(gin.Recovery())
	RegisterRoutes(engine, m)

	m.httpServer = &http.Server{
		Addr:    addr,
		Handler: engine,
	}
	log.Info().Str("addr", addr).Msg("HTTP server started")

	go func(srv *http.Server) {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(errors.WithStack(err)).Msg("Serve HTTP service failed")
		}
	}(m.httpServer)
}
