package test_crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"

	"github.com/cwloo/gonet/logs"
	BASE64 "github.com/cwloo/gonet/utils/codec/base64"
)

// padding ...
func padding(src []byte, blocksize int) []byte {
	padnum := blocksize - len(src)%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	return append(src, pad...)
}

// unpadding ...
func unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	return src[:n-unpadnum]
}

// AESEncrypt 加密
func AESEncrypt(src []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	src = padding(src, block.BlockSize())
	blockmode := cipher.NewCBCEncrypter(block, key)
	blockmode.CryptBlocks(src, src)
	return src
}

// AESDecrypt 解密
func AESDecrypt(src []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	blockmode := cipher.NewCBCDecrypter(block, key)
	blockmode.CryptBlocks(src, src)
	src = unpadding(src)
	return src
}

// Test ...
func Test002() {
	src := "{\"account\":\"test_0\",\"type\":0,\"timestamp\":1690643120}"
	key := "111362EE140F157D"
	crypted := AesEncrypt(src, key)
	//BASE64.URLEncode(crypted)
	logs.Debugf(">>> %s", BASE64.Encode(crypted))
	//x2 := AESDecrypt(x1, key)
	//fmt.Print(string(x2))
}
