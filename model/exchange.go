package model

import (
	"github.com/shopspring/decimal"
)

const (
	VolumePrecision = 4
	PricePrecision  = 4
)

type OtcDailyTradeReportResp struct {
	MarketID      string          `json:"market_id"`
	BuyerNum      int             `json:"buyer_num"`
	SellerNum     int             `json:"seller_num"`
	BuyAmountSum  decimal.Decimal `json:"buy_amount_sum"`
	SellAmountSum decimal.Decimal `json:"sell_amount_sum"`
}

type CTCTradeSummaryResult struct {
	Market    string          `json:"market"`
	VolumeSum decimal.Decimal `json:"volume_sum"`
	FundsSum  decimal.Decimal `json:"funds_sum"`
	AvgPrice  decimal.Decimal `json:"avg_price"`
	Cnt       int             `json:"cnt"`
	TraderNum int             `json:"trader_num"`
}

type DailyTraderNum struct {
	Date         string `json:"date"`
	TraderNum    int    `json:"trader_num"`
	NewTraderNum int    `json:"new_trader_num"`
}

type UserMetric struct {
	ID        int64 `json:"id"`
	Ts        int64 `json:"ts"`
	OnlineNum int64 `json:"online_num"`
	DailyUv   int64 `json:"daily_uv"`
}

type Ticker struct {
	Buy      string `json:"buy"`      // 最佳买入价
	Sell     string `json:"sell"`     // 最佳卖出价
	Low      string `json:"low"`      // 24小时最低价
	High     string `json:"high"`     // 24小时最高价
	Last     string `json:"last"`     // 上次成交价
	Vol      string `json:"volume"`   // 24小时成交量
	Increase string `json:"increase"` // 24小时涨幅
}

type MarketTicker struct {
	Market    string `json:"market"`    // 交易对ID
	Ticker    Ticker `json:"ticker"`    // 行情信息
	Timestamp int64  `json:"timestamp"` // 时间戳
}
