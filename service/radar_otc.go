package service

import (
	"fmt"
	"report-manager/config"
	"report-manager/proxy"
	"report-manager/util"
	"time"
)

func RadarOTCReport() error {
	// make a report
	reportContent, err := MakeRadarOTCReport()
	if err != nil {
		return fmt.Errorf("make report failed: %s", err.Error())
	}
	if config.GetServer().Template.RadarOTCReport.Destination.Console {
		fmt.Println("=====console report=====")
		fmt.Println(reportContent)
		fmt.Println("========================")
		return nil
	}
	// send report
	err = proxy.SendMessage(reportContent, config.GetServer().Template.RadarOTCReport.Destination.GroupID)
	if err != nil {
		return fmt.Errorf("send message failed: %s", err.Error())
	}
	return nil
}

//func RadarOTCNotice() error {
//	if restTime() {
//		return nil
//	}
//
//	reportContent, err := makeRadarOTCNotice()
//	if err != nil {
//		if err == report.DoNotReport {
//			return nil
//		}
//		return fmt.Errorf("make report failed: %s", err.Error())
//	}
//	if config.GetServer().Template.RadarOTCNotice.Destination.Console {
//		fmt.Println("=====console report=====")
//		fmt.Println(reportContent)
//		fmt.Println("========================")
//		return nil
//	}
//	// send report
//	err = proxy.SendMessage(reportContent, config.GetServer().Template.RadarOTCNotice.Destination.GroupID)
//	if err != nil {
//		return fmt.Errorf("send message failed: %s", err.Error())
//	}
//	return nil
//}

//func makeRadarOTCNotice() (string, error) {
//	reportMaker := report.NewMaker("radar notice", config.GetServer().Template.RadarOTCNotice.Content, []collector.Collector{
//		&collector.RadarWaitingRealNames{},
//		&collector.RadarFailedTransfer{},
//	})
//
//	return reportMaker.Make()
//}

func restTime() bool {
	now := time.Now().In(util.ShLoc())
	return now.Hour() < 9
}
