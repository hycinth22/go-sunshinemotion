package sunshinemotion

import (
	"strconv"
	"time"
)

const exchangeTimePattern = "2006-01-02 15:04:05"

func toRPCTimeStr(t time.Time) string {
	return t.Format(exchangeTimePattern)
}
func fromRPCTimeStr(s string) (time.Time, error) {
	return time.Parse(exchangeTimePattern, s)
}
func toRPCDistanceStr(d float64) string {
	return strconv.FormatFloat(d, 'f', 3, 64)
}
