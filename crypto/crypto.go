package crypto

import (
	"bytes"
	"encoding/base64"
	"log"
	"strconv"

	"github.com/inkedawn/go-sunshinemotion/crypto/aesExtend"
)

var (
	AESKey   = []byte("loaes2019*(#$cry")
	HashSalt = "androidmu3232chang*^12"
)

func CalcXTcode(userId uint, beginTime string, distance string) string {
	phrase := strconv.FormatUint(uint64(userId), 10) + beginTime + distance + HashSalt
	key := md5String(phrase)
	var xtCode bytes.Buffer
	xtCode.WriteByte(key[7])
	xtCode.WriteByte(key[3])
	xtCode.WriteByte(key[15])
	xtCode.WriteByte(key[24])
	xtCode.WriteByte(key[9])
	xtCode.WriteByte(key[17])
	xtCode.WriteByte(key[29])
	xtCode.WriteByte(key[23])
	e, err := aesExtend.AesEcbPkcs5PaddingEncrypt(xtCode.Bytes(), AESKey)
	if err != nil {
		panic(err)
	}
	return bytesToHexString(e)
}

func CalcLi(p0, p1 string) string {
	phrase := p0 + p1 + HashSalt
	key := md5String(phrase)
	var li bytes.Buffer
	li.WriteByte(key[7])
	li.WriteByte(key[3])
	li.WriteByte(key[11])
	li.WriteByte(key[20])
	li.WriteByte(key[9])
	li.WriteByte(key[14])
	li.WriteByte(key[29])
	li.WriteByte(key[21])
	e, err := aesExtend.AesEcbPkcs5PaddingEncrypt(li.Bytes(), AESKey)
	if err != nil {
		panic(err)
	}
	return bytesToHexString(e)
}

func EncryptBZ(bz string) string {
	// Base64
	bzBytes := []byte(bz)
	base64Result := make([]byte, base64.StdEncoding.EncodedLen(len(bzBytes)))
	base64.StdEncoding.Encode(base64Result, bzBytes)
	log.Println("bz", bz, "base64:", string(base64Result))

	// AesEcbPkcs5PaddingEncrypt
	e, err := aesExtend.AesEcbPkcs5PaddingEncrypt(base64Result, AESKey)
	if err != nil {
		panic(err)
	}
	return bytesToHexString(e)
}

func DecryptBZ(bz string) (string, error) {
	// AesEcbPkcs5PaddingDecrypt
	raw, err := hexStringToBytes(bz)
	if err != nil {
		return "", err
	}
	decrypted, err := aesExtend.AesEcbPkcs5PaddingDecrypt(raw, AESKey)
	if err != nil {
		return "", err
	}
	// Base64
	base64Bytes := decrypted
	lenBase64 := len(base64Bytes)
	base64decoded := make([]byte, base64.StdEncoding.DecodedLen(lenBase64))
	n, err := base64.StdEncoding.Decode(base64decoded, base64Bytes)
	if err != nil {
		return "", err
	}
	result := base64decoded
	return string(result[:n]), nil
}
