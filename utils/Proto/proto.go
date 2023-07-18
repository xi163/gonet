package Proto

import "google.golang.org/protobuf/proto"

func Encode(m proto.Message) ([]byte, error) {
	return proto.Marshal(m)
}

func Decode(a []byte, m proto.Message) error {
	return proto.Unmarshal(a, m)
}
