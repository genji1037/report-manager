package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"report-manager/db"
	"report-manager/server/http/respond"
	"time"
)

func GetSSNSReport(c *gin.Context) {
	date := c.Query("date")
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		respond.BadRequest(c, http.StatusBadRequest, err.Error())
		return
	}
	r := db.SSNSReport{Dat: date}
	r.GetByDate()
}
