package report

import (
	"fmt"
	"report-manager/collector"
	"report-manager/config"
	"report-manager/logger"
	"report-manager/util"
	"time"
)

type Report string

// MakeExchangeReport make a exchange report
// TODO implement through maker
func MakeExchangeReport() (string, error) {
	logger.Infof("[report] exchange report begin")
	defer logger.Infof("[report] exchange report done")
	loc := util.ShLoc()
	yesterdayDate := time.Now().Add(-24 * time.Hour).In(loc).Format("2006-01-02")
	yesterdayBeginTime, err := time.ParseInLocation("2006-01-02", yesterdayDate, loc)
	if err != nil {
		return "", fmt.Errorf("parse %s to Time.time failed: %s", yesterdayDate, err.Error())
	}
	yesterdayMiddleTs := yesterdayBeginTime.Unix() + int64(12*time.Hour/time.Second)
	yesterdayEndTs := yesterdayMiddleTs + int64(12*time.Hour/time.Second)
	// collect all data
	collectors := []collector.Collector{
		&collector.OTCDailyReport{
			Date: yesterdayDate,
		},
		&collector.CTCDailyReport{
			Date: yesterdayDate,
		},
		&collector.OTCDailyTraderNum{
			Date: yesterdayDate,
		},
		&collector.CTCDailyTraderNum{
			Date: yesterdayDate,
		},
		&collector.ExchangeUserMetrics{
			FromTs: int(yesterdayMiddleTs),
			ToTs:   int(yesterdayEndTs - 1),
		},
		&collector.CirculateAmount{},
		&collector.LatestPrice{},
		&collector.OTCFrozenAmount{},
		&collector.CTCFrozenAmount{},
	}
	collector.Collect(collectors)

	// render
	collectors = append(collectors, collector.NewStringRender("report_date", yesterdayDate))
	out := config.GetServer().Template.ExchangeDataReport.Content
	for i := range collectors {
		out = collectors[i].Render(out)
	}
	return out, nil
}

// TODO implement through maker
func MakeRadarOTCReport() (string, error) {
	logger.Infof("[report] radar otc report begin")
	defer logger.Infof("[report] radar otc report done")
	loc := util.ShLoc()
	now := time.Now().In(loc)
	end := now
	begin := end.Add(-24 * time.Hour)
	date := now.Format("2006-01-02")
	// collect all data
	collectors := []collector.Collector{
		&collector.RadarOTCDailyReport{
			Begin: begin,
			End:   end,
		},
		&collector.RadarOTCDailyTraderNum{
			Begin: begin,
			End:   end,
		},
		&collector.RadarOTCFrozenAmount{},
		&collector.RadarMerchantSummary{
			Begin: begin,
			End:   end,
		},
	}
	collector.Collect(collectors)

	// render
	collectors = append(collectors, collector.NewStringRender("report_date", date))
	out := config.GetServer().Template.RadarOTCReport.Content
	for i := range collectors {
		out = collectors[i].Render(out)
	}
	return out, nil
}

// TODO implement through maker
func MakeMallDestroyFailedListReport() (string, error) {
	logger.Infof("[report] exchange report begin")
	defer logger.Infof("[report] exchange report done")

	loc := util.ShLoc()
	today := time.Now().In(loc).Format("2006-01-02")

	// collect all data
	collectors := []collector.Collector{
		&collector.FailedOrderReport{},
	}
	collector.Collect(collectors)

	// render
	collectors = append(collectors, collector.NewStringRender("report_date", today))
	out := config.GetServer().Template.MallDestroyFailedReport.Content
	for i := range collectors {
		out = collectors[i].Render(out)
	}
	return out, nil
}

func MakeRadarOTCNotice() (string, error) {
	logger.Infof("[report] radar otc real name notice begin")
	defer logger.Infof("[report] radar otc real name notice report done")
	template := config.GetServer().Template.RadarOTCNotice.Content

	// collect real name
	realNameCollector := &collector.RadarWaitingRealNames{}
	err := realNameCollector.Collect()
	if err != nil {
		logger.Errorf("[report] radar otc real name collect waiting real name num failed: %v", err)
		return template, err
	}

	// collect failed or retry transfer

	// do not notice if no waiting real names
	if realNameCollector.Num <= 0 {
		logger.Infof("[report] radar otc no waiting real names")
		return template, DoNotReport
	}
	return realNameCollector.Render(template), nil
}
