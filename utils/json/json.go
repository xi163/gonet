package json

import (
	"encoding/json"
	"io/ioutil"
)

/*
*
json流转换成struct/map
*/
func ParseFile(filename string, v any) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return Parse(b, &v)
}

/*
*
json流转换成struct/map
*/
func Parse(b []byte, v any) error {
	return json.Unmarshal(b, v)
}

/*
*
json串转换成struct/map
*/
func ParseStr(s string, v any) error {
	return json.Unmarshal([]byte(s), v)
}

/*
*
struct/map转换成json流
*/
func Bytes(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err.Error())
	}
	return b
}

/*
*
struct/map转换成json串
*/
func String(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}

/*
*
map转换成struct
*/
func MapToStruct(m, v any) error {
	return Parse(Bytes(m), v)
}

/*
*
struct转换成map
*/
func StructToMap(v, m any) error {
	return Parse(Bytes(v), m)
}
