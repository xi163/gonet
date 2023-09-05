package tcpchannel

import (
	"errors"

	"github.com/cwloo/gonet/core/net/transmit"
	"github.com/cwloo/gonet/utils/codec"
	"github.com/cwloo/gonet/utils/conv"

	"net"
)

// TCP协议读写解析
type Channel struct {
}

func NewChannel() transmit.Channel {
	return &Channel{}
}

func (s *Channel) OnRecv(conn any) (any, error) {
	c, _ := conn.(net.Conn)
	if c == nil {
		panic(errors.New("error"))
	}
	buf := make([]byte, 1024)
	n, err := Read(c, buf)
	if err != nil {
		return nil, err
	}
	buf = buf[0:n]
	return buf, nil
}

func (s *Channel) OnSend(conn any, msg any) error {
	c, _ := conn.(net.Conn)
	if c == nil {
		panic(errors.New("error"))
	}
	switch msg := msg.(type) {
	case string:
		return WriteFull(c, conv.StrToByte(msg))
	case []byte:
		return WriteFull(c, msg)
	}
	b, _ := codec.Encode(msg)
	return WriteFull(c, b)
}
