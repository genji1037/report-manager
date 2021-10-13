package db

import "github.com/shopspring/decimal"

type SSNSReport struct {
	ID          uint   `gorm:"primary_key"`
	Dat         string `gorm:"size:16;not null;index"`
	Seq         int
	BonusAmount decimal.Decimal `gorm:"not null" sql:"type:decimal(32,16)"`
	LinkAmount  decimal.Decimal `gorm:"not null" sql:"type:decimal(32,16)"`
}

func (s *SSNSReport) Create() error {
	return gormDb.Create(s).Error
}

func (s *SSNSReport) GetByDateSeq() error {
	return gormDb.Model(s).Where("dat = ? and seq = ?", s.Dat, s.Seq).Last(s).Error
}
