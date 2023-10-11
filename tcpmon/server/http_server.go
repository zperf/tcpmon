package server

import (
	"fmt"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/zperf/tcpmon/tcpmon/storage"
	"github.com/zperf/tcpmon/tcpmon/tutils"
)

func RegisterRoutes(router *gin.Engine, mon *Monitor) {
	router.GET("/", GetHome)
	router.GET("/backup", GetBackup(mon))

	if mon.quorum != nil {
		router.GET("/members", GetMember(mon.quorum))
		router.POST("/members", JoinCluster(mon.quorum))
		router.POST("/members/leave", LeaveCluster(mon.quorum))
	}

	pprof.Register(router)
}

func GetHome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"service": "tcpmon"})
}

func GetBackup(mon *Monitor) func(c *gin.Context) {
	hostname := tutils.Hostname()
	filename := tutils.SafeFilename(fmt.Sprintf("tcpmon-datastore-%s.tar", hostname))

	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

		err := mon.datastore.NextFile()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, tutils.ErrorJSON(err))
			return
		}

		r, err := storage.NewDataStoreReader(storage.NewReaderConfig(mon.datastore.BaseDir()).WithSuffix(storage.SealFileSuffix))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, tutils.ErrorJSON(err))
			return
		}

		err = r.Package(c.Writer)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, tutils.ErrorJSON(err))
			return
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
