package collector

import (
	"fmt"
	"report-manager/config"
	"report-manager/db/exchange"
	"report-manager/model"
	"report-manager/proxy"
	"report-manager/util"
	"strconv"
	"strings"
	"time"
)

/*
1、OTC买单金额多少，交易人数多少(OTCDailyReport)
2、OTC卖单金额多少，交易人数多少(OTCDailyReport)
3、币币交易量多少，交易人数多少(CTCDailyReport)
4、交易所新增用户多少(OTCDailyTraderNum, CTCDailyTraderNum)
5、交易所活跃用户多少(OTCDailyTraderNum, CTCDailyTraderNum)
6、12-24点的平均在线人数(ExchangeUserMetrics)
7、流通量多少
8、收盘价格
*/

type OTCDailyReport struct {
	Date string
	Data []model.OtcDailyTradeReportResp
}

func (o *OTCDailyReport) Collect() error {
	begin, err := time.ParseInLocation("2006-01-02", o.Date, util.ShLoc())
	if err != nil {
		return err
	}
	trade := exchange.Trade{}
	rs, err := trade.DailyTraderNum(begin)
	if err != nil {
		return fmt.Errorf("exchange.Trade.DailyTraderNum(%s) failed: %s", begin, err.Error())
	}

	o.Data = rs
	return nil
}

func (o *OTCDailyReport) Render(ori string) string {

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

type CTCDailyReport struct {
	Date string
	Data []model.CTCTradeSummaryResult
}

func (c *CTCDailyReport) Collect() error {
	begin, err := time.ParseInLocation("2006-01-02", c.Date, util.ShLoc())
	if err != nil {
		return err
	}
	ctcTrade := exchange.CTCTrade{}
	rs, err := ctcTrade.TradeSummaryDaily(begin)
	if err != nil {
		return fmt.Errorf("exchange.CTCTrade.TradeSummaryDaily(%s) failed: %s", begin, err.Error())
	}

	c.Data = rs
	return nil
}

func (c *CTCDailyReport) Render(ori string) string {

	lineTemp := config.GetServer().Template.CtcDailyReportLine
	lineArr := make([]string, 0, len(c.Data))
	for _, v := range c.Data {

		lineArr = append(lineArr, render(lineTemp, map[string]string{
			"market":     v.Market,
			"volume_sum": v.VolumeSum.Truncate(model.VolumePrecision).String(),
			"trader_num": strconv.Itoa(v.TraderNum),
		}))

	}

	return render(ori, map[string]string{
		"ctc_report": strings.Join(lineArr, ""),
	})
}

type OTCDailyTraderNum struct {
	Date string
	Data model.DailyTraderNum
}

func (o *OTCDailyTraderNum) Collect() error {
	rs, err := proxy.OTCDailyTraderNum(o.Date, o.Date)
	if err != nil {
		return fmt.Errorf("proxy.OTCDailyTraderNum(%s, %s) failed: %s", o.Date, o.Date, err.Error())
	}
	if len(rs) > 0 {
		o.Data = rs[0]
	}
	return nil
}

func (o *OTCDailyTraderNum) Render(ori string) string {
	return render(ori, map[string]string{
		"otc_new_trader_num": strconv.Itoa(o.Data.NewTraderNum),
		"otc_trader_num":     strconv.Itoa(o.Data.TraderNum),
	})
}

type CTCDailyTraderNum struct {
	Date string
	Data model.DailyTraderNum
}

func (o *CTCDailyTraderNum) Collect() error {
	rs, err := proxy.CTCDailyTraderNum(o.Date, o.Date)
	if err != nil {
		return fmt.Errorf("proxy.CTCDailyTraderNum(%s, %s) failed: %s", o.Date, o.Date, err.Error())
	}
	if len(rs) > 0 {
		o.Data = rs[0]
	}
	return nil
}

func (o *CTCDailyTraderNum) Render(ori string) string {
	return render(ori, map[string]string{
		"ctc_new_trader_num": strconv.Itoa(o.Data.NewTraderNum),
		"ctc_trader_num":     strconv.Itoa(o.Data.TraderNum),
	})
}

type ExchangeUserMetrics struct {
	FromTs int
	ToTs   int
	Data   []model.UserMetric
}

func (e *ExchangeUserMetrics) Collect() error {
	// fetch metrics between from_ts and to_ts, so we set op = 1.
	// step every 10 minutes is enough, so we set step = 600
	rs, err := proxy.ExchangeUserMetrics(0, 0, 600, e.FromTs, e.ToTs, 1)
	if err != nil {
		return fmt.Errorf("proxy.ExchangeUserMetrics(0, 0, 600, %d, %d, 1) failed: %s", e.FromTs, e.ToTs, err.Error())
	}
	e.Data = rs
	return nil
}

func (e *ExchangeUserMetrics) Render(ori string) string {
	var numerator, denominator int64
	for _, v := range e.Data {
		numerator += v.OnlineNum
		denominator++
	}
	if denominator <= 0 {
		return ori
	}
	avgOnlineNum := numerator / denominator
	return render(ori, map[string]string{
		"half_bottom_avg_online": strconv.Itoa(int(avgOnlineNum)),
	})
}

type CirculateAmount struct {
	// notice: open-platform API only support latest circulate amount,
	// no way to get circulate amount by specific date.
	Data model.CirculateAmount
}

func (c *CirculateAmount) Collect() error {
	c.Data.Token = "SIE"
	circulateAmount, err := proxy.LatestCirculateAmount()
	if err != nil {
		return fmt.Errorf("proxy.LatestCirculateAmount() failed: %s", err.Error())
	}
	c.Data.Amount = circulateAmount
	return nil
}

func (c *CirculateAmount) Render(ori string) string {
	line := config.GetServer().Template.CtcCirculationAmountReportLine
	report := render(line, map[string]string{
		"token":            c.Data.Token,
		"circulate_amount": c.Data.Amount.Truncate(model.VolumePrecision).String(),
	})
	return render(ori, map[string]string{
		"ctc_circulation_amount_report": report,
	})
}

type LatestPrice struct {
	Data []model.MarketTicker
}

func (l *LatestPrice) Collect() error {
	tickers, err := proxy.ExchangeTickers()
	if err != nil {
		return fmt.Errorf("proxy.ExchangeTickers() failed: %s", err.Error())
	}
	l.Data = tickers
	return nil
}

func (l *LatestPrice) Render(ori string) string {

	lineTemp := config.GetServer().Template.CtcClosingPriceReportLine
	lineArr := make([]string, 0, len(l.Data))
	for _, v := range l.Data {
		lineArr = append(lineArr, render(lineTemp, map[string]string{
			"market":        v.Market,
			"closing_price": v.Ticker.Last,
		}))
	}

	return render(ori, map[string]string{
		"ctc_closing_price_report": strings.Join(lineArr, ""),
	})
}
