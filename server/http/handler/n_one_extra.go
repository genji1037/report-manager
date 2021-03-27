package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"report-manager/server/http/respond"
	"report-manager/service"
)

type ReceiveNOneExtraRequest struct {
	Date       string              `json:"date" binding:"required"`
	NOneExtras []service.NOneExtra `json:"n_one_extras"`
}

func ReceiveNOneExtra(c *gin.Context) {
	var req ReceiveNOneExtraRequest
	if err := c.ShouldBind(&req); err != nil {
		respond.Error(c, http.StatusBadRequest, http.StatusBadRequest, err.Error())
		return
	}

	nOneIssuer := service.SIECountNOneIssuer{Extras: req.NOneExtras}
	if err := service.CountSIEDefault(nOneIssuer, req.Date); err != nil {
		respond.Error(c, http.StatusInternalServerError, http.StatusInternalServerError, err.Error())
		return
	}

	respond.Success(c, nil)
}
