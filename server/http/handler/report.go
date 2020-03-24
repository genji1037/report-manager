package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"report-manager/server/http/respond"
	"report-manager/service"
)

func DoReport(c *gin.Context) {
	name := c.Param("name")
	switch name {
	case "exchange_data_report":
		err := service.ExchangeReport()
		if err != nil {
			respond.Error(c, http.StatusInternalServerError, http.StatusInternalServerError, err.Error())
			return
		}
		respond.Success(c, "ok")
	case "mall_destroy_failed_report":
		err := service.MallDestroyFailedReport()
		if err != nil {
			respond.Error(c, http.StatusInternalServerError, http.StatusInternalServerError, err.Error())
			return
		}
		respond.Success(c, "ok")
	default:
		respond.Error(c, http.StatusBadRequest, http.StatusBadRequest, "report not found")
	}
}
