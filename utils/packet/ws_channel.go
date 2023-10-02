package packet

import (
	"encoding/binary"
	"errors"

	"github.com/cwloo/gonet/core/net/transmit"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils/codec"
	"github.com/gorilla/websocket"
)

// websocket协议读写解析
type WSChannel struct {
}

func NewWSChannel() transmit.Channel {
	return &WSChannel{}
}

func (s *WSChannel) OnRecv(conn any) (int, any, error) {
	c, _ := conn.(*websocket.Conn)
	if c == nil {
		logs.Fatalf("error")
	}
	//len+CRC，4字节
	// c.SetReadLimit(4)
	msgType, b, err := c.ReadMessage()
	if err != nil {
		// logs.Errorf(err.Error())
		return msgType, nil, err
	}
	switch msgType {
	//case websocket.PingMessage:
	//	return msgType, nil, errors.New("error PingMessage")
	//case websocket.TextMessage:
	//	return msgType, nil, errors.New("error TextMessage")
	case websocket.CloseMessage:
		return msgType, nil, errors.New("peer closed")
	}
	return msgType, b, err
}

func (s *WSChannel) OnSend(conn any, msg any, msgType int) error {
	c, _ := conn.(*websocket.Conn)
	if c == nil {
		logs.Fatalf("error")
	}
	switch msgType {
	case websocket.TextMessage:
		// c.SetWriteDeadline(time.Now().Add(time.Duration(60) * time.Second))
		switch msg := msg.(type) {
		case string:
			return c.WriteMessage(msgType, []byte(msg))
		case []byte:
			return c.WriteMessage(msgType, msg)
		}
	case websocket.BinaryMessage:
		// c.SetWriteDeadline(time.Now().Add(time.Duration(60) * time.Second))
		switch msg := msg.(type) {
		case string:
			return c.WriteMessage(msgType, []byte(msg))
		case []byte:
			return c.WriteMessage(msgType, msg)
		case *Msg:
			b, _ := Pack(msg, binary.LittleEndian)
			return c.WriteMessage(msgType, b)
		default:
			b, _ := codec.Encode(msg)
			return c.WriteMessage(msgType, b)
		}
	default:
		logs.Fatalf("msg type")
	}
	panic("error")
}
