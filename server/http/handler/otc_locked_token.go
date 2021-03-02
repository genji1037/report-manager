package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"report-manager/server/http/respond"
	"report-manager/service"
)

func PersistOTCLockedToken(c *gin.Context) {
	err := service.PersistsOTCLockedTokens()
	if err != nil {
		respond.Error(c, http.StatusInternalServerError, http.StatusInternalServerError, err.Error())
		return
	}
	respond.Success(c, nil)
}
