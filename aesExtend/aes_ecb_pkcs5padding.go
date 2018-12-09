package aesExtend

import (
	"crypto/aes"
)

// AES-128/ECB/PKCS5Padding
func AES_ECB_PKCS5PaddingEncrypt(raw []byte, key []byte) ([]byte, error) {
	aesCipherBlock, err := aes.NewCipher(key[:16])
	if err != nil {
		return nil, err
	}

	afterPaddingContent := PKCS5Padding([]byte(raw), aesCipherBlock.BlockSize())
	buffer := make([]byte, len(afterPaddingContent))
	NewECBEncrypter(aesCipherBlock).CryptBlocks(buffer, afterPaddingContent)
	return buffer, nil
}

// AES-128/ECB/PKCS5Padding
func AES_ECB_PKCS5PaddingDecrypt(raw []byte, key []byte) ([]byte, error) {
	aesCipherBlock, err := aes.NewCipher(key[:16])
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, len(raw))
	NewECBDecrypter(aesCipherBlock).CryptBlocks(buffer, raw)
	unpaddingContent := PKCS5UnPadding(buffer)
	return unpaddingContent, nil
}
