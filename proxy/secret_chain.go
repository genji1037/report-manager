package proxy

import (
	"errors"
	"io/ioutil"
	"net/http"
	"report-manager/config"
	"report-manager/logger"
)

var ErrSecretChainPledgeSnapshotNotReady = errors.New("secret chain snapshot not ready")

func GetSecretChainPledgeSnapshot(date string) ([]byte, error) {
	url := config.GetServer().Proxy.SecretChain.BaseURI + "/sec_chain/i/pledge/snapshot/" + date
	logger.Debugf("proxy.GetSecretChainPledgeSnapshot Get %s", url)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrSecretChainPledgeSnapshotNotReady
	}
	return ioutil.ReadAll(resp.Body)
}
