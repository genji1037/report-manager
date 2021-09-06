package exchange

import (
	"report-manager/alg"
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

func (CTCTrade) SumFrozenAmount(includeUIDs, excludeUIDs []string) ([]model.Frozen, error) {
	result := make([]model.Frozen, 0)
	sql := genSumFrozenAmountSQL(includeUIDs, excludeUIDs, false)
	err := gormDb.Raw(sql).Scan(&result).Error
	return result, err
}

func (CTCTrade) SumFrozenAmountGroupByUID(includeUIDs, excludeUIDs []string) ([]model.UserFrozen, error) {
	result := make([]model.UserFrozen, 0)
	sql := genSumFrozenAmountSQL(includeUIDs, excludeUIDs, true)
	err := gormDb.Raw(sql).Scan(&result).Error
	return result, err
}

func genSumFrozenAmountSQL(includeUIDs, excludeUIDs []string, groupByUID bool) string {
	var inStmt string
	if len(includeUIDs) > 0 {
		inStmt = " AND uid in " + alg.SQLIn(includeUIDs)
	} else if len(excludeUIDs) > 0 {
		inStmt = " AND uid not in " + alg.SQLIn(excludeUIDs)
	}

	var maybeUIDAndComma string
	if groupByUID {
		maybeUIDAndComma = "uid,"
	}

	sql := `
SELECT 
    ` + maybeUIDAndComma + `IF(quote = 'bid', bid, ask) AS token, SUM(` + "`locked`" + `) amount
FROM
    ctc_orders
WHERE
    state = 'wait' AND is_robot = 0 ` + inStmt + `
GROUP BY ` + maybeUIDAndComma + `token`

	return sql
}

func (CTCTrade) SumFrozenAmountByUID() ([]model.UserTokenAmount, error) {
	result := make([]model.UserTokenAmount, 0)
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
