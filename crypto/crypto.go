package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"strconv"
)

var (
	AESKey   = []byte("loaes2019*(#$cry")
	HashSalt = "androidmu3232chang*^12"
	DBHashSalt = "locat2018*(#$crymukai"
	aesBlock cipher.Block
)

func init() {
	var err error
	aesBlock, err = aes.NewCipher(AESKey)
	if err != nil {
		panic(err)
	}
}

func CalcDBXTcode(userId int64, beginTime string, distance string) string {
	encrypter := NewECBPKCS5Encrypter(aesBlock)
	phrase := strconv.FormatInt(userId, 10) + beginTime + distance + DBHashSalt
	key := md5String(phrase)
	var xtCode bytes.Buffer
	xtCode.WriteByte(key[4])
	xtCode.WriteByte(key[9])
	xtCode.WriteByte(key[13])
	xtCode.WriteByte(key[26])
	xtCode.WriteByte(key[7])
	xtCode.WriteByte(key[18])
	xtCode.WriteByte(key[30])
	xtCode.WriteByte(key[21])
	// AES128-ECB-PKCS5Padding Encrypt
	result := encrypter.CryptBlocks(xtCode.Bytes())
	return bytesToHexString(result)
}

func CalcXTcode(userId int64, beginTime string, distance string) string {
	encrypter := NewECBPKCS5Encrypter(aesBlock)
	phrase := strconv.FormatInt(userId, 10) + beginTime + distance + HashSalt
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
	encrypter := NewECBPKCS5Encrypter(aesBlock)
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

// AES128-ECB-PKCS5Padding Algorithm
func EncryptString(raw string) string {
	encrypter := NewECBPKCS5Encrypter(aesBlock)
	return bytesToHexString(encrypter.CryptBlocks([]byte(raw)))
}

// AES128-ECB-PKCS5Padding Algorithm
func DecryptString(raw string) (string, error) {
	decrypter := NewECBPKCS5Decrypter(aesBlock)
	data, err := hexStringToBytes(raw)
	if err != nil {
		return "", errors.New("Parsing Encrypted String in HEX format Failed" + err.Error())
	}
	return string(decrypter.CryptBlocks(data)), nil
}
