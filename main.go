package main

import (
	"github.com/vrecan/death"
	_ "net/http/pprof"
	"report-manager/config"
	"report-manager/db"
	"report-manager/db/defi_fund"
	"report-manager/db/exchange"
	"report-manager/db/open"
	"report-manager/db/types"
	"report-manager/job"
	"report-manager/logger"
	serverHttp "report-manager/server/http"
	"syscall"
)

func main() {

	err := config.LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	env := config.GetServer().Env
	if env == "dev" || env == "test" {
		logger.CreateLoggerOnce(logger.DebugLevel, logger.DebugLevel)
	} else {
		logger.CreateLoggerOnce(logger.InfoLevel, logger.InfoLevel)
	}

	// open exchange database
	serverCfg := config.GetServer()
	err = exchange.Open(types.NewConnection(serverCfg.Proxy.Exchange.Database))
	if err != nil {
		logger.Panicf("Failed to open exchange database, %v", err)
	}
	logger.Infof("exchange db connected")

	//err = radar_otc.Open(types.NewConnection(serverCfg.Proxy.RadarOTC.Database))
	//if err != nil {
	//	logger.Panicf("Failed to open radar otc database, %v", err)
	//}
	//logger.Infof("radar otc db connected")

	err = open.Open(types.NewConnection(serverCfg.Proxy.OpenPlatform.Database))
	if err != nil {
		logger.Panicf("Failed to open open platform database, %v", err)
	}
	logger.Infof("open platform db connected")

	err = defi_fund.Open(types.NewConnection(serverCfg.Proxy.DefiFund.Database))
	if err != nil {
		logger.Panicf("Failed to open defi_fund database, %v", err)
	}
	logger.Infof("defi_fund db connected")

	err = db.Open(types.NewConnection(serverCfg.Database))
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
