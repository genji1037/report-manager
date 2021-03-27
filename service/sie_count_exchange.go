package service

import (
	"report-manager/db"
)

type SIECountExchange struct {
	SIERawData []SIECountRawData
}

func (s SIECountExchange) Type() string {
	return db.SieCountExchange
}

func (s SIECountExchange) Prepared(date string) bool {
	// we call this func explicitly after data prepare. so always return true.
	return true
}

func (s SIECountExchange) RawData(date string) ([]SIECountRawData, error) {
	return s.SIERawData, nil
}
