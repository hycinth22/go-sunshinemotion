package ssmt

import (
	"strconv"
	"time"
)

const exchangeTimePattern = "2006-01-02 15:04:05"

var (
	TimeZoneCST = time.FixedZone("CST", int((8 * time.Hour).Seconds()))
)

func toServiceStdTime(t time.Time) string {
	return t.In(TimeZoneCST).Format(exchangeTimePattern)
}
func fromServiceStdTime(s string) (time.Time, error) {
	return time.ParseInLocation(exchangeTimePattern, s, TimeZoneCST)
}

func toServiceStdDistance(d float64) string {
	return strconv.FormatFloat(d, 'f', 3, 64)
}
func fromServiceStdDistance(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func toServiceStdInt64(n int64) string {
	return strconv.FormatInt(n, 10)
}
