package respond

import (
	"encoding/json"
	"net/http"
	"report-manager/logger"
	"runtime"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	jsb, ok := data.([]byte)
	if !ok {
		var err error
		jsb, err = json.Marshal(data)
		if err != nil {
			logger.Warnf("Respond data error, %s", err)
			return
		}
	}
	result := struct {
		OK     bool            `json:"ok"`
		Result json.RawMessage `json:"result"`
	}{OK: true, Result: jsb}
	c.JSONP(http.StatusOK, &result)
	return
}

func InternalError(c *gin.Context, err error) {
	_, file, line, _ := runtime.Caller(1)
	logger.Errorf("%v at %s:%d", err, file, line)
	Error(c, 500, 500, "internal server error")
}

func Error(c *gin.Context, status, code int, description interface{}) {
	c.JSON(status, gin.H{
		"ok":          false,
		"error_code":  code,
		"description": description,
	})
}

func BadRequest(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"ok":          false,
		"error_code":  code,
		"description": msg,
	})
}
