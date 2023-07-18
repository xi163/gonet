package codec

import (
	"bytes"
	"encoding/gob"
)

/*
*
struct/map/any编码成byte
*/
func Encode(v any) ([]byte, error) {
	// b := bytes.NewBuffer(nil)
	// _ = gob.NewEncoder(b).Encode(v)
	var b bytes.Buffer
	// binary.Write(&b, binary.LittleEndian, v)
	err := gob.NewEncoder(&b).Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), err
}

/*
*
byte解码成struct/map/any
*/
func Decode(a []byte, v any) error {
	b := bytes.NewBuffer(a)
	err := gob.NewDecoder(b).Decode(v)
	return err
}
