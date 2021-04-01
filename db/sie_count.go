package db

import "github.com/shopspring/decimal"

const (
	SieCountTypeSIEReward = "sugar"
	SieCountTypeWallet    = "wallet"
	SieCountExchange      = "exchange"
	SieCountNOneBuy       = "nOneBuy"
	SieCountNOneIssuer    = "nOneIssuer"
	SieCountShopDestroy   = "shopDestroy"
)

type SieCount struct {
	ID          uint            `gorm:"primary_key"`
	Dat         string          `gorm:"size:16;not null"`
	Token       string          `gorm:"size:16;not null"`
	Typ         string          `gorm:"size:16;not null"`
	FinaAmount  decimal.Decimal `gorm:"not null" sql:"type:decimal(32,16)"`
	WhiteAmount decimal.Decimal `gorm:"not null" sql:"type:decimal(32,16)"`
	UserAmount  decimal.Decimal `gorm:"not null" sql:"type:decimal(32,16)"`
}

func (s *SieCount) Create() error {
	return gormDb.Model(s).Create(s).Error
}
