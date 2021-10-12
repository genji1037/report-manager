package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
	"report-manager/alg"
	"strings"
)

const validPrefix = "/server-secret-social-network/file"

var addr string

func main() {
	flag.StringVar(&addr, "addr", "", "")
	e := gin.New()
	e.GET(validPrefix+"/:path", func(c *gin.Context) {
		path := validPrefix + c.Param("path")
		i := strings.LastIndex(path, "/")
		fPath, err := alg.GetFilePath(path[:i], path[i+1:])
		if err != nil {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.File(fPath)
	})
	e.Run(":17078")
}
