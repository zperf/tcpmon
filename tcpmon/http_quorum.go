package tcpmon

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
)

func GetMember(cluster *Quorum) func(c *gin.Context) {
	return func(c *gin.Context) {
		members := make([]string, 0)
		for _, member := range cluster.Members() {
			members = append(members, member.Address())
		}
		c.JSON(http.StatusOK, gin.H{
			"len":     len(members),
			"members": members,
		})
	}
}

func JoinCluster(q *Quorum) func(c *gin.Context) {
	return func(c *gin.Context) {
		buf, err := io.ReadAll(c.Request.Body)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, errors.WithStack(err))
			return
		}

		clusterIPAddr := strings.TrimSpace(string(buf))

		_, err = q.TryJoin([]string{clusterIPAddr})
		if err != nil {
			c.JSON(http.StatusOK, ErrorJSON(errors.WithStack(err)))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"error": nil,
		})
	}
}

func LeaveCluster(q *Quorum) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := q.mlist.Leave(3 * time.Second)
		c.JSON(http.StatusOK, ErrorJSON(err))
	}
}
