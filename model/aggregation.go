package model

import "github.com/shopspring/decimal"

type AggregationResult struct {
	FinanceAmount decimal.Decimal
	WhiteAmount   decimal.Decimal
	UserAmount    decimal.Decimal
}
