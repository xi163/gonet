package wschannel

import (
	"errors"

	"github.com/cwloo/gonet/core/net/transmit"
	logs "github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils/codec"
	"github.com/gorilla/websocket"
)

// websocket协议读写解析
type Channel struct {
}

func NewChannel() transmit.Channel {
	return &Channel{}
}

func (s *Channel) OnRecv(conn any) (int, any, error) {
	c, _ := conn.(*websocket.Conn)
	if c == nil {
		logs.Fatalf("error")
	}
	// c.SetReadLimit(1024)
	msgType, b, err := c.ReadMessage()
	if err != nil {
		logs.Errorf(err.Error())
		return msgType, nil, err
	}
	switch msgType {
	// case websocket.PingMessage:
	// 	return msgType, nil, errors.New("error PingMessage")
	// case websocket.TextMessage:
	// 	return msgType, nil, errors.New("error TextMessage")
	case websocket.CloseMessage:
		return msgType, nil, errors.New("peer closed")
	}
	return msgType, b, err
}

func (s *Channel) OnSend(conn any, msg any, msgType int) error {
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
		default:
			b, _ := codec.Encode(msg)
			return c.WriteMessage(msgType, b)
		}
	default:
		logs.Fatalf("msg type")
	}
	panic("error")
}
