package handler

import (
	"github.com/gin-gonic/gin"
	"report-manager/server/http/respond"
	"report-manager/service"
)

func SIECountSugar(c *gin.Context) {
	err := service.CountSIESugar()
	if err != nil {
		respond.Error(c, 400, 400, err.Error())
		return
	}
	respond.Success(c, nil)
}

func SIECountNOneBuy(c *gin.Context) {
	err := service.CountSIENOneBuy()
	if err != nil {
		respond.Error(c, 400, 400, err.Error())
		return
	}
	respond.Success(c, nil)
}

func SIECountShopDestroy(c *gin.Context) {
	err := service.CountShopDestroy()
	if err != nil {
		respond.Error(c, 400, 400, err.Error())
		return
	}
	respond.Success(c, nil)
}
