package model

import "github.com/shopspring/decimal"

type RadarMerchantSummary struct {
	UID            string          `json:"uid"`
	Market         string          `json:"market"`
	SellVolume     decimal.Decimal `json:"sell_volume"`
	BuyVolume      decimal.Decimal `json:"buy_volume"`
	DealTradeCount int             `json:"deal_trade_count"`
}
