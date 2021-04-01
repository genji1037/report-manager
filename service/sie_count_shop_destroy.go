package service

import (
	"fmt"
	"report-manager/alg"
	"report-manager/db"
	"report-manager/db/open"
	"time"
)

type SIECountShopDestroy struct {
}

func (s SIECountShopDestroy) Type() string {
	return db.SieCountShopDestroy
}

func (s SIECountShopDestroy) Prepared(date string) bool {
	return true
}

func (s SIECountShopDestroy) RawData(date string) ([]SIECountRawData, error) {
	endTime, err := alg.NewShTime(date)
	if err != nil {
		return nil, fmt.Errorf("bad date: %v", err)
	}
	endTime = endTime.Add(CountBoundOffset)
	beginTime := endTime.Add(-24 * time.Hour)

	// query shop user
	shopUser := open.ShopUser{CreatedAt: beginTime}
	err = shopUser.GetByCreatedAt()
	if err != nil {
		return nil, err
	}
	userBegin := shopUser.ID

	shopUser.ID = 0
	shopUser.CreatedAt = endTime
	err = shopUser.GetByCreatedAt()
	if err != nil {
		return nil, err
	}
	userEnd := shopUser.ID

	userDestroy, err := open.ShopUser{}.Query(userBegin, userEnd)
	if err != nil {
		return nil, fmt.Errorf("open.ShopUser.Query failed: %v", err)
	}

	// query shop boss
	shopBoss := open.ShopBoss{CreatedAt: beginTime}
	err = shopBoss.GetByCreatedAt()
	if err != nil {
		return nil, err
	}
	bossBegin := shopBoss.ID

	shopBoss.ID = 0
	shopBoss.CreatedAt = endTime
	err = shopBoss.GetByCreatedAt()
	if err != nil {
		return nil, err
	}
	bossEnd := shopBoss.ID

	bossDestroy, err := open.ShopBoss{}.Query(bossBegin, bossEnd)
	if err != nil {
		return nil, fmt.Errorf("open.ShopBoss.Query failed: %v", err)
	}

	result := make([]SIECountRawData, 0, len(userDestroy)+len(bossDestroy))
	for _, ud := range userDestroy {
		result = append(result, SIECountRawData{
			UID:    ud.UID,
			Token:  "SIE",
			Amount: ud.Amount,
		})
	}
	for _, bd := range bossDestroy {
		result = append(result, SIECountRawData{
			UID:    bd.UID,
			Token:  "SIE",
			Amount: bd.Amount,
		})
	}

	return result, nil
}
