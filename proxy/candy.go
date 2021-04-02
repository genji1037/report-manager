package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"io/ioutil"
	"net/http"
	"report-manager/config"
)

func LatestCirculateAmount() (decimal.Decimal, error) {
	baseURI := config.GetServer().Proxy.OpenPlatform.BaseURI
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

func GetRewardFileName(date string) (string, string, error) {
	baseURI := config.GetServer().Proxy.Candy.BaseURI
	rsp, err := get(baseURI+"/sugar/reward_file_name?date="+date, nil)
	if err != nil {
		return "", "", err
	}
	result := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Reward1Name string `json:"reward_1_name"`
			Reward2Name string `json:"reward_2_name"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(rsp.RawMessage, &result)
	if err != nil {
		return "", "", fmt.Errorf("unmarshal %s failed: %s", string(rsp.RawMessage), err.Error())
	}
	if result.Code != http.StatusOK {
		return "", "", fmt.Errorf("failed with code %d msg %s", result.Code, result.Msg)
	}
	return result.Data.Reward1Name, result.Data.Reward2Name, nil
}

func DownloadSugarFile(fileName string) (io.ReadCloser, error) {
	baseURI := config.GetServer().Proxy.Candy.BaseURI
	resp, err := http.Get(baseURI + "/sugar/download/" + fileName)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
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
