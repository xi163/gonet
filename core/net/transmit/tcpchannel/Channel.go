package tcpchannel

import (
	"github.com/cwloo/gonet/core/net/transmit"
	logs "github.com/cwloo/gonet/logs"
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

func (s *Channel) OnRecv(conn any) (int, any, error) {
	c, _ := conn.(net.Conn)
	if c == nil {
		logs.Fatalf("error")
	}
	buf := make([]byte, 1024)
	n, err := Read(c, buf)
	if err != nil {
		return 0, nil, err
	}
	buf = buf[0:n]
	return 0, buf, nil
}

func (s *Channel) OnSend(conn any, msg any, msgType int) error {
	c, _ := conn.(net.Conn)
	if c == nil {
		logs.Fatalf("error")
	}
	switch msg := msg.(type) {
	case string:
		return WriteFull(c, conv.StrToByte(msg))
	case []byte:
		return WriteFull(c, msg)
	default:
		b, _ := codec.Encode(msg)
		return WriteFull(c, b)
	}
}
