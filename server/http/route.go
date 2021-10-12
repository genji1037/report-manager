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
		// curl -XPOST http://localhost:18096/test/sie_count/sugar -H 'Password:ba9b89sbs9yys9bys9bd8'
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

	reportAdminGroup := router.Group("/rep_admin")
	{
		reportAdminGroup.POST("/special_user/create", handler.CreateSpecialUser)
		reportAdminGroup.POST("/special_user/delete", handler.DeleteSpecialUser)
		reportAdminGroup.POST("/special_user/update", handler.UpdateSpecialUser)
		reportAdminGroup.GET("/special_users", handler.ListSpecialUsers)
		reportAdminGroup.GET("/special_user/report", handler.GetSpecialUserReport)

		reportAdminGroup.GET("/sugars", handler.GetSugars)
		reportAdminGroup.GET("/defi_fund/quota", handler.GetDefiFundQuota)
		reportAdminGroup.GET("/defi_fund/platform_snapshot", handler.GetDefiFundPlatformSnapshot)
		reportAdminGroup.GET("/secret_chain/pledge", handler.GetSecretChainPledge)
		reportAdminGroup.GET("/n_one/out", handler.GetNOneOut)
		reportAdminGroup.GET("/ssns/report", handler.GetSSNSReport)
	}

	err := router.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		logger.Errorf("router Run failed: %s", err.Error())
	}
}
