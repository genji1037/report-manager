package open

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"time"
)

type ShopUser struct {
	ID        uint
	UID       string
	CreatedAt time.Time
	Amount    decimal.Decimal
}

func (ShopUser) TableName() string {
	return "shop_user"
}

func (s *ShopUser) GetByCreatedAt() error {
	err := gormDb.Model(s).Where("created_at < ?", s.CreatedAt).Last(s).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func (ShopUser) Query(begin, end uint) ([]ShopUser, error) {
	result := make([]ShopUser, 0)
	err := gormDb.Model(new(ShopUser)).Where("id > ? and id <= ?", begin, end).Scan(&result).Error
	return result, err
}
