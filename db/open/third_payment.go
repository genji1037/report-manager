package open

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"time"
)

type ThirdPayment struct {
	ID         uint
	AppID      string
	UID        string
	CreateTime time.Time
	Token      string
	Amount     decimal.Decimal
}

func (ThirdPayment) TableName() string {
	return "third_payment"
}

func (t *ThirdPayment) GetByCreatedAt() error {
	err := gormDb.Model(t).Where("create_time < ?", t.CreateTime).Last(t).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func (ThirdPayment) GetNOneBuy(begin, end uint) ([]ThirdPayment, error) {
	result := make([]ThirdPayment, 0)
	db := gormDb.Model(new(ThirdPayment))
	err := queryNOne(db, 20, begin, end).Scan(&result).Error
	return result, err
}

func (ThirdPayment) GetNOneReward(begin, end uint) ([]ThirdPayment, error) {
	result := make([]ThirdPayment, 0)
	db := gormDb.Model(new(ThirdPayment))
	err := queryNOne(db, 22, begin, end).Scan(&result).Error
	return result, err
}

func queryNOne(db *gorm.DB, payType int, begin, end uint) *gorm.DB {
	criteria := fmt.Sprintf("id > ? and id <= ? and app_id = 'zd7b3n7nazce89bf' and pay_type = %d and state = 1", payType)
	return db.Where(criteria, begin, end)
}
