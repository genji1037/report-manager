package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"report-manager/config"
)

func SendMessage(msg, cid string) error {
	url := fmt.Sprintf("%s%s", config.GetServer().Proxy.OpenPlatform.BaseURI, "/manager/im/group/msg/send")
	return alert(msg, url, cid)
}

func alert(msg, url, cid string) error {
	form := make(map[string]interface{})
	form["text"] = msg
	form["cid"] = cid
	resp, err := openPlatformPost(url, form)
	if err != nil {
		return err
	}
	if resp.Code != 200 {
		return fmt.Errorf("call[%s] returned code[%d], msg[%s]", url, resp.Code, resp.Msg)
	}
	return nil
}

// Result is receive post func return value
type OpenPlatformResult struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
	Msg  string                 `json:"msg"`
}

func openPlatformPost(url string, m map[string]interface{}) (rsp *OpenPlatformResult, err error) {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	for k, v := range m {
		vStr := getStringFromGivenType(v)
		writer.WriteField(k, vStr)
	}
	writer.Close()
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return
	}
	req.Header.Add("content-type", writer.FormDataContentType())
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	rsp = new(OpenPlatformResult)
	body, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, rsp); err != nil {
		return rsp, fmt.Errorf("response: body:%s", string(body))
	}
	return
}
