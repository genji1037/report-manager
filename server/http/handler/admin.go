package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"report-manager/db"
	"report-manager/server/http/respond"
)

type GetSIECountRequest struct {
	Date  string `form:"date"` // e.g. 2021-03-31
	Token string `form:"token"`
}

type GetSIECountResponse struct {
	Entries []SIECountEntry `json:"entries"`
}

type SIECountEntry struct {
	Date       string          `json:"date"`
	Token      string          `json:"token"`
	Type       string          `json:"type"`
	FinaAmount decimal.Decimal `json:"fina_amount"`
	UserAmount decimal.Decimal `json:"user_amount"`
}

func GetSIECount(c *gin.Context) {
	var req GetSIECountRequest
	if err := c.ShouldBind(&req); err != nil {
		respond.Error(c, http.StatusBadRequest, http.StatusBadRequest, err.Error())
		return
	}
	if req.Date == "" && req.Token == "" {
		respond.Error(c, http.StatusBadRequest, http.StatusBadRequest, "require date or token")
		return
	}
	criteria := db.SieCount{
		Dat:   req.Date,
		Token: req.Token,
	}
	result, err := criteria.Query()
	if err != nil {
		respond.Error(c, http.StatusInternalServerError, http.StatusInternalServerError, err.Error())
		return
	}
	entries := make([]SIECountEntry, 0, len(result))
	for i := range result {
		entries = append(entries, SIECountEntry{
			Date:       result[i].Dat,
			Token:      result[i].Token,
			Type:       result[i].Typ,
			FinaAmount: result[i].FinaAmount,
			UserAmount: result[i].UserAmount,
		})
	}

	respond.Success(c, GetSIECountResponse{
		Entries: entries,
	})
}
