package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"report-manager/db/defi_fund"
	"report-manager/model"
	"report-manager/proxy"
	"report-manager/server/http/respond"
	"time"
)

type GetDefiFundQuotaReq struct {
	Date string `form:"date" binding:"required"`
}

type GetDefiFundQuotaResp struct {
	Quota []model.Quota `json:"quota"`
}

func GetDefiFundQuota(c *gin.Context) {
	date := c.Query("date")
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		respond.BadRequest(c, http.StatusBadRequest, err.Error())
		return
	}
	quota, err := defi_fund.FundQuota{}.QueryByDate(date)
	if err != nil {
		respond.InternalError(c, err)
		return
	}

	respond.Success(c, GetDefiFundQuotaResp{Quota: defi_fund.FundQuota{}.ToEntities(quota)})
}

type GetDefiFundPlatformSnapshotReq struct {
	Date string `form:"date" binding:"required"`
}

func GetDefiFundPlatformSnapshot(c *gin.Context) {
	date := c.Query("date")
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		respond.BadRequest(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := proxy.GetPlatformSnapshot(date)
	if err != nil {
		respond.InternalError(c, err)
		return
	}

	respond.Success(c, resp)
}
