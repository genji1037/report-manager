package alg

import "time"

var shLoc *time.Location

func init() {
	shLoc, _ = time.LoadLocation("Asia/Shanghai")
}

func NowDate() string {
	return time.Now().Format("2006-01-02")
}

func NewShTime(date string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", date, shLoc)
}
