package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

// 公钥加密
func PublicEncrypt(pt []byte, path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	// pem 解码
	block, _ := pem.Decode(buf)
	// x509 解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err.Error())
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	ct, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, pt)
	if err != nil {
		panic(err.Error())
	}
	//返回密文
	return ct
}

// 私钥解密
func PrivateDecrypt(ct []byte, path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	// pem 解码
	block, _ := pem.Decode(buf)
	// X509 解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err.Error())
	}
	//对密文进行解密
	pt, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ct)
	if err != nil {
		panic(err.Error())
	}
	//返回明文
	return pt
}
