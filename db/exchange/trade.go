package exchange

import (
	"report-manager/model"
	"time"
)

type Trade struct {
}

const MaxRangeNum = 31

// begin, end time is at Asia/Shanghai location
func (t *Trade) DailyTraderNum(begin time.Time) ([]model.OtcDailyTradeReportResp, error) {
	end := begin.Add(time.Hour * 24).Add(-1 * time.Second)
	rs := make([]model.OtcDailyTradeReportResp, 0, MaxRangeNum)
	err := gormDb.Raw(`
SELECT 
    market_id,
	COUNT(DISTINCT buyer_id) buyer_num,
	COUNT(DISTINCT seller_id) seller_num,
	sum(buy_amount) buy_amount_sum,
	sum(sell_amount) sell_amount_sum
FROM
    (SELECT 
        market_id,
		buyer_id,
		seller_id,
		if(`+"`type`"+`= 1, amount, 0) buy_amount,
		if(`+"`type`"+`= 1, 0, amount) sell_amount
    FROM
        otctrade.trades
    WHERE
        succ_time BETWEEN ? AND ?) t
GROUP BY market_id
`, begin, end).Scan(&rs).Error
	return rs, err
}
