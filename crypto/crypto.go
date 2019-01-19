package crypto

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
)

var (
	AESKey   = []byte("loaes2019*(#$cry")
	HashSalt = "androidmu3232chang*^12"

	encrypter *ecbPKCS5Encrypter // AES128-ECB-PKCS5Padding Encrypter
	decrypter *ecbPKCS5Decrypter // AES128-ECB-PKCS5Padding Decrypter
)

func init() {
	aesBlock, err := aes.NewCipher(AESKey)
	if err != nil {
		panic(err)
	}
	encrypter = NewECBPKCS5Encrypter(aesBlock)
	decrypter = NewECBPKCS5Decrypter(aesBlock)
}

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
	// AES128-ECB-PKCS5Padding Encrypt
	result := encrypter.CryptBlocks(xtCode.Bytes())
	return bytesToHexString(result)
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
	// AES128-ECB-PKCS5Padding Encrypt
	result := encrypter.CryptBlocks(li.Bytes())
	return bytesToHexString(result)
}

func EncryptBZ(bz string) string {
	// Base64 Encode
	bzBytes := []byte(bz)
	base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(bzBytes)))
	base64.StdEncoding.Encode(base64Text, bzBytes)
	log.Println("bz", bz, "base64:", string(base64Text))
	// AES128-ECB-PKCS5Padding Encrypt
	result := encrypter.CryptBlocks(base64Text)
	return bytesToHexString(result)
}

func DecryptBZ(bz string) (string, error) {
	// AES128-ECB-PKCS5Padding Decrypt
	origin, err := hexStringToBytes(bz)
	if err != nil {
		return "", err
	}
	base64Bytes := decrypter.CryptBlocks(origin)
	// Base64 Decode
	lenBase64 := len(base64Bytes)
	base64decoded := make([]byte, base64.StdEncoding.DecodedLen(lenBase64))
	n, err := base64.StdEncoding.Decode(base64decoded, base64Bytes)
	if err != nil {
		return "", errors.New("Decode Base64 Data" + fmt.Sprint(string(base64Bytes)) + ": " + err.Error())
	}
	result := base64decoded
	return string(result[:n]), nil
}
