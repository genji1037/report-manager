package db

import "github.com/shopspring/decimal"

type SecretChainPledgeSnapshot struct {
	ID        uint            `gorm:"primary_key"`
	Dat       string          `gorm:"size:10;not null;index"`
	SIEVolume decimal.Decimal `sql:"type:decimal(32,16)"`
	GASVolume decimal.Decimal `sql:"type:decimal(32,16)"`
}

func (s *SecretChainPledgeSnapshot) Create() error {
	return gormDb.Create(s).Error
}

func (s *SecretChainPledgeSnapshot) GetByDate() error {
	return gormDb.Model(s).Where("dat = ?", s.Dat).Last(s).Error
}
