package db

import "github.com/shopspring/decimal"

type SSNSReport struct {
	ID          uint            `gorm:"primary_key"`
	Dat         string          `gorm:"size:16;not null"`
	BonusAmount decimal.Decimal `gorm:"not null" sql:"type:decimal(32,16)"`
	LinkAmount  decimal.Decimal `gorm:"not null" sql:"type:decimal(32,16)"`
}

func (s *SSNSReport) Create() error {
	return gormDb.Create(s).Error
}

func (s *SSNSReport) GetByDate() error {
	return gormDb.Model(s).Where("dat = ?", s.Dat).Last(s).Error
}
