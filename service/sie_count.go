package service

import (
	"fmt"
	"github.com/shopspring/decimal"
	"report-manager/alg"
	"report-manager/config"
	"report-manager/db"
	"report-manager/logger"
	"time"
)

// 统计时间边界为11:50 p.m.
var CountBoundOffset = -10 * time.Minute

func CountSIESugar() error {
	return CountSIEDefault(SIECountSugar{}, alg.NowDate())
}

func CountSIENOneBuy() error {
	return CountSIEDefault(SIECountNOneBuy{}, alg.NowDate())
}

func CountShopDestroy() error {
	return CountSIEDefault(SIECountShopDestroy{}, alg.NowDate())
}

func CountSIEDefault(sieCount SIECount, date string) error {
	cfg := config.GetServer()
	return CountSIE(sieCount, date, cfg)
}

func CountSIE(sieCount SIECount, date string, cfg config.Server) error {
	return countSIE(sieCount, date, alg.NewStrMapFromSlice(cfg.ExchangeFinaUIDs), alg.NewStrMapFromSlice(cfg.WhiteUIDs))
}

type SIECountRawData struct {
	UID    string
	Token  string
	Amount decimal.Decimal
}

type SIECount interface {
	Type() string
	Prepared(date string) bool
	RawData(date string) ([]SIECountRawData, error)
}

type aggregationVal struct {
	fina  decimal.Decimal
	white decimal.Decimal
	user  decimal.Decimal
}

func countSIE(sieCount SIECount, date string, finas, whites *alg.StrMap) error {
	delay := time.Minute
	maxRetry := int(2 * time.Hour / delay)
	retry := 0
	for !sieCount.Prepared(date) {
		retry++
		if date != alg.NowDate() { // abort
			return fmt.Errorf("request date is not today")
		}
		logger.Infof("CountSIE %s not ready, wait a minute", sieCount.Type())
		if retry >= maxRetry {
			return fmt.Errorf("retry too many times")
		}
		time.Sleep(time.Minute)
	}
	rawData, err := sieCount.RawData(date)
	if err != nil {
		return fmt.Errorf("CountSIE %s get raw data failed: %v", sieCount.Type(), err)
	}
	// aggregation
	result := make(map[string]aggregationVal)

	add := func(typ, token string, val decimal.Decimal) {
		vals := result[token]
		switch typ {
		case "fina":
			vals.fina = vals.fina.Add(val)
		case "white":
			vals.white = vals.white.Add(val)
		case "user":
			vals.user = vals.user.Add(val)
		}
		result[token] = vals
	}

	for _, d := range rawData {
		isUser := true
		if finas.Contain(d.UID) {
			isUser = false
			add("fina", d.Token, d.Amount)
		}
		if whites.Contain(d.UID) {
			isUser = false
			add("white", d.Token, d.Amount)
		}
		if isUser {
			add("user", d.Token, d.Amount)
		}
	}

	// persist
	for k, v := range result {
		sieCount := db.SieCount{
			Dat:         date,
			Token:       k,
			Typ:         sieCount.Type(),
			FinaAmount:  v.fina,
			WhiteAmount: v.white,
			UserAmount:  v.user,
		}
		if err := sieCount.Create(); err != nil {
			logger.Errorf("create sie count %+v failed: %v", sieCount, err)
		}
	}
	return nil
}
