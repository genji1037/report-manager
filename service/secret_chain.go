package service

import (
	"bufio"
	"bytes"
	"github.com/shopspring/decimal"
	"io"
	"report-manager/logger"
	"report-manager/model"
	"report-manager/proxy"
)

func GetSecretChainPledge(date string) (decimal.Decimal, decimal.Decimal, error) {
	bs, err := proxy.GetSecretChainPledgeSnapshot(date)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}
	rd := bufio.NewReader(bytes.NewBuffer(bs))
	var sieSum, gasSum decimal.Decimal
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return decimal.Zero, decimal.Zero, err
		}
		var pledge model.Pledge
		err = pledge.Unmarshal(line)
		if err != nil {
			logger.Errorf("unmarshal %s to %T failed: %v", string(line), pledge, err)
			return decimal.Zero, decimal.Zero, err
		}
		sieSum = sieSum.Add(pledge.SIEVolume).Add(pledge.InactiveSIEVolume)
		gasSum = gasSum.Add(pledge.GASVolume).Add(pledge.InactiveGASVolume)
	}
	return sieSum, gasSum, nil
}
