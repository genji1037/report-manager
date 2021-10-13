package proxy

import (
	"errors"
	"io/ioutil"
	"net/http"
	"report-manager/config"
)

var ErrFileNotFound = errors.New("file not found")

func GetSSNSFile(dir, filePrefix string) ([]byte, error) {
	resp, err := http.Get(config.GetServer().Proxy.SSNSFile.BaseURI + "/file?path=" + dir + "/" + filePrefix)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrFileNotFound
	}
	return ioutil.ReadAll(resp.Body)
}
