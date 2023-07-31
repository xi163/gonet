package sign

import (
	"math/rand"
	"time"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils/codec/base64"
	"github.com/cwloo/gonet/utils/crypto/aes"
	"github.com/cwloo/gonet/utils/json"
	"github.com/cwloo/gonet/utils/random"
)

type Sign struct {
	Rand      string `json:"rand" form:"rand"`
	Data      any    `json:"data" form:"data"`
	Timestamp int64  `json:"timestamp" form:"timestamp"`
	Expired   int64  `json:"expired" form:"expired"`
}

func Encode(data any, exp time.Time, secret []byte) string {
	token := Sign{
		Rand:      random.CharStr(rand.Int() % 10),
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
		Expired:   exp.UnixMilli(),
	}
	logs.Errorf("==>>> %s", json.String(token))
	encrypt := aes.ECBEncryptPKCS7(json.Bytes(token), secret, secret)
	//encrypt := test_crypto.AESEncrypt(json.Bytes(token), secret)
	strBase64 := base64.URLEncode(encrypt)
	//return uri.URLEncode(strBase64)
	return strBase64
}

func Decode(s string, secret []byte) (v any, exp int64) {
	//strBase64 := uri.URLDecode(s)
	strBase64 := s
	encrypt := base64.URLDecode(strBase64)
	jsonStr := aes.ECBDecryptPKCS7(encrypt, secret, secret)
	//jsonStr := test_crypto.AESDecrypt(encrypt, secret)
	model := Sign{}
	err := json.Parse(jsonStr, &model)
	if err != nil {
		logs.Errorf(err.Error())
		return nil, 0
	}
	return model.Data, model.Expired
}
