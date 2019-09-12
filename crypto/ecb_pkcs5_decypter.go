package crypto

import (
	. "crypto/cipher"

	"github.com/inkedawn/go-sunshinemotion/v3/crypto/cipherExtend"
)

type ecbPKCS5Decrypter struct {
	ecbDecrypter BlockMode
	blockSize    int
}

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBPKCS5Decrypter(b Block) *ecbPKCS5Decrypter {
	return &ecbPKCS5Decrypter{cipherExtend.NewECBDecrypter(b), b.BlockSize()}
}

func (x *ecbPKCS5Decrypter) CryptBlocks(src []byte) []byte {
	aesBuffer := make([]byte, len(src))
	x.ecbDecrypter.CryptBlocks(aesBuffer, src)
	unpadding := cipherExtend.PKCS5UnPadding(aesBuffer)
	return unpadding
}
