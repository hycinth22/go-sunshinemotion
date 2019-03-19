package ssmt

import (
	"github.com/inkedawn/go-sunshinemotion/crypto"
)

var (
	AESKey   = []byte("loaes2019*(#$cry")
	HashSalt = "androidmu3232chang*^12"
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
