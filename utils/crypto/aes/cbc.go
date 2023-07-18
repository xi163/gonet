package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/cwloo/gonet/utils/crypto/padding"
)

func CBCEncryptPKCS5(pt, key, IV []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := cipher.NewCBCEncrypter(block, IV)
	padder := padding.NewPkcs5Padding()
	pt, err = padder.Pad(pt) // padd last block of plaintext if block size less than block cipher size
	if err != nil {
		panic(err.Error())
	}
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct
}

func CBCDecryptPKCS5(ct, key, IV []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := cipher.NewCBCDecrypter(block, IV)
	pt := make([]byte, len(ct))
	mode.CryptBlocks(pt, ct)
	padder := padding.NewPkcs5Padding()
	pt, err = padder.Unpad(pt) // unpad plaintext after decryption
	if err != nil {
		panic(err.Error())
	}
	return pt
}

// Key size for AES is either: 16 bytes (128 bits), 24 bytes (192 bits) or 32 bytes (256 bits)
func CBCEncryptPKCS7(pt, key, IV []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := cipher.NewCBCEncrypter(block, IV)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Pad(pt) // padd last block of plaintext if block size less than block cipher size
	if err != nil {
		panic(err.Error())
	}
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct
}

func CBCDecryptPKCS7(ct, key, IV []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := cipher.NewCBCDecrypter(block, IV)
	pt := make([]byte, len(ct))
	mode.CryptBlocks(pt, ct)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Unpad(pt) // unpad plaintext after decryption
	if err != nil {
		panic(err.Error())
	}
	return pt
}

// AES encryption with CBC and PKCS7 padding
func CBCTest() {
	pt := []byte("Some plain text")
	// aes_128_cbc
	key := []byte("secretkey16bytes")

	ct := CBCEncryptPKCS7(pt, key, key)
	fmt.Printf("Ciphertext: %x\n", ct)

	recoveredPt := CBCDecryptPKCS7(ct, key, key)
	fmt.Printf("Recovered plaintext: %s\n", recoveredPt)
}
