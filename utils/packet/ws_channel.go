package packet

import (
	"encoding/binary"
	"errors"

	"github.com/cwloo/gonet/core/net/transmit"
	"github.com/gorilla/websocket"
)

// websocket协议读写解析
type WSChannel struct {
}

func NewWSChannel() transmit.Channel {
	return &WSChannel{}
}

func (s *WSChannel) OnRecv(conn any) (any, error) {
	c, _ := conn.(*websocket.Conn)
	if c == nil {
		panic("error")
	}
	//len+CRC，4字节
	// c.SetReadLimit(4)
	msgType, b, err := c.ReadMessage()
	if err != nil {
		// logs.Errorf(err.Error())
		return nil, err
	}
	switch msgType {
	case websocket.PingMessage:
		return nil, errors.New("error PingMessage")
	case websocket.TextMessage:
		return nil, errors.New("error TextMessage")
	case websocket.CloseMessage:
		return nil, errors.New("peer closed")
	}
	return b, err
}

func (s *WSChannel) OnSend(conn any, msg any) error {
	c, _ := conn.(*websocket.Conn)
	if c == nil {
		panic("error")
	}
	// c.SetWriteDeadline(time.Now().Add(time.Duration(60) * time.Second))
	switch msg := msg.(type) {
	case *Msg:
		b, _ := Pack(msg, binary.LittleEndian)
		return c.WriteMessage(websocket.BinaryMessage, b)
	}
	panic("error")
}
