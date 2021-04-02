package http

import (
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"report-manager/logger"
	"report-manager/metrics"
	"report-manager/server/http/middle"

	"github.com/gin-gonic/gin"
	"report-manager/server/http/handler"
)

func Run(host string, port int) {
	router := gin.Default()

	ginpprof.Wrap(router)

	metrics.InitCustomerMetrics(router)

	// 报告
	router.POST("/report/:name", handler.DoReport)

	testGroup := router.Group("/test")
	testGroup.Use(middle.SimplePassword)
	{
		testGroup.POST("/otc_locked_token", handler.PersistOTCLockedToken)
		testGroup.POST("/ctc_locked_token", handler.PersistCTCLockedToken)
		testGroup.POST("/sie_count/sugar", handler.SIECountSugar)
		testGroup.POST("/sie_count/n_one_buy", handler.SIECountNOneBuy)
		testGroup.POST("/sie_count/shop_destroy", handler.SIECountShopDestroy)
	}

	internalGroup := router.Group("/internal")
	{
		internalGroup.POST("/wallet/balance", handler.ReceiveWalletBalance)
		internalGroup.POST("/n_one/extra", handler.ReceiveNOneExtra)
	}

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/count", handler.GetSIECount)
	}

	err := router.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		logger.Errorf("router Run failed: %s", err.Error())
	}
}
