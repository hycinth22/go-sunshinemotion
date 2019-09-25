package ssmt

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Float64Range struct {
	Min float64
	Max float64
}
type IntRange struct {
	Min int
	Max int
}

// in range [Min, Max)
func (r *Float64Range) In(d, epsilon float64) bool {
	return d-r.Min >= -epsilon && d-r.Max < +epsilon
}

func init() {
	// TODO: 寻找更好的随机数种子方法，如/dev/random
	rand.Seed(time.Now().UnixNano())
}

// min and max is both in the range
func randRange(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// [min, max)
func randRangeFloat(min, max float64) float64 {
	return min + (max-min)*rand.Float64()
}

func MD5String(raw string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
}

func PasswordHash(raw string) string {
	return MD5String(raw)
}

func GenerateIMEI() string {
	r1, r2 := 10000+randRange(0, 89999), 1000000+randRange(0, 8999999)
	input := "86" + strconv.Itoa(r1) + strconv.Itoa(r2)
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
	index2 := randRange(0, len(letters)-1)
	return letters[index:index+1] + letters[index2:index2+1] + strconv.Itoa(randRange(1, 999)) + "." + strconv.Itoa(randRange(1, 9))
}

func RandScreen() string {
	screens := []string{"1080x1920", "1080x2028", "1080x2030"}
	index := randRange(0, len(screens)-1)
	return screens[index]
}

// %v the value in a default format, adds field names
func DumpStructValue(data interface{}) string {
	return fmt.Sprintf("%+v", data)
}

// %#v	a Go-syntax representation of the value
func DumpStruct(data interface{}) string {
	return fmt.Sprintf("%#v", data)
}
