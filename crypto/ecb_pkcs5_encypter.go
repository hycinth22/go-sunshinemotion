package crypto

import (
	. "crypto/cipher"

	"github.com/inkedawn/go-sunshinemotion/v3/crypto/cipherExtend"
)

type ecbPKCS5Encrypter struct {
	ecbEncrypter BlockMode
	blockSize    int
}

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBPKCS5Encrypter(b Block) *ecbPKCS5Encrypter {
	return &ecbPKCS5Encrypter{cipherExtend.NewECBEncrypter(b), b.BlockSize()}
}

func (x *ecbPKCS5Encrypter) CryptBlocks(src []byte) []byte {
	paddingSrc := cipherExtend.PKCS5Padding(src, x.ecbEncrypter.BlockSize())
	buffer := make([]byte, len(paddingSrc))
	x.ecbEncrypter.CryptBlocks(buffer, paddingSrc)
	return buffer
}
