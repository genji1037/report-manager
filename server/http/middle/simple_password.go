package middle

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"report-manager/server/http/respond"
)

func SimplePassword(c *gin.Context) {
	if c.GetHeader("Password") != "ba9b89sbs9yys9bys9bd8" {
		respond.Error(c, http.StatusUnauthorized, http.StatusUnauthorized, "unAuthorized")
		c.Abort()
		return
	}
	c.Next()
}
