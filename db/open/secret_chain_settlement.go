package open

import (
	"github.com/shopspring/decimal"
	"time"
)

type Settlement struct {
	ID              uint `gorm:"primary_key"`
	CreatedAt       time.Time
	Dat             string          `gorm:"uniqueIndex;size:30"`
	BlockNumber     int64           `gorm:"not null"`
	GasRewardVolume decimal.Decimal `gorm:"type:decimal(24,8);not null"`
	State           string          `gorm:"size:20;not null"`
	ErrorReason     string          `gorm:"size:500;not null"`
}

func (s *Settlement) GetByDate() error {
	return gormDb.Table("secret_chain.settlements").Model(s).Where("dat = ?", s.Dat).Last(s).Error
}
