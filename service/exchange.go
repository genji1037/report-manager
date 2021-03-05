package service

import (
	"fmt"
	"report-manager/config"
	"report-manager/proxy"
	"report-manager/report"
)

func ExchangeReport() error {
	// make a report
	reportContent, err := report.MakeExchangeReport()
	if err != nil {
		return fmt.Errorf("make report failed: %s", err.Error())
	}
	if config.GetServer().Template.ExchangeDataReport.Destination.Console {
		fmt.Println("=====console report=====")
		fmt.Println(reportContent)
		fmt.Println("========================")
		return nil
	}
	// send report
	err = proxy.SendMessage(reportContent, config.GetServer().Template.ExchangeDataReport.Destination.GroupID)
	if err != nil {
		return fmt.Errorf("send message failed: %s", err.Error())
	}
	return nil
}

func ExchangeLockedTokensReport(console bool) error {
	// make a report
	reportContent, err := report.MakeExchangeLockedTokensReport()
	if err != nil {
		return fmt.Errorf("make report failed: %s", err.Error())
	}
	if console || config.GetServer().Template.ExchangeLockedTokensReport.Destination.Console {
		fmt.Println("=====console report=====")
		fmt.Println(reportContent)
		fmt.Println("========================")
		return nil
	}
	// send report
	err = proxy.SendMessage(reportContent, config.GetServer().Template.ExchangeLockedTokensReport.Destination.GroupID)
	if err != nil {
		return fmt.Errorf("send message failed: %s", err.Error())
	}
	return nil
}
