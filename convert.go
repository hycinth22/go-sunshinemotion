package lib

import (
	"strconv"
	"time"
)

const exchangeTimePattern = "2006-01-02 15:04:05"

func toExchangeTimeStr(t time.Time) string {
	return t.Format(exchangeTimePattern)
}
func fromExchangeTimeStr(s string) (time.Time, error) {
	return time.Parse(exchangeTimePattern, s)
}

func toExchangeDistanceStr(d float64) string {
	return strconv.FormatFloat(d, 'f', 3, 64)
}
func fromExchangeDistanceStr(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func toExchangeInt64Str(n int64) string {
	return strconv.FormatInt(n, 10)
}