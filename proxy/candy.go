package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"report-manager/config"
)

func LatestCirculateAmount() (decimal.Decimal, error) {
	baseURI := config.GetServer().Proxy.Candy.BaseURI
	rsp, err := get(baseURI+"/manager/exchange/currency", nil)
	if err != nil {
		return decimal.Zero, err
	}
	circulateAmount := struct {
		Currency float64 `json:"currency"`
	}{}
	err = json.Unmarshal(rsp.RawMessage, &circulateAmount)
	if err != nil {
		return decimal.Zero, fmt.Errorf("unmarshal %s failed: %s", string(rsp.RawMessage), err.Error())
	}
	return decimal.NewFromFloat(circulateAmount.Currency), nil
}

// Result is receive post func return value
type CandyResult struct {
	json.RawMessage
}

func get(url string, header map[string]string) (*CandyResult, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rsp := new(CandyResult)
	body, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, rsp); err != nil {
		return rsp, fmt.Errorf("response body: %s", string(body))
	}
	return rsp, nil
}
