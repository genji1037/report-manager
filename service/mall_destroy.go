package service

import (
	"fmt"
	"report-manager/config"
	"report-manager/proxy"
)

func MallDestroyFailedReport() error {
	// make a report
	reportContent, err := MakeMallDestroyFailedListReport()
	if err != nil {
		return fmt.Errorf("make report failed: %s", err.Error())
	}
	if config.GetServer().Template.MallDestroyFailedReport.Destination.Console {
		fmt.Println("=====console report=====")
		fmt.Println(reportContent)
		fmt.Println("========================")
		return nil
	}
	// send report
	err = proxy.SendMessage(reportContent, config.GetServer().Template.MallDestroyFailedReport.Destination.GroupID)
	if err != nil {
		return fmt.Errorf("send message failed: %s", err.Error())
	}
	return nil
}
