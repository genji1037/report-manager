package http

import (
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"report-manager/logger"
	"report-manager/metrics"

	"github.com/gin-gonic/gin"
	"report-manager/server/http/handler"
)

func Run(host string, port int) {
	router := gin.Default()

	ginpprof.Wrap(router)

	metrics.InitCustomerMetrics(router)

	// 根据AppID获取特定APP信息
	router.POST("/report/:name", handler.DoReport)

	err := router.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		logger.Errorf("router Run failed: %s", err.Error())
	}
}
