package exchange

import "report-manager/model"

type OrderTrade struct{}

func (OrderTrade) SumFrozenAmount() ([]model.Frozen, error) {
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
        ` + "`type`" + ` = 2 AND ` + "`status`" + ` IN (2 , 3)
    GROUP BY market_id UNION ALL SELECT 
        market_id, SUM(volume) AS frozen_amount
    FROM
        trades
    WHERE
	` + "`type`" + ` = 2 AND state IN (3 , 4)
    GROUP BY market_id) summary
GROUP BY market_id
`
	err := gormDb.Raw(sql).Scan(&result).Error
	return result, err
}
