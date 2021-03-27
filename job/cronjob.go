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

	//mustAddFunc(c, "@every 30s", func() { logger.Infof("[cron] still alive") })
	mustAddFunc(c, "0 2 0 * * *", withErr(service.ExchangeReport))
	mustAddFunc(c, "0 0 21 * * *", withErr(service.RadarOTCReport))
	mustAddFunc(c, "0 3 0 * * *", withErr(service.MallDestroyFailedReport))
	// 手雷OTC提醒（待审核商户提醒，失败或待重试的转账）
	mustAddFunc(c, "@every 30m", withErr(service.RadarOTCNotice))
	mustAddFunc(c, "0 0 0 * * *", withErr(service.PersistsOTCLockedTokens))
	mustAddFunc(c, "0 0 0 * * *", withErr(service.PersistsCTCLockedTokens))
	mustAddFunc(c, "0 0 0 * * *", withErr(func() error { return service.ExchangeLockedTokensReport(false) }))
	mustAddFunc(c, "0 10 0 * * *", withErr(service.CountSIESugar))
	mustAddFunc(c, "0 0 0 * * *", withErr(service.CountSIENOneBuy))

	logger.Infof("[cron] started")
	c.Run()
}

func mustAddFunc(c *cron.Cron, spec string, cmd func()) {
	if err := c.AddFunc(spec, cmd); err != nil {
		logger.Warnf("add func failed: %s", err.Error())
	}
}

func withErr(f func() error) func() {
	return func() {
		if err := f(); err != nil {
			logger.Errorf("[cron] failed: %s", err.Error())
		}
	}
}
