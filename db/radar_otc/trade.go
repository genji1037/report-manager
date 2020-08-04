package radar_otc

import (
	"report-manager/model"
	"time"
)

type Trade struct{}

// begin, end time is at Asia/Shanghai location
func (Trade) TradeSummary(begin, end time.Time) ([]model.OtcDailyTradeReportResp, error) {
	rs := make([]model.OtcDailyTradeReportResp, 0)
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
        radar_otc.trades
    WHERE
        state = 1 and updated_at BETWEEN ? AND ?) t
GROUP BY market_id
`, begin, end).Scan(&rs).Error
	return rs, err
}

func (Trade) TraderNum(begin, end time.Time) (int, error) {
	result := struct {
		Cnt int
	}{}
	sql := `
SELECT 
    COUNT(0) as cnt
FROM
    (SELECT DISTINCT
        buyer_id AS trader_id
    FROM
        radar_otc.trades
    WHERE
        created_at BETWEEN ? AND ? UNION SELECT DISTINCT
        seller_id AS trader_id
    FROM
        radar_otc.trades
    WHERE
        created_at BETWEEN ? AND ?) uids`
	err := gormDb.Raw(sql, begin, end, begin, end).Scan(&result).Error
	return result.Cnt, err
}

func (Trade) NewTraderNum(begin, end time.Time) (int, error) {
	result := struct {
		Cnt int
	}{}
	sql := `
SELECT 
    COUNT(0)
FROM
    (SELECT 
        newuids.trader_id
    FROM
        (SELECT DISTINCT
        buyer_id AS trader_id
    FROM
        radar_otc.trades
    WHERE
        created_at BETWEEN ? AND ? UNION SELECT DISTINCT
        seller_id AS trader_id
    FROM
        radar_otc.trades
    WHERE
        created_at BETWEEN ? AND ?) newuids
    LEFT JOIN (SELECT DISTINCT
        buyer_id AS trader_id
    FROM
        radar_otc.trades
    WHERE
        created_at < ? UNION SELECT DISTINCT
        seller_id AS trader_id
    FROM
        radar_otc.trades
    WHERE
        created_at < ?) olduids ON newuids.trader_id = olduids.trader_id
    WHERE
        olduids.trader_id IS NULL) uids`
	err := gormDb.Raw(sql, begin, end, begin, end, begin, begin).Scan(&result).Error
	return result.Cnt, err
}
