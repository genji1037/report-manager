package collector

import (
	"fmt"
	"github.com/shopspring/decimal"
	"report-manager/config"
	"report-manager/db/radar_otc"
	"report-manager/model"
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
	Data  []RadarMerchantSummary
}

func (r *RadarMerchantSummary) Collect() error {
	// TODO
	return nil
}

func (r *RadarMerchantSummary) Render(ori string) string {
	return ""
}
