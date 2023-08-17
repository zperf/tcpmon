package tcpmon

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitRoutes(router *gin.Engine, mon *Monitor) {
	ds := mon.datastore
	gs := mon.gossipServer
	router.GET("/metrics", GetMetrics(ds))
	router.GET("/metrics/:type", GetMetrics(ds))
	router.GET("/backup", GetBackup(ds))
	router.GET("/members", GetMember(gs))
	pprof.Register(router)
}

func GetMember(gs *GossipServer) func(c *gin.Context) {
	return func(c *gin.Context) {
		addrs := []string{}
		for _, member := range gs.cluster.Members() {
			addrs = append(addrs, member.Address())
		}
		c.JSON(http.StatusOK, gin.H{
			"len":     len(gs.cluster.Members()),
			"members": addrs,
		})
	}
}

func GetBackup(ds *Datastore) func(c *gin.Context) {
	hostname := Hostname()
	filename := SafeFilename(fmt.Sprintf("tcpmon-dump-%s-%s.bak", hostname, time.Now().Format(time.DateOnly)))

	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		_, err := ds.Backup(c.Writer, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorJSON(err))
			return
		}
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
			keys, err := ds.GetKeys()
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
