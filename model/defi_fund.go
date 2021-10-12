package model

import (
	"github.com/shopspring/decimal"
)

type Quota struct {
	MerchantUUID string `json:"merchant_uuid"`
	Dat          string `json:"dat"`
	//Amount       decimal.Decimal `json:"amount"`
	//Balance      decimal.Decimal `json:"balance"`
	Used decimal.Decimal `json:"used"`
}
