package cipherExtend

// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Electronic Code Book (ECB) mode.

// ECB provides confidentiality by assigning a fixed ciphertext block to each
// plaintext block.

// See NIST SP 800-38A, pp 08-09

import . "crypto/cipher"

type ecbBlock struct {
	b         Block
	blockSize int
}

func newECBBlock(b Block) *ecbBlock {
	return &ecbBlock{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

func (x *ecbBlock) BlockSize() int { return x.blockSize }
