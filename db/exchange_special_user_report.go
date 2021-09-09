package db

import (
	"bytes"
	"github.com/prometheus/common/log"
	"github.com/shopspring/decimal"
	"report-manager/model"
	"time"
)

type ExchangeSpecialUserReport struct {
	ID            uint `gorm:"primary_key"`
	CreatedAt     time.Time
	Dat           string          `gorm:"size:10;not null;index:dat_uid"`
	UID           string          `gorm:"size:36;not null;index:dat_uid"`
	Token         string          `gorm:"size:10;not null"`
	OutcomeAmount decimal.Decimal `sql:"type:decimal(32,16)"` // 支付金额
	IncomeAmount  decimal.Decimal `sql:"type:decimal(32,16)"` // 收入金额
	LockedAmount  decimal.Decimal `sql:"type:decimal(32,16)"` // 冻结金额
}

func (r *ExchangeSpecialUserReport) Create() error {
	return gormDb.Create(r).Error
}

func (ExchangeSpecialUserReport) CreateBatch(all []ExchangeSpecialUserReport) error {
	batchCreate := func(subs []ExchangeSpecialUserReport) {
		var vals []interface{}
		sqlbuf := bytes.Buffer{}
		sqlbuf.WriteString("INSERT INTO exchange_special_user_reports (created_at,dat,uid,token,outcome_amount,income_amount,locked_amount) VALUES ")
		for _, d := range subs {
			sqlbuf.WriteString("(?,?,?,?,?,?,?),")
			vals = append(vals, time.Now(), d.Dat, d.UID, d.Token, d.OutcomeAmount, d.IncomeAmount, d.LockedAmount)
		}
		sqlStr := sqlbuf.String()
		// trim the last ,
		sqlStr = sqlStr[0 : len(sqlStr)-1]
		err := gormDb.Exec(sqlStr, vals...).Error
		if err != nil {
			log.Warnf("batch create exchange_special_user_report failed: %v", err)
		}
	}

	totalLen := len(all)
	index := 0
	chunkSize := 500

	for {
		if totalLen-index <= chunkSize { // last batch
			if index >= totalLen {
				break
			}
			sub := all[index:]
			batchCreate(sub)
			break
		}

		next := index + chunkSize
		sub := all[index:next]
		index += chunkSize
		batchCreate(sub)
	}

	return nil
}

func (ExchangeSpecialUserReport) QueryByDatUID(date, uid string) ([]model.ExchangeSpecialUserReport, error) {
	rs := make([]model.ExchangeSpecialUserReport, 0)
	err := gormDb.Model(new(ExchangeSpecialUserReport)).Where("dat = ? and uid = ?", date, uid).Scan(&rs).Error
	return rs, err
}
