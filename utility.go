package sunshinemotion

import (
	"crypto/md5"
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
