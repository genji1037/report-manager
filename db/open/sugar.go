package open

import (
	"report-manager/model"
	"time"
)

type Sugar struct {
	CreateTime        time.Time `json:"-"`
	Sugar             float64   `json:"sugar"`               // 当日糖果发放量
	Currency          float64   `json:"currency"`            // 当日流通量
	RealCurrency      float64   `json:"real_currency"`       // 当日实际流通量
	ShopSIE           float64   `json:"shop_sie"`            // 当前商户SIE
	ShopUsedSIE       float64   `json:"shop_used_sie"`       // 当前商户销毁SIE
	AccountIn         float64   `json:"account_in"`          // 差值账户in
	AccountOut        float64   `json:"account_out"`         // 差值账户out
	AvgGrowthRate     float64   `json:"avg_growth_rate"`     // 用户平均增长率
	Dat               string    `json:"dat"`                 // 日期
	DefiPlegde        float64   `json:"defi_plegde"`         // defi质押
	DefiSaving        float64   `json:"defi_saving"`         // defi定存
	SecretChainPledge float64   `json:"secret_chain_pledge"` // secret chain质押
}

func (s Sugar) TableName() string {
	return "sugars"
}

func (s *Sugar) GetByDate() error {
	return gormDb.Model(s).Where("dat = ?", s.Dat).Last(s).Error
}

func (Sugar) QueryDescPage(req model.PageReq) ([]Sugar, error) {
	offset, limit := req.OffsetLimit()
	sugars := make([]Sugar, 0, limit)
	err := gormDb.Model(new(Sugar)).Order("id desc").Offset(offset).Limit(limit).Find(&sugars).Error
	return sugars, err
}
