package service

import (
	"fmt"
	"report-manager/config"
	"report-manager/proxy"
	"report-manager/report"
)

func RadarOTCReport() error {
	// make a report
	reportContent, err := report.RadarOTCReport()
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

func RadarOTCWaitingRealNames() error {
	reportContent, err := report.RadarOTCWaitingRealNames()
	if err != nil {
		if err == report.DoNotReport {
			return nil
		}
		return fmt.Errorf("make report failed: %s", err.Error())
	}
	if config.GetServer().Template.RadarOTCWaitingRealNames.Destination.Console {
		fmt.Println("=====console report=====")
		fmt.Println(reportContent)
		fmt.Println("========================")
		return nil
	}
	// send report
	err = proxy.SendMessage(reportContent, config.GetServer().Template.RadarOTCWaitingRealNames.Destination.GroupID)
	if err != nil {
		return fmt.Errorf("send message failed: %s", err.Error())
	}
	return nil
}
