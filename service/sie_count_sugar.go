package service

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"io"
	"os"
	"report-manager/alg"
	"report-manager/db"
	"report-manager/db/open"
	"report-manager/logger"
	"report-manager/proxy"
	"time"
)

type SIECountSugar struct{}

func (s SIECountSugar) Type() string {
	return db.SieCountTypeSIEReward
}

func (s SIECountSugar) Prepared(date string) bool {
	return checkSIESugarDone(date)
}

func (s SIECountSugar) RawData(date string) ([]SIECountRawData, error) {
	reward1, reward2, err := downLoadSIESugarRewardFile(date)
	if err != nil {
		return nil, fmt.Errorf("downLoadSIESugarRewardFile failed: %v", err)
	}
	defer reward1.Close()
	defer reward2.Close()

	// unzip reward file

	rewards1, err := solveSIESugarRewardFile(reward1)
	if err != nil {
		return nil, fmt.Errorf("solveSIESugarRewardFile failed: %v", err)
	}
	rewards2, err := solveSIESugarRewardFile(reward2)
	if err != nil {
		return nil, fmt.Errorf("solveSIESugarRewardFile failed: %v", err)
	}
	rewardAll := make([]SIECountRawData, 0, len(rewards1)+len(rewards2))
	rewardAll = append(rewardAll, rewards1...)
	rewardAll = append(rewardAll, rewards2...)

	return rewardAll, nil
}

// 检查SIE糖果计算是否完成
func checkSIESugarDone(date string) bool {
	var r open.RewardDetail
	if err := r.Last(); err != nil {
		logger.Warnf("checkSIESugarDone get last reward detail failed: %v", err)
		return false
	}
	dateTime, err := alg.NewShTime(date)
	if err != nil {
		logger.Warnf("checkSIESugarDone NewShTime failed: %v", err)
		return false
	}
	// since we write reward record after
	return r.CreateTime.After(dateTime)
}

// 下载SIE糖果奖励文件
func downLoadSIESugarRewardFile(date string) (*os.File, *os.File, error) {
	logger.Infof("sieCount begin get reward file name")
	reward1Name, reward2Name, err := proxy.GetRewardFileName()
	if err != nil {
		return nil, nil, fmt.Errorf("proxy.GetRewardFileName failed: %v", err)
	}
	logger.Infof("sieCount got reward file name %s %s, start download.", reward1Name, reward2Name)
	rd1, err := proxy.DownloadSugarFile(reward1Name)
	if err != nil {
		return nil, nil, fmt.Errorf("DownloadSugarFile %s failed: %v", reward1Name, err)
	}
	logger.Infof("sieCount down load %s ok", reward1Name)
	f1Name := fmt.Sprintf("reward_1_%d.zip", time.Now().UnixNano())
	f1, err := alg.WriteTemp(rd1, f1Name)
	if err != nil {
		return nil, nil, fmt.Errorf("WriteTemp %s failed: %v", f1Name, err)
	}
	logger.Infof("sieCount write tmp file %s ok", f1Name)

	rd2, err := proxy.DownloadSugarFile(reward2Name)
	if err != nil {
		return nil, nil, fmt.Errorf("DownloadSugarFile %s failed: %v", reward2Name, err)
	}
	logger.Infof("sieCount down load %s ok", reward2Name)

	f2Name := fmt.Sprintf("reward_2_%d.zip", time.Now().UnixNano())
	f2, err := alg.WriteTemp(rd2, f2Name)
	if err != nil {
		return nil, nil, fmt.Errorf("WriteTemp %s failed: %v", f1Name, err)
	}
	logger.Infof("sieCount write tmp file %s ok", f2Name)

	return f1, f2, nil
}

func solveSIESugarRewardFile(reader *os.File) ([]SIECountRawData, error) {
	logger.Infof("sieCount start solving reward file")
	defer logger.Infof("sieCount solve reward file ok")
	// unzip first
	rd, err := alg.Unzip(reader)
	defer rd.Close()
	if err != nil {
		return nil, fmt.Errorf("unzip failed: %v", err)
	}

	result := make([]SIECountRawData, 0)
	r := bufio.NewReader(rd)
	for {
		line, _, err := r.ReadLine()
		if err == nil {
			bs := bytes.Split(line, []byte(","))
			if len(bs) == 2 {
				uid := string(bytes.TrimSpace(bs[0]))
				balance, err := decimal.NewFromString(string(bs[1]))
				if err != nil {
					err = errors.Wrap(err, "parse string to float")
					return result, err
				}
				result = append(result, SIECountRawData{
					UID:    uid,
					Token:  "SIE",
					Amount: balance,
				})
			} else {
				return result, err
			}
		} else if err == io.EOF {
			break
		} else {
			err = errors.Wrap(err, "read line")
			return result, err
		}
	}
	return result, nil
}
