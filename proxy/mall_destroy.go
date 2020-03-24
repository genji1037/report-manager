package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"report-manager/config"
	"report-manager/model"
)

func ListFailedOrder() ([]model.FailedOrderResp, error) {
	url := config.GetServer().Proxy.MallDestroy.BaseURI + "/destroyer/orders/failed"
	result, err := mallDestroyGet(url, nil)
	if err != nil {
		return nil, fmt.Errorf("mallDestroyGet failed: %s", err.Error())
	}
	orders := make([]model.FailedOrderResp, 0)
	err = json.Unmarshal(result.Result, &orders)
	if err != nil {
		return nil, fmt.Errorf("unmarshal %s to []model.FailedOrderResp failed: %s", string(result.Result), err.Error())
	}
	return orders, nil
}

func mallDestroyGet(url string, params map[string]string) (rsp *ExchangeResult, err error) {

	first := true
	for k, v := range params {
		prefix := "&"
		if first {
			prefix = "?"
			first = false
		}
		url += prefix + k + "=" + v
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	rsp = new(ExchangeResult)
	body, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, rsp); err != nil {
		return rsp, fmt.Errorf("response: body: %s", string(body))
	}
	return
}
