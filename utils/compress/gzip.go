package compress

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func Gzip(msg []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(msg)
	if err != nil {
		return msg, err
	}
	_ = gz.Close()
	msg = b.Bytes()
	return msg, err
}

func Gunzip(msg []byte) ([]byte, error) {
	b := bytes.NewBuffer(msg)
	reader, err := gzip.NewReader(b)
	if err != nil {
		return msg, err
	}
	msg, err = ioutil.ReadAll(reader)
	if err != nil {
		return msg, err
	}
	_ = reader.Close()
	return msg, err
}
