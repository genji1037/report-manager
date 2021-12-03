package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
)

type RadarOTCResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type GetRetryOrFailedTransferResponse struct {
	Transfers []BlockChainTransfer `json:"transfers"`
}

type BlockChainTransfer struct {
	ID         uint             `json:"id"`
	CreatedAt  int64            `json:"created_at"`
	UpdatedAt  int64            `json:"updated_at"`
	TransferAt int64            `json:"transfer_at"` // 转账时间
	Source     string           `json:"source"`
	OrderID    uint             `json:"order_id"`
	OrderBizID string           `json:"order_biz_id"`
	TxHash     string           `json:"tx_hash"`
	From       string           `json:"from"`
	FromUID    string           `json:"from_uid"`
	To         string           `json:"to"`
	ToUID      string           `json:"to_uid"`
	Token      string           `json:"token"`
	Amount     *decimal.Decimal `json:"amount"`
	State      string           `json:"state"`
	Error      string           `json:"error"`
}

//func GetRetryOrFailedTransfer() ([]BlockChainTransfer, error) {
//	baseURI := config.GetServer().Proxy.RadarOTC.BaseURI
//	data, err := radarOTCGet(baseURI + "/v3/otc/manager/transfer/retryOrFailed")
//	if err != nil {
//		return nil, fmt.Errorf("radar OTC get failed: %v", err)
//	}
//	response := GetRetryOrFailedTransferResponse{}
//	err = json.Unmarshal(data, &response)
//	if err != nil {
//		return nil, err
//	}
//	return response.Transfers, nil
//}

func radarOTCGet(uri string) (json.RawMessage, error) {
	resp, err := httpClient.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("http get failed: %v", err)
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := RadarOTCResponse{}
	err = json.Unmarshal(bs, &response)
	if err != nil {
		return nil, err
	}
	if response.Code >= http.StatusBadRequest {
		return nil, fmt.Errorf("[%d] %s", response.Code, response.Data)
	}
	return response.Data, nil
}
