package cipherExtend

import "bytes"

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := PKCS5PaddingLen(cipherText, blockSize)
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, paddingText...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := PKCS5UnPaddingLen(origData)
	// 去掉最后一个字节 unpadding 次
	return origData[:(length - unpadding)]
}

func PKCS5PaddingLen(cipherText []byte, blockSize int) int {
	return blockSize - len(cipherText)%blockSize
}

func PKCS5UnPaddingLen(origData []byte) int {
	return int(origData[len(origData)-1])
}
