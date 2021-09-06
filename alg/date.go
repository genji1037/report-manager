package alg

import "time"

var shLoc *time.Location

func init() {
	shLoc, _ = time.LoadLocation("Asia/Shanghai")
}

func NowDate() string {
	return time.Now().In(shLoc).Format("2006-01-02")
}

func YesterdayDate() string {
	return time.Now().Add(-24 * time.Hour).In(shLoc).Format("2006-01-02")
}

func NewShTime(date string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", date, shLoc)
}

func ParseSHDate(date string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", date, shLoc)
}
