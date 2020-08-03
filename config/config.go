package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// config root
type Server struct {
	Host     string   `yaml:"host"` // 主机地址
	Port     int      `yaml:"port"` // 端口号
	Template Template `yaml:"template"`
	Proxy    Proxy    `yaml:"proxy"`
}

type Template struct {
	ExchangeDataReport             Report `yaml:"exchange_data_report"`
	OtcDailyReportLine             string `yaml:"otc_daily_report_line"`
	CtcDailyReportLine             string `yaml:"ctc_daily_report_line"`
	CtcCirculationAmountReportLine string `yaml:"ctc_circulation_amount_report_line"`
	CtcClosingPriceReportLine      string `yaml:"ctc_closing_price_report_line"`
	OTCFrozenAmountLine            string `yaml:"otc_frozen_amount_line"`
	CTCFrozenAmountLine            string `yaml:"ctc_frozen_amount_line"`
	MallDestroyFailedReport        Report `yaml:"mall_destroy_failed_report"`
	MallDestroyFailedLine          string `yaml:"mall_destroy_failed_line"`
}

type Report struct {
	Content     string      `yaml:"content"`
	Destination Destination `yaml:"destination"`
}

type Destination struct {
	Console bool   `yaml:"console"` // if true output to console
	GroupID string `yaml:"group_id"`
}

type Proxy struct {
	Exchange     Exchange     `yaml:"exchange"`
	Candy        Candy        `yaml:"candy"`
	OpenPlatform OpenPlatform `yaml:"open_platform"`
	MallDestroy  MallDestroy  `yaml:"mall_destroy"`
}

type Exchange struct {
	BaseURI  string `yaml:"base_uri"`
	Database MySQL  `yaml:"database"`
}

type Candy struct {
	BaseURI string `yaml:"base_uri"`
}

type OpenPlatform struct {
	BaseURI string `yaml:"base_uri"`
}

type MallDestroy struct {
	BaseURI string `yaml:"base_uri"`
}

type MySQL struct {
	Host         string `yaml:"host"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Database     string `yaml:"database"`
	Charset      string `yaml:"charset"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

var srvManger *Manger

type Manger struct {
	server *Server
}

// GetServer 获取配置信息
func GetServer() Server {
	return *srvManger.server
}

// 加载配置文件
func LoadConfig(path string) error {
	// 读取基本配置
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	serve := Server{}
	err = yaml.Unmarshal(data, &serve)
	if err != nil {
		return err
	}

	err = serve.Validate()
	if err != nil {
		return fmt.Errorf("validate config failed: %s", err.Error())
	}

	serve.Init()

	srvManger = &Manger{
		server: &serve,
	}
	return nil
}

func (s Server) Validate() error {
	return nil
}

func (s *Server) Init() {
}
