package db

import (
	"fmt"
	"time"
)

const ExchangeSpecialUserRoleFina = "fina"

type ExchangeSpecialUser struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Role      string `gorm:"size:10"`
	UID       string `gorm:"size:36;not null;unique_index"`
	Email     string `gorm:"size:128;not null"`
	Remark    string `gorm:"size:512;not null"`
}

func (a *ExchangeSpecialUser) Create() error {
	return gormDb.Create(a).Error
}

func (a *ExchangeSpecialUser) GetByUID() error {
	return gormDb.Model(a).Where("uid = ?", a.UID).Last(a).Error
}

func (ExchangeSpecialUser) DeleteByUID(uid string) error {
	return gormDb.Delete(new(ExchangeSpecialUser), "uid = ?", uid).Error
}

func (a ExchangeSpecialUser) UpdateByUID() error {
	return gormDb.Model(new(ExchangeSpecialUser)).Where("uid = ?", a.UID).Updates(map[string]interface{}{
		"email":  a.Email,
		"remark": a.Remark,
	}).Error
}

func (a ExchangeSpecialUser) QueryMapsByUIDs(uids []string) (map[string]ExchangeSpecialUser, error) {
	results, err := a.QueryByUIDs(uids)
	if err != nil {
		return nil, err
	}
	m := make(map[string]ExchangeSpecialUser, len(results))
	for _, r := range results {
		m[r.UID] = r
	}
	return m, nil
}

func (a ExchangeSpecialUser) QueryByUIDs(uids []string) ([]ExchangeSpecialUser, error) {
	if len(uids) > 1000 {
		return nil, fmt.Errorf("too many uids")
	}
	rs := make([]ExchangeSpecialUser, 0, len(uids))
	err := gormDb.Model(a).Where("uid in ?", uids).Find(&rs).Error
	return rs, err
}

type SpecialUser struct {
	UID    string `json:"uid"`
	Email  string `json:"email"`
	Remark string `json:"remark"`
}

// List lists all exchange special user.
func (a ExchangeSpecialUser) List(role string) ([]SpecialUser, error) {
	rs := make([]SpecialUser, 0)
	err := gormDb.Model(a).Where("role = ?", role).Scan(&rs).Error
	return rs, err
}
