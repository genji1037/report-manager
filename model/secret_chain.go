package model

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

type Pledge struct {
	UID               string
	SIEVolume         decimal.Decimal
	GASVolume         decimal.Decimal
	InactiveSIEVolume decimal.Decimal
	InactiveGASVolume decimal.Decimal
	SIERecordVolume   decimal.Decimal
}

func (p Pledge) Marshal() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s",
		p.UID, p.SIEVolume.String(), p.GASVolume.String(), p.InactiveSIEVolume.String(),
		p.InactiveGASVolume.String(), p.SIERecordVolume.String())
}

func (p *Pledge) Unmarshal(bs []byte) error {
	str := string(bs)
	arr := strings.Split(str, ",")
	if len(arr) < 6 {
		return fmt.Errorf("expect at least 6 elements got %d", len(arr))
	}
	var err error
	p.SIEVolume, err = decimal.NewFromString(arr[1])
	if err != nil {
		return fmt.Errorf("parse SIE volume failed: %v", err)
	}
	p.GASVolume, err = decimal.NewFromString(arr[2])
	if err != nil {
		return fmt.Errorf("parse GAS volume failed: %v", err)
	}
	p.InactiveSIEVolume, err = decimal.NewFromString(arr[3])
	if err != nil {
		return fmt.Errorf("parse inactive SIE volume failed: %v", err)
	}
	p.InactiveGASVolume, err = decimal.NewFromString(arr[4])
	if err != nil {
		return fmt.Errorf("parse inactive GAS volume failed: %v", err)
	}
	p.SIERecordVolume, err = decimal.NewFromString(arr[5])
	if err != nil {
		return fmt.Errorf("parse SIE record volume failed: %v", err)
	}
	return nil
}
