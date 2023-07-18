// Copyright 2016 Andre Burgaud. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Electronic Code Book (ECB) mode.

// Implemented for legacy purpose only. ECB should be avoided
// as a mode of operation. Favor other modes available
// in the Go crypto/cipher package (i.e. CBC, GCM, CFB, OFB or CTR).

// See NIST SP 800-38A, pp 9

// The source code in this file is a modified copy of
// https://golang.org/src/crypto/cipher/cbc.go
// and released under the following
// Go Authors copyright and license:

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found at https://golang.org/LICENSE

// Package ecb implements block cipher mode of encryption ECB (Electronic Code
// Book) functions. This is implemented for legacy purposes only and should not
// be used for any new encryption needs. Use CBC (Cipher Block Chaining) instead.
package ecb

import (
	"crypto/cipher"
)

type ecb struct {
	b         cipher.Block
	blockSize int
	iv        []byte
	tmp       []byte
}

func newECB(b cipher.Block, iv []byte) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
		iv:        dup(iv),
		tmp:       make([]byte, b.BlockSize()),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in elecronic codebook (ECB)
// mode, using the given Block (Cipher).
func NewECBEncrypter(b cipher.Block, iv []byte) cipher.BlockMode {
	if len(iv) != b.BlockSize() {
		panic("ecb.NewECBEncrypter: IV length must equal block size")
	}
	return (*ecbEncrypter)(newECB(b, iv))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	// if len(src) < aes.BlockSize {
	// 	panic("crypto/cipher: input too small")
	// }
	// if len(src)%aes.BlockSize != 0 {
	// 	panic("crypto/cipher: input not full blocks")
	// }
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	if InexactOverlap(dst[:len(src)], src) {
		panic("crypto/cipher: invalid buffer overlap")
	}

	for len(src) > 0 {
		x.b.Encrypt(dst[:x.blockSize], src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}

	// iv := x.iv

	// for len(src) > 0 {
	// 	// Write the xor to dst, then encrypt in place.
	// 	xorBytes(dst[:x.blockSize], src[:x.blockSize], iv)
	// 	x.b.Encrypt(dst[:x.blockSize], dst[:x.blockSize])

	// 	// Move to the next block with this block as the next iv.
	// 	iv = dst[:x.blockSize]
	// 	src = src[x.blockSize:]
	// 	dst = dst[x.blockSize:]
	// }

	// // Save the iv for the next CryptBlocks call.
	// copy(x.iv, iv)
}

func (x *ecbEncrypter) SetIV(iv []byte) {
	if len(iv) != len(x.iv) {
		panic("cipher: incorrect length IV")
	}
	copy(x.iv, iv)
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic codebook (ECB)
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block, iv []byte) cipher.BlockMode {
	if len(iv) != b.BlockSize() {
		panic("ecb.NewECBDecrypter: IV length must equal block size")
	}
	return (*ecbDecrypter)(newECB(b, iv))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	// if len(src) < aes.BlockSize {
	// 	panic("crypto/cipher: input too small")
	// }
	// if len(src)%aes.BlockSize != 0 {
	// 	panic("crypto/cipher: input not full blocks")
	// }
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	if InexactOverlap(dst[:len(src)], src) {
		panic("crypto/cipher: invalid buffer overlap")
	}
	// if len(src) == 0 {
	// 	return
	// }

	for len(src) > 0 {
		x.b.Decrypt(dst[:x.blockSize], src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}

	// For each block, we need to xor the decrypted data with the previous block's ciphertext (the iv).
	// To avoid making a copy each time, we loop over the blocks BACKWARDS.
	// end := len(src)
	// start := end - x.blockSize
	// prev := start - x.blockSize

	// // Copy the last block of ciphertext in preparation as the new iv.
	// copy(x.tmp, src[start:end])

	// // Loop over all but the first block.
	// for start > 0 {
	// 	x.b.Decrypt(dst[start:end], src[start:end])
	// 	xorBytes(dst[start:end], dst[start:end], src[prev:start])

	// 	end = start
	// 	start = prev
	// 	prev -= x.blockSize
	// }

	// // The first block is special because it uses the saved iv.
	// x.b.Decrypt(dst[start:end], src[start:end])
	// xorBytes(dst[start:end], dst[start:end], x.iv)

	// // Set the new iv to the first block we copied earlier.
	// x.iv, x.tmp = x.tmp, x.iv
}

func (x *ecbDecrypter) SetIV(iv []byte) {
	if len(iv) != len(x.iv) {
		panic("cipher: incorrect length IV")
	}
	copy(x.iv, iv)
}
