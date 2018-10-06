package lib

import (
	"crypto/md5"
	"fmt"
	"math/rand"
)

func randRange(min int, max int) int {
	return min + rand.Int()%(max-min+1)
}

func MD5String(raw string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
}

func PasswordHash(raw string) string {
	return MD5String(raw)
}
