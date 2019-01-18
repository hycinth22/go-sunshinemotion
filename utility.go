package sunshinemotion

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// min and max is both in the range
func randRange(min int, max int) int {
	return min + rand.Int()%(max-min+1)
}

// The result is lowercase
func md5String(raw string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
}

func parseHTTPDate(date string) (t time.Time, format string, err error) {
	// try RFC1123 first
	t, err = time.Parse(time.RFC1123, date)
	if err == nil {
		return t, time.RFC1123, nil
	}
	// try RFC822 (updated by RFC1123)
	t, err = time.Parse(time.RFC822, date)
	if err == nil {
		return t, time.RFC822, nil
	}
	// try RFC850 (obsoleted by RFC 1036)
	t, err = time.Parse(time.RFC850, date)
	if err == nil {
		return t, time.RFC850, nil
	}
	// try ANSIC (ANSI C's asctime() format)
	t, err = time.Parse(time.ANSIC, date)
	if err == nil {
		return t, time.ANSIC, nil
	}
	return time.Unix(0,0 ), "", errors.New("unknown HTTP Date Header format")
}