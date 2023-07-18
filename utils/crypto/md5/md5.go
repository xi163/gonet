package md5

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Md5Byte(b []byte, upper bool, salt ...string) string {
	h := md5.New()
	h.Write(b)
	if len(salt) > 0 {
		h.Write([]byte(salt[0]))
	}
	switch upper {
	case true:
		// logs.Debugf(fmt.Sprintf("%X", h.Sum(nil)))
		return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	}
	// logs.Debugf(fmt.Sprintf("%x", h.Sum(nil)))
	return strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}

func Md5(s string, upper bool, salt ...string) string {
	return Md5Byte([]byte(s), upper, salt...)
}
