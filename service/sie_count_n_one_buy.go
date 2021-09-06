package service

import (
	"fmt"
	"report-manager/alg"
	"report-manager/db"
	"report-manager/db/open"
	"time"
)

type SIECountNOneBuy struct {
}

func (s SIECountNOneBuy) Type() string {
	return db.SieCountNOneBuy
}

func (s SIECountNOneBuy) Prepared(date string) bool {
	return true
}

func (s SIECountNOneBuy) RawData(date string) ([]SIECountRawData, error) {
	endTime, err := alg.NewShTime(date)
	if err != nil {
		return nil, fmt.Errorf("bad date: %v", err)
	}
	endTime = endTime.Add(CountBoundOffset)
	beginTime := endTime.Add(-24 * time.Hour)

	rs, err := open.ThirdPayment{}.CreatedAt2ID([]time.Time{beginTime, endTime})
	if err != nil {
		return nil, err
	}

	begin := rs[0]
	end := rs[1]

	buys, err := open.ThirdPayment{}.GetNOneBuy(begin, end)
	if err != nil {
		return nil, fmt.Errorf("get n-one buys failed: %v", err)
	}

	result := make([]SIECountRawData, 0, len(buys))
	for _, buy := range buys {
		result = append(result, SIECountRawData{
			UID:    buy.UID,
			Token:  buy.Token,
			Amount: buy.Amount,
		})
	}

	return result, nil
}
