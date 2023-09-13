package tcpmon

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"

	"github.com/zperf/tcpmon/tcpmon/tutils"
)

func GetMember(q *Quorum) func(c *gin.Context) {
	return func(c *gin.Context) {
		members := make(map[string]any)
		for _, member := range q.Members() {
			addr := member.Address()
			meta, err := q.GetMemberMeta(addr)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, tutils.ErrorJSON(err))
				return
			}
			members[member.Address()] = meta
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

		member := strings.TrimSpace(string(buf))
		members := make(map[string]string)
		members[member] = ""

		_, err = q.TryJoin(members)
		if err != nil {
			c.JSON(http.StatusOK, tutils.ErrorJSON(errors.WithStack(err)))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"error": nil,
		})
	}
}

func LeaveCluster(q *Quorum) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := q.Leave(3 * time.Second)
		c.JSON(http.StatusOK, tutils.ErrorJSON(err))
	}
}
