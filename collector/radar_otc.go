package collector

import (
	"fmt"
	"github.com/shopspring/decimal"
	"report-manager/config"
	"report-manager/db/radar_otc"
	"report-manager/model"
	"report-manager/proxy"
	"strconv"
	"strings"
	"time"
)

type RadarOTCDailyReport struct {
	Begin time.Time
	End   time.Time
	Data  []model.OtcDailyTradeReportResp
}

func (o *RadarOTCDailyReport) Collect() error {
	rs, err := radar_otc.Trade{}.TradeSummary(o.Begin, o.End)
	if err != nil {
		return fmt.Errorf("exchange.Trade.DailyTraderNum(%s, %s) failed: %s", o.Begin, o.End, err.Error())
	}
	o.Data = rs
	return nil
}

func (o *RadarOTCDailyReport) Render(ori string) string {

	lineTemp := config.GetServer().Template.OtcDailyReportLine
	lineArr := make([]string, 0, len(o.Data))
	for _, v := range o.Data {

		lineArr = append(lineArr, render(lineTemp, map[string]string{
			"market_id":   v.MarketID,
			"buy_amount":  v.BuyAmountSum.String(),
			"buyer_num":   strconv.Itoa(v.BuyerNum),
			"sell_amount": v.SellAmountSum.String(),
			"seller_num":  strconv.Itoa(v.SellerNum),
		}))

	}

	return render(ori, map[string]string{
		"otc_report": strings.Join(lineArr, ""),
	})
}

type RadarOTCDailyTraderNum struct {
	Begin time.Time
	End   time.Time
	Data  model.DailyTraderNum
}

func (o *RadarOTCDailyTraderNum) Collect() error {
	traderNum, err := radar_otc.Trade{}.TraderNum(o.Begin, o.End)
	if err != nil {
		return fmt.Errorf("radar_otc.Trade{}.TraderNum(%s, %s) failed: %s", o.Begin, o.End, err.Error())
	}
	newTraderNum, err := radar_otc.Trade{}.NewTraderNum(o.Begin, o.End)
	if err != nil {
		return fmt.Errorf("radar_otc.Trade{}.NewTraderNum(%s, %s) failed: %s", o.Begin, o.End, err.Error())
	}
	o.Data.TraderNum = traderNum
	o.Data.NewTraderNum = newTraderNum
	return nil
}

func (o *RadarOTCDailyTraderNum) Render(ori string) string {
	return render(ori, map[string]string{
		"otc_new_trader_num": strconv.Itoa(o.Data.NewTraderNum),
		"otc_trader_num":     strconv.Itoa(o.Data.TraderNum),
	})
}

type RadarOTCFrozenAmount struct {
	Data []model.OTCFrozen
}

func (o *RadarOTCFrozenAmount) Collect() error {
	summarizedByMarket, err := radar_otc.OrderTrade{}.SumFrozenAmount()
	if err != nil {
		return fmt.Errorf("exchange.OrderTrade{}.SumFrozenAmount() failed: %v", err)
	}
	// re-summary by token
	summarizedByToken := make([]model.OTCFrozen, 0, len(summarizedByMarket))
	summarizedByTokenMapper := make(map[string]decimal.Decimal)
	getTokenFromMarket := func(market string) string {
		tmpArr := strings.Split(market, "/")
		if len(tmpArr) < 1 {
			return ""
		}
		// market formatted like BTC/CNY, we expect the base currency, so return index 0 when correctly formatted.
		// or return the whole word in case no '/' found at market string.
		return tmpArr[0]
	}
	for _, entry := range summarizedByMarket {
		token := getTokenFromMarket(entry.Market)
		sum, ok := summarizedByTokenMapper[token]
		if ok {
			summarizedByTokenMapper[token] = sum.Add(entry.Amount)
		} else {
			summarizedByTokenMapper[token] = entry.Amount
		}
	}

	// convert map to array
	for token, sum := range summarizedByTokenMapper {
		summarizedByToken = append(summarizedByToken, model.OTCFrozen{
			Frozen: model.Frozen{
				Token:  token,
				Amount: sum,
			},
		})
	}
	o.Data = summarizedByToken

	return nil
}

func (o RadarOTCFrozenAmount) Render(ori string) string {
	lineTemp := config.GetServer().Template.OTCFrozenAmountLine
	lineArr := make([]string, 0, len(o.Data))
	for _, v := range o.Data {
		lineArr = append(lineArr, render(lineTemp, map[string]string{
			"token":  v.Token,
			"amount": v.Amount.String(),
		}))
	}

	return render(ori, map[string]string{
		"otc_frozen_amount": strings.Join(lineArr, ""),
	})
}

type RadarMerchantSummary struct {
	Begin time.Time
	End   time.Time
	Data  []model.RadarMerchantSummary
}

func (r *RadarMerchantSummary) Collect() error {
	buys, err := radar_otc.Trade{}.MerchantSummaryBuy(r.Begin, r.End)
	if err != nil {
		return fmt.Errorf("radar_otc.Trade{}.MerchantSummaryBuy(%s, %s) failed: %v", r.Begin, r.End, err)
	}
	sells, err := radar_otc.Trade{}.MerchantSummarySell(r.Begin, r.End)
	if err != nil {
		return fmt.Errorf("radar_otc.Trade{}.MerchantSummarySell(%s, %s) failed: %v", r.Begin, r.End, err)
	}
	// merge buys and sells
	getKey := func(sum model.RadarMerchantSummary) string {
		return sum.UID + "," + sum.Market
	}
	results := make([]model.RadarMerchantSummary, 0, len(buys))
	mapper := make(map[string]model.RadarMerchantSummary)
	for _, buy := range buys {
		key := getKey(buy)
		sum, ok := mapper[key]
		if ok {
			sum.BuyVolume = buy.BuyVolume
			sum.BuyDealTradeCount = buy.BuyDealTradeCount
		} else {
			mapper[key] = buy
		}
	}
	for _, sell := range sells {
		key := getKey(sell)
		sum, ok := mapper[key]
		if ok {
			sum.BuyVolume = sell.BuyVolume
			sum.BuyDealTradeCount = sell.BuyDealTradeCount
		} else {
			mapper[key] = sell
		}
	}
	for _, v := range mapper {
		results = append(results, v)
	}

	r.Data = results

	return nil
}

func (r *RadarMerchantSummary) Render(ori string) string {
	lineTemp := config.GetServer().Template.RadarMerchantSummaryLine
	lineArr := make([]string, 0, len(r.Data))
	for _, v := range r.Data {
		lineArr = append(lineArr, render(lineTemp, map[string]string{
			"uid":                   v.UID,
			"market":                v.Market,
			"sell_volume":           v.SellVolume.String(),
			"buy_volume":            v.BuyVolume.String(),
			"sell_deal_trade_count": strconv.Itoa(v.SellDealTradeCount),
			"buy_deal_trade_count":  strconv.Itoa(v.BuyDealTradeCount),
		}))
	}

	return render(ori, map[string]string{
		"radar_merchant_summary_line": strings.Join(lineArr, ""),
	})
}

type RadarWaitingRealNames struct {
	Num int
}

func (r *RadarWaitingRealNames) Collect() error {
	cnt, err := radar_otc.RealName{}.CountWaitingRealNames()
	r.Num = cnt
	return err
}

func (r *RadarWaitingRealNames) Render(ori string) string {
	return render(ori, map[string]string{
		"waiting_real_num": strconv.Itoa(r.Num),
	})
}

func (r *RadarWaitingRealNames) Ignore() bool {
	return r.Num == 0
}

type RadarFailedTransfer struct {
	Num int
}

func (r *RadarFailedTransfer) Collect() error {
	transfers, err := proxy.GetRetryOrFailedTransfer()
	r.Num = len(transfers)
	return err
}

func (r *RadarFailedTransfer) Render(s string) string {
	return render(s, map[string]string{
		"failed_transfer_num": strconv.Itoa(r.Num),
	})
}

func (r *RadarFailedTransfer) Ignore() bool {
	return r.Num == 0
}
