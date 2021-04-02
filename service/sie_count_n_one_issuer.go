package service

import (
	"fmt"
	"github.com/shopspring/decimal"
	"report-manager/alg"
	"report-manager/db"
	"report-manager/db/open"
	"report-manager/logger"
	"time"
)

type NOneExtra struct {
	UID    string          `json:"uid"`
	Amount decimal.Decimal `json:"amount"`
}

type SIECountNOneIssuer struct {
	Extras []NOneExtra
}

func (s SIECountNOneIssuer) Type() string {
	return db.SieCountNOneIssuer
}

func (s SIECountNOneIssuer) Prepared(date string) bool {
	return true
}

func (s SIECountNOneIssuer) RawData(date string) ([]SIECountRawData, error) {
	endTime, err := alg.NewShTime(date)
	if err != nil {
		return nil, fmt.Errorf("bad date: %v", err)
	}
	endTime = endTime.Add(CountBoundOffset)
	beginTime := endTime.Add(-24 * time.Hour)

	payment := open.ThirdPayment{CreateTime: beginTime}
	err = payment.GetByCreatedAt()
	if err != nil {
		return nil, err
	}
	begin := payment.ID

	payment.ID = 0
	payment.CreateTime = endTime
	err = payment.GetByCreatedAt()
	if err != nil {
		return nil, err
	}
	end := payment.ID

	logger.Infof("start get n_one reward from %d to %d", begin, end)
	t0 := time.Now()
	rewards, err := open.ThirdPayment{}.GetNOneReward(begin, end)
	if err != nil {
		return nil, fmt.Errorf("get n-one buys failed: %v", err)
	}

	logger.Infof("get n_one reward cost %s", time.Now().Sub(t0))

	result := make([]SIECountRawData, 0, len(rewards))
	for _, reward := range rewards {
		result = append(result, SIECountRawData{
			UID:    reward.UID,
			Token:  reward.Token,
			Amount: reward.Amount,
		})
	}

	for _, extra := range s.Extras {
		result = append(result, SIECountRawData{
			UID:    extra.UID,
			Token:  "SIE",
			Amount: extra.Amount,
		})
	}

	return result, nil
}
