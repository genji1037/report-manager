package alg

import "time"

var shLoc *time.Location

func init() {
	shLoc, _ = time.LoadLocation("Asia/Shanghai")
}

func NowDate() string {
	return time.Now().In(shLoc).Format("2006-01-02")
}

func DateAdd(date string, days int) (string, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", err
	}
	return t.Add(time.Duration(days) * time.Hour).Format("2006-01-02"), nil
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
