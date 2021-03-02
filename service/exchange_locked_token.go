package service

import (
	"bytes"
	"fmt"
	"report-manager/db/exchange"
	"report-manager/db/file"
	"strconv"
	"time"
)

func PersistsOTCLockedTokens() error {
	userFrozen, err := exchange.OrderTrade{}.SumFrozenAmountByUID()
	if err != nil {
		return err
	}
	path := genCSVFileName("OTCLockedTokens")
	buf := bytes.Buffer{}
	for _, uf := range userFrozen {
		buf.WriteString(fmt.Sprintf("%s,%s,%s\n", uf.UID, uf.Token, uf.Amount.String()))
	}
	return file.DefaultDB.SaveFile(path, buf.Bytes())
}

func PersistsCTCLockedTokens() error {
	userFrozen, err := exchange.CTCTrade{}.SumFrozenAmountByUID()
	if err != nil {
		return err
	}
	path := genCSVFileName("CTCLockedTokens")
	buf := bytes.Buffer{}
	for _, uf := range userFrozen {
		buf.WriteString(fmt.Sprintf("%s,%s,%s\n", uf.UID, uf.Token, uf.Amount.String()))
	}
	return file.DefaultDB.SaveFile(path, buf.Bytes())
}

func genCSVFileName(prefix string) string {
	return prefix + time.Now().Format("20060102") + strconv.Itoa(int(time.Now().UnixNano())) + ".csv"
}
