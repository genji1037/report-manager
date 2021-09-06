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

func (ThirdPayment) CreatedAt2ID(createdAts []time.Time) ([]uint, error) {
	rs := make([]uint, len(createdAts))
	for i, createdAt := range createdAts {
		var t ThirdPayment
		t.CreateTime = createdAt
		err := t.GetByCreatedAt()
		if err != nil {
			return nil, err
		}
		rs[i] = t.ID
	}
	return rs, nil
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

type Summary struct {
	UID     string
	PayType int
	Token   string
	Amount  decimal.Decimal
}

func (s Summary) IsIncome() bool {
	return s.PayType == 21 || s.PayType == 22
}

func (s Summary) IsOutcome() bool {
	return s.PayType == 20
}

func (ThirdPayment) Summary(begin, end time.Time, appID string, UIDs []string) ([]Summary, error) {
	rs, err := ThirdPayment{}.CreatedAt2ID([]time.Time{begin, end})
	if err != nil {
		return nil, err
	}
	return ThirdPayment{}.SummaryByID(rs[0], rs[1], appID, UIDs)
}

func (ThirdPayment) SummaryByID(begin, end uint, appID string, UIDs []string) ([]Summary, error) {
	result := make([]Summary, 0, 3) // 3 pay types
	err := gormDb.Raw(`
SELECT
	uid, pay_type, token, sum(amount) as amount
FROM third_payment
WHERE
	id > ? AND id <= ?
AND app_id = ?
AND state = 1
AND uid in (?)
GROUP BY uid, pay_type, token 
`, begin, end, appID, UIDs).Scan(&result).Error
	return result, err
}
