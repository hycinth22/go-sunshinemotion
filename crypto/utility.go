package crypto

import (
	"crypto/md5"
	"fmt"
)

func bytesToHexString(bytes []byte) string {
	return fmt.Sprintf("%X", bytes)
}

func hexStringToBytes(hex string) (bytes []byte, err error) {
	_, err = fmt.Sscanf(hex, "%X", &bytes)
	return
}

// The result is lowercase
func md5String(raw string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
}
