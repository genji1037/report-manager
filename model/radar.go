package model

import "github.com/shopspring/decimal"

type RadarMerchantSummary struct {
	UID                string          `json:"uid"`
	Market             string          `json:"market"`
	SellVolume         decimal.Decimal `json:"sell_volume"`
	BuyVolume          decimal.Decimal `json:"buy_volume"`
	SellDealTradeCount int             `json:"sell_deal_trade_count"`
	BuyDealTradeCount  int             `json:"buy_deal_trade_count"`
}
