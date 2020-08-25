package job

import (
	"github.com/robfig/cron"
	"report-manager/logger"
	"report-manager/service"
	"time"
)

func StartCronJob() {
	logger.Infof("[cron] starting")
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logger.Warnf("load location failed")
		loc = time.Local
	}
	c := cron.NewWithLocation(loc)

	err = c.AddFunc("0 2 0 * * *", withErr(service.ExchangeReport))
	if err != nil {
		logger.Warnf("add func failed: %s", err.Error())
	}

	err = c.AddFunc("0 0 21 * * *", withErr(service.RadarOTCReport))
	if err != nil {
		logger.Warnf("add func failed: %s", err.Error())
	}

	err = c.AddFunc("0 3 0 * * *", withErr(service.MallDestroyFailedReport))
	if err != nil {
		logger.Warnf("add func failed: %s", err.Error())
	}

	// 手雷OTC待审核商户提醒
	err = c.AddFunc("@every 30m", withErr(service.RadarOTCWaitingRealNames))
	if err != nil {
		logger.Warnf("add func failed: %s", err.Error())
	}

	logger.Infof("[cron] started")
	c.Run()
}

func withErr(f func() error) func() {
	return func() {
		if err := f(); err != nil {
			logger.Errorf("[cron] failed: %s", err.Error())
		}
	}
}
