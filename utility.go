package lib

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// min and max is both in the range
func randRange(min int, max int) int {
	return min + rand.Int()%(max-min+1)
}

func MD5String(raw string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
}

func PasswordHash(raw string) string {
	return MD5String(raw)
}

func GenerateIMEI() string {
	r1, r2 := 1000000+randRange(0, 8999999), 1000000+randRange(0, 8999999)
	input := strconv.Itoa(r1) + strconv.Itoa(r2)
	var a, b int32 = 0, 0
	for i, tt := range input {
		if i%2 == 0 {
			a = a + tt
		} else {
			temp := tt * 2
			b = b + temp/10 + temp%10
		}
	}
	last := (a + b) % 10
	if last == 0 {
		last = 0
	} else {
		last = 10 - last
	}
	return input + strconv.Itoa(int(last))
}

func RandModel() string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	index := randRange(0, len(letters)-1)
	return letters[index:index+1] + strconv.Itoa(randRange(1, 999))
}
