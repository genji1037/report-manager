package util

import (
	"report-manager/logger"
	"time"
)

func ShLoc() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logger.Errorf("[util] load Asia/Shanghai location failed: %s", err.Error())
		return time.Local
	}
	return loc
}
