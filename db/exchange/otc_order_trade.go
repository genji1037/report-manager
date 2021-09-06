package exchange

import (
	"report-manager/alg"
	"report-manager/model"
)

type OrderTrade struct{}

func (OrderTrade) SumFrozenAmount(includeUIDs, excludeUIDs []string) ([]model.Frozen, error) {
	var orderInStmt, tradeInStmt string
	if len(includeUIDs) > 0 {
		orderInStmt = " AND open_id in " + alg.SQLIn(includeUIDs)
		tradeInStmt = " AND seller_id in " + alg.SQLIn(includeUIDs)
	} else if len(excludeUIDs) > 0 {
		orderInStmt = " AND open_id not in " + alg.SQLIn(excludeUIDs)
		tradeInStmt = " AND seller_id not in " + alg.SQLIn(excludeUIDs)
	}

	result := make([]model.Frozen, 0)
	sql := `
SELECT 
    market_id as market, sum(frozen_amount) as amount
FROM
    (SELECT 
        market_id, SUM(balance) AS frozen_amount
    FROM
        orders
    WHERE
        ` + "`type`" + ` = 2 AND ` + "`status`" + ` IN (2 , 3) ` + orderInStmt + `
    GROUP BY market_id UNION ALL SELECT 
        market_id, SUM(volume) AS frozen_amount
    FROM
        trades
    WHERE
	` + "`type`" + ` = 2 AND state IN (3 , 4) ` + tradeInStmt + `
    GROUP BY market_id) summary
GROUP BY market_id
`
	err := gormDb.Raw(sql).Scan(&result).Error
	return result, err
}

func (OrderTrade) SumFrozenAmountByUID() ([]model.UserTokenAmount, error) {
	result := make([]model.UserTokenAmount, 0)
	sql := `
SELECT 
    uid, token, sum(frozen_amount) as amount
FROM
    (SELECT 
        substring(market_id, 1,locate('/', market_id)-1) as token, SUM(balance) AS frozen_amount, open_id as uid
    FROM
        orders
    WHERE
        ` + "`type`" + ` = 2 AND ` + "`status`" + ` IN (2 , 3)
    GROUP BY uid, token UNION ALL SELECT 
        substring(market_id, 1,locate('/', market_id)-1) as token, SUM(volume) AS frozen_amount, seller_id as uid
    FROM
        trades
    WHERE
	` + "`type`" + ` = 2 AND state IN (3 , 4)
    GROUP BY uid, token) summary
GROUP BY uid, token
`
	err := gormDb.Raw(sql).Scan(&result).Error
	return result, err
}
