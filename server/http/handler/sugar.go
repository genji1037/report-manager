package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"report-manager/alg"
	"report-manager/db/open"
	"report-manager/model"
	"report-manager/server/http/respond"
	"time"
)

type GetSugarsReq struct {
	Date string `form:"date"`
	model.PageReq
}

type GetSugarsResp struct {
	Sugars []open.Sugar `json:"sugars"`
}

func GetSugars(c *gin.Context) {
	var req GetSugarsReq
	if err := c.ShouldBind(&req); err != nil {
		respond.BadRequest(c, http.StatusBadRequest, err.Error())
		return
	}

	date := req.Date
	if date != "" { // query by date
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			respond.BadRequest(c, http.StatusBadRequest, fmt.Sprintf("bad date format: %v", err))
			return
		}
		sugar, err := getSugarByDate(date)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				respond.BadRequest(c, http.StatusBadRequest, fmt.Sprintf("sugar record not found at %s", date))
				return
			}
			respond.InternalError(c, err)
			return
		}
		respond.Success(c, GetSugarsResp{Sugars: []open.Sugar{sugar}})
		return
	}

	// page query
	sugars, err := open.Sugar{}.QueryDescPage(req.PageReq)
	if err != nil {
		respond.InternalError(c, err)
		return
	}
	respond.Success(c, GetSugarsResp{Sugars: sugars})
}

func getSugarByDate(date string) (open.Sugar, error) {
	var sugar open.Sugar
	sugar.Dat = date
	err := sugar.GetByDate()
	return sugar, err
}

type GetNOneOutResp struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

func GetNOneOut(c *gin.Context) {
	date := c.Query("date")
	yesterdayDate, err := alg.DateAdd(date, -1)
	if err != nil {
		respond.BadRequest(c, http.StatusBadRequest, err.Error())
		return
	}
	sugar := open.Sugar{Dat: date}
	yesterdaySugar := open.Sugar{Dat: yesterdayDate}

	if err := sugar.GetByDate(); err != nil {
		if err == gorm.ErrRecordNotFound {
			respond.BadRequest(c, http.StatusBadRequest, "record not found at "+sugar.Dat)
			return
		}
		respond.InternalError(c, err)
		return
	}

	err = yesterdaySugar.GetByDate()
	if err != nil && err != gorm.ErrRecordNotFound {
		respond.InternalError(c, err)
		return
	}

	respond.Success(c, GetNOneOutResp{
		Date:   date,
		Amount: yesterdaySugar.AccountOut - sugar.AccountOut,
	})
}
