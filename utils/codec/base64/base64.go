package base64

import (
	BASE64 "encoding/base64"
)

func Encode(b []byte) string {
	return BASE64.StdEncoding.EncodeToString(b)
}

func Decode(s string) []byte {
	b, err := BASE64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err.Error())
	}
	return b
}

func RawEncode(b []byte) string {
	return BASE64.RawStdEncoding.EncodeToString(b)
}

func RawDecode(s string) []byte {
	b, err := BASE64.RawStdEncoding.DecodeString(s)
	if err != nil {
		panic(err.Error())
	}
	return b
}

func URLEncode(b []byte) string {
	return BASE64.URLEncoding.EncodeToString(b)
}

func URLDecode(s string) []byte {
	b, err := BASE64.URLEncoding.DecodeString(s)
	if err != nil {
		panic(err.Error())
	}
	return b
}
