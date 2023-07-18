package packet

import (
	"errors"
)

type RootMsg struct {
	Cmd  uint32
	Data []byte
}

func Decode(msg any) (uint32, any, error) {
	if msg == nil {
		panic("error")
	}
	// switch msg := msg.(type) {
	// case *RootMsg:
	// 	data, err := codec.Decode(msg.Data)
	// 	return msg.Cmd, data, err
	// }
	return 0, nil, errors.New("error msgtype")
}
