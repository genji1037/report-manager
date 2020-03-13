package proxy

import (
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
	"time"
)

func init() {
	http.DefaultClient.Timeout = time.Minute
}

// GetStringFromType could get string from given type
func getStringFromGivenType(v interface{}) string {
	var str string
	switch v.(type) {
	case uint:
		str = strconv.FormatInt(int64(v.(uint)), 10)
	case int:
		str = strconv.Itoa(v.(int))
	case int64:
		str = strconv.FormatInt(v.(int64), 10)
	case string:
		str, _ = v.(string)
	case decimal.Decimal:
		str = v.(decimal.Decimal).String()
	default:
		str = ""
	}
	return str
}
