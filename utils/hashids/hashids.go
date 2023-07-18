package hashids

import (
	go_hashids "github.com/speps/go-hashids"
)

func Encrypt(salt string, minLength int, params []int) string {
	hd := go_hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, err := go_hashids.NewWithData(hd)
	if err == nil {
		e, err := h.Encode(params)
		if err == nil {
			return e
		}
	}
	return ""
}

func Decrypt(salt string, minLength int, hash string) []int {
	hd := go_hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, err := go_hashids.NewWithData(hd)
	if err == nil {
		e, err := h.DecodeWithError(hash)
		if err == nil {
			return e
		}
	}
	return []int{}
}
