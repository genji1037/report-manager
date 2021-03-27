package open

import (
	"time"
)

type RewardDetail struct {
	ID         uint `gorm:"primary_key"`
	CreateTime time.Time
}

func (RewardDetail) TableName() string {
	return "reward_detail"
}

func (r *RewardDetail) Last() error {
	return gormDb.Model(r).Debug().Last(r).Error
}
