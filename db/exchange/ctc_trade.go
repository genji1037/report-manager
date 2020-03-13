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
