package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"./aesExtend"
)

const (
	AESKey   = "loaes2019*(#$cry"
	HashSalt = "androidmu3232chang*^12"
)

// AES-128 Cipher
var aesCipherBlock cipher.Block
var aesEcbBlockMode cipher.BlockMode

func init() {
	var err error
	key := []byte(AESKey)[:16]
	aesCipherBlock, err = aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	aesEcbBlockMode = aesExtend.NewECBEncrypter(aesCipherBlock)
}

// AES/ECB/PKCS5Padding
func aes_ecb_pkcs5padding_encrypt(raw string) string {
	afterPaddingContent := aesExtend.PKCS5Padding([]byte(raw), aesCipherBlock.BlockSize())
	buffer := make([]byte, len(afterPaddingContent))
	aesEcbBlockMode.CryptBlocks(buffer, afterPaddingContent)
	return fmt.Sprintf("%X", buffer)
}