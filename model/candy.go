package model

import "github.com/shopspring/decimal"

type CirculateAmount struct {
	Token  string          `json:"token"`
	Amount decimal.Decimal `json:"amount"`
}
