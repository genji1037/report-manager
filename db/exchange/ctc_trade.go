package exchange

import (
	"report-manager/model"
	"time"
)

type CTCTrade struct {
}

func (c *CTCTrade) TradeSummaryDaily(begin time.Time) ([]model.CTCTradeSummaryResult, error) {
	rs := make([]model.CTCTradeSummaryResult, 0)
	end := begin.Add(time.Hour * 24).Add(-1 * time.Second)
	err := gormDb.Raw(`
select
	market,
	sum(volume) volume_sum,
	sum(funds)/sum(volume) avg_price,
	sum(funds) funds_sum,
	count(*) cnt,
	COUNT(DISTINCT ask_uid)+COUNT(DISTINCT bid_uid) trader_num
from
	ctc_trades
where
	created_at >= ?
	and created_at < ?
    and is_robot = 0
group by
	market
`, begin, end).Scan(&rs).Error

	return rs, err
}

func (CTCTrade) SumFrozenAmount() ([]model.CTCFrozen, error) {
	result := make([]model.CTCFrozen, 0)
	sql := `
SELECT 
    IF(quote = 'bid', bid, ask) AS token, SUM(` + "`locked`" + `) amount
FROM
    ctc_orders
WHERE
    state = 'wait' AND is_robot = 0
GROUP BY token`
	err := gormDb.Raw(sql).Scan(&result).Error
	return result, err
}

func (CTCTrade) SumFrozenAmountByUID() ([]model.UserFrozen, error) {
	result := make([]model.UserFrozen, 0)
	sql := `
SELECT 
    IF(quote = 'bid', bid, ask) AS token, SUM(` + "`locked`" + `) amount, uid
FROM
    ctc_orders
WHERE
    state = 'wait' AND is_robot = 0
GROUP BY token, uid`
	err := gormDb.Raw(sql).Scan(&result).Error
	return result, err
}
