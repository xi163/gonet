package blowfish

import (
	"fmt"

	"github.com/cwloo/gonet/utils/crypto/aes/ecb"
	"github.com/cwloo/gonet/utils/crypto/padding"
	"golang.org/x/crypto/blowfish"
)

func ECBEncryptPKCS5(pt, key, IV []byte) []byte {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBEncrypter(block, IV)
	padder := padding.NewPkcs5Padding()
	pt, err = padder.Pad(pt) // pad last block of plaintext if block size less than block cipher size
	if err != nil {
		panic(err.Error())
	}
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct
}

func ECBDecryptPKCS5(ct, key, IV []byte) []byte {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBDecrypter(block, IV)
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
func ECBEncryptPKCS7(pt, key, IV []byte) []byte {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBEncrypter(block, IV)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Pad(pt) // padd last block of plaintext if block size less than block cipher size
	if err != nil {
		panic(err.Error())
	}
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct
}

func ECBDecryptPKCS7(ct, key, IV []byte) []byte {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBDecrypter(block, IV)
	pt := make([]byte, len(ct))
	mode.CryptBlocks(pt, ct)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Unpad(pt) // unpad plaintext after decryption
	if err != nil {
		panic(err.Error())
	}
	return pt
}

// Blowfish encryption with ECB and PKCS5 padding
func ECBTest() {
	pt := []byte("Some plain text")
	key := []byte("a_very_secret_key")

	ct := ECBEncryptPKCS5(pt, key, key)
	fmt.Printf("Ciphertext: %x\n", ct)

	recoveredPt := ECBDecryptPKCS5(ct, key, key)
	fmt.Printf("Recovered plaintext: %s\n", recoveredPt)
}
