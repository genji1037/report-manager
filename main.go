package main

import (
	"github.com/vrecan/death"
	_ "net/http/pprof"
	"report-manager/config"
	"report-manager/db"
	"report-manager/db/exchange"
	"report-manager/db/open"
	"report-manager/db/radar_otc"
	"report-manager/db/types"
	"report-manager/job"
	"report-manager/logger"
	serverHttp "report-manager/server/http"
	"syscall"
)

func main() {
	logger.CreateLoggerOnce(logger.InfoLevel, logger.DebugLevel)

	err := config.LoadConfig("config.yaml")
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
	logger.Infof("exchange db connected")

	err = radar_otc.Open(types.Connection{
		Host:         serverCfg.Proxy.RadarOTC.Database.Host,
		User:         serverCfg.Proxy.RadarOTC.Database.User,
		Password:     serverCfg.Proxy.RadarOTC.Database.Password,
		Database:     serverCfg.Proxy.RadarOTC.Database.Database,
		Charset:      serverCfg.Proxy.RadarOTC.Database.Charset,
		MaxIdleConns: serverCfg.Proxy.RadarOTC.Database.MaxIdleConns,
		MaxOpenConns: serverCfg.Proxy.RadarOTC.Database.MaxOpenConns,
	})
	if err != nil {
		logger.Panicf("Failed to open radar otc database, %v", err)
	}
	logger.Infof("radar otc db connected")

	err = open.Open(types.Connection{
		Host:         serverCfg.Proxy.OpenPlatform.Database.Host,
		User:         serverCfg.Proxy.OpenPlatform.Database.User,
		Password:     serverCfg.Proxy.OpenPlatform.Database.Password,
		Database:     serverCfg.Proxy.OpenPlatform.Database.Database,
		Charset:      serverCfg.Proxy.OpenPlatform.Database.Charset,
		MaxIdleConns: serverCfg.Proxy.OpenPlatform.Database.MaxIdleConns,
		MaxOpenConns: serverCfg.Proxy.OpenPlatform.Database.MaxOpenConns,
	})
	if err != nil {
		logger.Panicf("Failed to open open platform database, %v", err)
	}
	logger.Infof("open platform db connected")

	err = db.Open(types.Connection{
		Host:         serverCfg.Database.Host,
		User:         serverCfg.Database.User,
		Password:     serverCfg.Database.Password,
		Database:     serverCfg.Database.Database,
		Charset:      serverCfg.Database.Charset,
		MaxIdleConns: serverCfg.Database.MaxIdleConns,
		MaxOpenConns: serverCfg.Database.MaxOpenConns,
	})
	if err != nil {
		logger.Panicf("Failed to open database, %v", err)
	}
	logger.Infof("db connected")

	// 开启http服务
	go serverHttp.Run(serverCfg.Host, serverCfg.Port)

	go job.StartCronJob()

	// 捕捉退出信号
	d := death.NewDeath(syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL,
		syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGALRM)
	d.WaitForDeathWithFunc(func() {
		logger.Infof("report-manager server stopped.")
	})
}
