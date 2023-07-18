package endian

import (
	"bytes"
	"encoding/binary"
)

/*
*
struct转换成byte二进制流
*/
func Encode(v any, order binary.ByteOrder) ([]byte, error) {
	// b := bytes.NewBuffer([]byte{})
	var b bytes.Buffer
	err := binary.Write(&b, order, v)
	return b.Bytes(), err
}

/*
*
byte二进制流转换成struct
*/
func Decode(a []byte, v any, order binary.ByteOrder) error {
	b := bytes.NewBuffer(a)
	return binary.Read(b, order, v)
}
