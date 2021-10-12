package defi_fund

import (
	"github.com/shopspring/decimal"
	"report-manager/model"
	"time"
)

type FundQuota struct {
	ID           uint `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	QuotaID      string `gorm:"size:50;not null;index"`
	Deadline     time.Time
	ExpiredAt    *time.Time      `gorm:"index"`
	MerchantUUID string          `gorm:"not null;index:idx_dat_meruid,unique"` // list quotas can use this index
	Dat          string          `gorm:"not null;index:idx_dat_meruid,unique"`
	Amount       decimal.Decimal `gorm:"type:decimal(24,8);not null"`
	Balance      decimal.Decimal `gorm:"type:decimal(24,8);not null"`
}

func (FundQuota) QueryByDate(date string) ([]FundQuota, error) {
	quotas := make([]FundQuota, 0)
	err := gormDb.Model(new(FundQuota)).Where("dat = ?", date).Scan(&quotas).Error
	return quotas, err
}

func (FundQuota) ToEntities(fq []FundQuota) []model.Quota {
	result := make([]model.Quota, 0, len(fq))
	for i := range fq {
		result = append(result, fq[i].ToEntity())
	}
	return result
}

func (f FundQuota) ToEntity() model.Quota {
	return model.Quota{
		MerchantUUID: f.MerchantUUID,
		Dat:          f.Dat,
		Used:         f.Amount.Sub(f.Balance),
	}
}
