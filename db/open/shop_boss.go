package open

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"time"
)

type ShopBoss struct {
	ID        uint
	UID       string
	CreatedAt time.Time
	Amount    decimal.Decimal
}

func (ShopBoss) TableName() string {
	return "shop_boss"
}

func (s *ShopBoss) GetByCreatedAt() error {
	err := gormDb.Model(s).Where("created_at < ?", s.CreatedAt).Last(s).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func (ShopBoss) Query(begin, end uint) ([]ShopBoss, error) {
	result := make([]ShopBoss, 0)
	err := gormDb.Model(new(ShopBoss)).Where("id > ? and id <= ?", begin, end).Scan(&result).Error
	return result, err
}
