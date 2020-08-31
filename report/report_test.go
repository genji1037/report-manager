package report

import (
	"fmt"
	"report-manager/config"
	"report-manager/db/exchange"
	"report-manager/db/types"
	"report-manager/logger"
	"testing"
)

func TestMain(m *testing.M) {
	logger.CreateLoggerOnce(logger.DebugLevel, logger.DebugLevel)

	err := config.LoadConfig("../config.yaml")
	if err != nil {
		panic(err)
	}

	// open exchange database
	serverCfg := config.GetServer()
	err = exchange.Open(types.Connection{
		Host:         serverCfg.Proxy.Exchange.Database.Host,
		User:         serverCfg.Proxy.Exchange.Database.User,
		Password:     serverCfg.Proxy.Exchange.Database.Password,
		Database:     serverCfg.Proxy.Exchange.Database.Database,
		Charset:      serverCfg.Proxy.Exchange.Database.Charset,
		MaxIdleConns: serverCfg.Proxy.Exchange.Database.MaxIdleConns,
		MaxOpenConns: serverCfg.Proxy.Exchange.Database.MaxOpenConns,
	})
	if err != nil {
		logger.Panicf("Failed to open exchange database, %v", err)
	}
	m.Run()
}

func TestExchangeReport(t *testing.T) {
	content, err := MakeExchangeReport()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(content)
}
