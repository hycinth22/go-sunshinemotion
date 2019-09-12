package ssmt

import (
	"github.com/inkedawn/go-sunshinemotion/v3/crypto"
)

func CalcXTcode(userId int64, beginTime string, distance string) string {
	return crypto.CalcXTcode(userId, beginTime, distance)
}

func EncodeString(r string) string {
	return crypto.EncryptString(r)
}
func DecodeString(r string) (string, error) {
	return crypto.DecryptString(r)
}
