package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"report-manager/config"
	"report-manager/model"
	"strconv"
)

func OTCDailyTraderNum(fromDate, toDate string) ([]model.DailyTraderNum, error) {
	rs := make([]model.DailyTraderNum, 0)
	baseURI := config.GetServer().Proxy.Exchange.BaseURI
	params := map[string]string{
		"from_date": fromDate,
		"to_date":   toDate,
	}
	rsp, err := exchangeGet(baseURI+"/api/manager/user/trade/daily", params)
	if err != nil {
		return rs, fmt.Errorf("%s/api/manager/user/trade/daily failed: %s", baseURI, err.Error())
	}
	err = json.Unmarshal(rsp.Result, &rs)
	if err != nil {
		return rs, fmt.Errorf("umarshal %s failed: %s", string(rsp.Result), err.Error())
	}
	return rs, nil
}

func CTCDailyTraderNum(fromDate, toDate string) ([]model.DailyTraderNum, error) {
	rs := make([]model.DailyTraderNum, 0)
	baseURI := config.GetServer().Proxy.Exchange.BaseURI
	params := map[string]string{
		"from_date": fromDate,
		"to_date":   toDate,
	}
	rsp, err := exchangeGet(baseURI+"/api/manager/ctc/user/trade/daily", params)
	if err != nil {
		return rs, fmt.Errorf("%s/api/manager/ctc/user/trade/daily failed: %s", baseURI, err.Error())
	}
	err = json.Unmarshal(rsp.Result, &rs)
	if err != nil {
		return rs, fmt.Errorf("umarshal %s failed: %s", string(rsp.Result), err.Error())
	}
	return rs, nil
}

// /api/manager/metrics
func ExchangeUserMetrics(limit, id, step, fromTs, toTs, op int) ([]model.UserMetric, error) {
	var cap int
	if op == 1 {
		cap = (toTs - fromTs) / step
	} else {
		cap = limit
	}
	rs := make([]model.UserMetric, 0, cap)
	baseURI := config.GetServer().Proxy.Exchange.BaseURI
	params := map[string]string{
		"limit":   strconv.Itoa(limit),
		"id":      strconv.Itoa(id),
		"step":    strconv.Itoa(step),
		"op":      strconv.Itoa(op),
		"from_ts": strconv.Itoa(fromTs),
		"to_ts":   strconv.Itoa(toTs),
	}
	rsp, err := exchangeGet(baseURI+"/api/manager/metrics", params)
	if err != nil {
		return rs, fmt.Errorf("%s/api/manager/metrics failed: %s", baseURI, err.Error())
	}
	err = json.Unmarshal(rsp.Result, &rs)
	if err != nil {
		return rs, fmt.Errorf("umarshal %s failed: %s", string(rsp.Result), err.Error())
	}
	return rs, nil
}

func ExchangeTickers() ([]model.MarketTicker, error) {
	rs := make([]model.MarketTicker, 0)
	baseURI := config.GetServer().Proxy.Exchange.BaseURI
	rsp, err := exchangeGet(baseURI+"/api/ctc/tickers", nil)
	if err != nil {
		return rs, fmt.Errorf("%s/api/ctc/tickers failed: %s", baseURI, err.Error())
	}
	err = json.Unmarshal(rsp.Result, &rs)
	if err != nil {
		return rs, fmt.Errorf("umarshal %s failed: %s", string(rsp.Result), err.Error())
	}
	return rs, nil
}

// Result is receive post func return value
type ExchangeResult struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int             `json:"error_code"`
	Description string          `json:"description"`
}

func exchangeGet(url string, params map[string]string) (rsp *ExchangeResult, err error) {

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
	req.Header.Add("otc-session-id", "isecret")

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

func exchangePost(url string, header map[string]string, m map[string]interface{}) (rsp *ExchangeResult, err error) {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	for k, v := range m {
		vStr := getStringFromGivenType(v)
		err := writer.WriteField(k, vStr)
		if err != nil {
			return nil, err
		}
	}
	writer.Close()
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return
	}
	req.Header.Add("content-type", writer.FormDataContentType())
	req.Header.Add("otc-session-id", "isecret")

	for k, v := range header {
		req.Header.Add(k, v)
	}
	defer req.Body.Close()

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
