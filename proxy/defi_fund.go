package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"report-manager/config"
)

type DefiFundResp struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type GetPlatformSnapshotResp struct {
	PlatformTotalAmount           decimal.Decimal `json:"platform_total_amount"`
	PlatformIssuedBonusAmount     decimal.Decimal `json:"platform_issued_bonus_amount"`
	PlatformTotalAmountIncr       decimal.Decimal `json:"platform_total_amount_incr"`
	PlatformIssuedBonusAmountIncr decimal.Decimal `json:"platform_issued_bonus_amount_incr"`
}

func GetPlatformSnapshot(date string) (GetPlatformSnapshotResp, error) {
	url := config.GetServer().Proxy.DefiFund.BaseURI + "/platform_snapshot/" + date
	resp, err := httpClient.Get(url)
	if err != nil {
		return GetPlatformSnapshotResp{}, err
	}
	defer resp.Body.Close()
	respBs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GetPlatformSnapshotResp{}, err
	}
	var defiResp DefiFundResp
	err = json.Unmarshal(respBs, &defiResp)
	if err != nil {
		return GetPlatformSnapshotResp{}, err
	}

	if defiResp.Code != http.StatusOK {
		return GetPlatformSnapshotResp{}, fmt.Errorf("[%d] %s", defiResp.Code, defiResp.Msg)
	}

	var getPlatformSnapshotResp GetPlatformSnapshotResp
	err = json.Unmarshal(defiResp.Data, &getPlatformSnapshotResp)

	return getPlatformSnapshotResp, err
}
