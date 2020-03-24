package collector

import (
	"fmt"
	"report-manager/config"
	"report-manager/model"
	"report-manager/proxy"
	"strings"
)

type FailedOrderReport struct {
	Data []model.FailedOrderResp
}

func (f *FailedOrderReport) Collect() error {
	orders, err := proxy.ListFailedOrder()
	if err != nil {
		return fmt.Errorf("proxy.ListFailedOrder failed: %s", err.Error())
	}
	f.Data = orders
	return nil
}

func (f *FailedOrderReport) Render(ori string) string {

	lineTemp := config.GetServer().Template.MallDestroyFailedLine
	lineArr := make([]string, 0, len(f.Data))
	for _, v := range f.Data {
		var errMsg string
		if v.State == "failed" {
			errMsg += " 交易所购买失败"
		}
		if v.DestroyState == "failed" {
			errMsg += " 通知积分系统失败"
		}
		errMsg += fmt.Sprintf("[%s] %s", v.ErrCode, v.ErrDesc)
		lineArr = append(lineArr, render(lineTemp, map[string]string{
			"created_at":    v.CreatedAt.String(),
			"sub_system":    v.SubSystem,
			"order_id":      v.OrderID,
			"origin_token":  v.OriginToken,
			"amount":        v.OriginAmount.String(),
			"destroy_token": v.DestroyToken,
			"open_id":       v.OpenID,
			"shop_open_id":  v.ShopOpenID,
			"error":         errMsg,
		}))
	}

	return render(ori, map[string]string{
		"content": strings.Join(lineArr, ""),
	})
}
