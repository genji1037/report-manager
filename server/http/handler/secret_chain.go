package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"net/http"
	"report-manager/db"
	"report-manager/db/open"
	"report-manager/logger"
	"report-manager/proxy"
	"report-manager/server/http/respond"
	"report-manager/service"
	"time"
)

type GetSecretChainPledgeResp struct {
	Date            string          `json:"date"`
	SIEVolume       decimal.Decimal `json:"sie_volume"`
	GASVolume       decimal.Decimal `json:"gas_volume"`
	GASRewardVolume decimal.Decimal `json:"gas_reward_volume"`
}

func GetSecretChainPledge(c *gin.Context) {
	date := c.Query("date")
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		respond.BadRequest(c, http.StatusBadRequest, err.Error())
		return
	}

	var cacheNotFound bool
	var snapshot db.SecretChainPledgeSnapshot
	snapshot.Dat = date
	if err := snapshot.GetByDate(); err != nil {
		if err == gorm.ErrRecordNotFound {
			cacheNotFound = true
		} else {
			respond.InternalError(c, err)
			return
		}
	}

	if cacheNotFound {
		sieSum, gasSum, err := service.GetSecretChainPledge(date)
		if err != nil {
			if err == proxy.ErrSecretChainPledgeSnapshotNotReady {
				respond.BadRequest(c, http.StatusBadRequest, "secret chain snapshot not ready for "+date)
				return
			}
			respond.InternalError(c, err)
			return
		}
		snapshot.SIEVolume = sieSum
		snapshot.GASVolume = gasSum
		err = snapshot.Create()
		if err != nil {
			logger.Warnf("create secret chain snapshot failed: %v", err)
		}
	}

	settlement := open.Settlement{Dat: date}
	if err := settlement.GetByDate(); err != nil {
		logger.Warnf("GetSecretChainPledge get settlement by date failed: %v", err)
	}

	respond.Success(c, GetSecretChainPledgeResp{
		Date:            snapshot.Dat,
		SIEVolume:       snapshot.SIEVolume,
		GASVolume:       snapshot.GASVolume,
		GASRewardVolume: settlement.GasRewardVolume,
	})
}
