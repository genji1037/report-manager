package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"report-manager/db"
	"report-manager/logger"
	"report-manager/server/http/respond"
)

type WalletBalance struct {
	Type   int             `json:"type"`   // 1: 财务, 2: 用户.
	Token  string          `json:"token"`  // 币种
	Amount decimal.Decimal `json:"amount"` // 金额
}

type WalletBalanceRequest struct {
	Date           string          `json:"date" binding:"required"` // 日期
	WalletBalances []WalletBalance `json:"balances"`
}

func ReceiveWalletBalance(c *gin.Context) {
	var req WalletBalanceRequest
	if err := c.ShouldBind(&req); err != nil {
		respond.Error(c, http.StatusBadRequest, http.StatusBadRequest, err.Error())
		return
	}
	m := make(map[string]db.SieCount)
	for _, data := range req.WalletBalances {
		token := data.Token
		sieCount, ok := m[token]
		if !ok {
			sieCount = db.SieCount{
				Dat:         req.Date,
				Token:       token,
				Typ:         db.SieCountTypeWallet,
				FinaAmount:  decimal.Decimal{},
				WhiteAmount: decimal.Decimal{},
				UserAmount:  decimal.Decimal{},
			}
		}
		switch data.Type {
		case 1:
			sieCount.FinaAmount = sieCount.FinaAmount.Add(data.Amount)
		case 2:
			sieCount.UserAmount = sieCount.UserAmount.Add(data.Amount)
		default:
			continue
		}
		m[token] = sieCount
	}

	go func() {
		for _, v := range m {
			err := v.Create()
			if err != nil {
				logger.Errorf("ReceiveWalletBalance create sie count failed: %v", err)
			}
		}
	}()

	respond.Success(c, nil)
}
