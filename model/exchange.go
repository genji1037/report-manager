package model

import (
	"github.com/shopspring/decimal"
)

type Quote string

const (
	VolumePrecision       = 4
	PricePrecision        = 4
	Ask             Quote = "ask"
	Bid             Quote = "bid"
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

type Frozen struct {
	Market string          `json:"-"`      // 市场对
	Amount decimal.Decimal `json:"amount"` // 冻结金额
	Token  string          `json:"token"`  // 币种
}

func (Frozen) FromUserFrozen(ufs []UserFrozen) []Frozen {
	m := make(map[string]Frozen)
	for _, uf := range ufs {
		token := uf.Token
		f, ok := m[token]
		if !ok {
			f.Token = token
		}
		f.Amount = f.Amount.Add(uf.Amount)
		m[token] = f
	}
	fs := make([]Frozen, 0)
	for _, v := range m {
		fs = append(fs, v)
	}
	return fs
}

type UserFrozen struct {
	UID string
	Frozen
}

type UserAmount struct {
	UID    string          `json:"uid"`
	Amount decimal.Decimal `json:"amount"` // 金额
}

type UserTokenAmount struct {
	UserAmount
	Token string `json:"token"` // 币种
}

type UserMarketAmount struct {
	UserAmount
	Market string `json:"market"` // 市场对
	Quote  Quote  `json:"quote"`  // 报价方式
}

type ExchangeSpecialUserReport struct {
	Dat           string          `json:"dat"`
	UID           string          `json:"uid"`
	Token         string          `json:"token"`
	OutcomeAmount decimal.Decimal `json:"outcome_amount"` // 支付金额
	IncomeAmount  decimal.Decimal `json:"income_amount"`  // 收入金额
	LockedAmount  decimal.Decimal `json:"locked_amount"`  // 冻结金额
}
