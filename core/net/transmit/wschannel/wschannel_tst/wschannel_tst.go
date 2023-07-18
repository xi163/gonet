package wschannel_tst

import (
	"errors"

	"github.com/cwloo/gonet/core/net/transmit"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils/codec"
	"github.com/gorilla/websocket"
)

// <summary>
// WSChannel websocket协议读写解析
// <summary>
type WSChannel struct {
}

func NewWSChannel() transmit.Channel {
	return &WSChannel{}
}

func (s *WSChannel) OnRecv(conn any) (any, error) {
	c, _ := conn.(*websocket.Conn)
	if c == nil {
		panic(errors.New("error"))
	}
	// c.SetReadLimit(1024)
	msgType, b, err := c.ReadMessage()
	if err != nil {
		logs.Errorf(err.Error())
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
	// logs.Warnf("%v", string(msg.([]byte)))
	c, _ := conn.(*websocket.Conn)
	if c == nil {
		panic(errors.New("error"))
	}
	// c.SetWriteDeadline(time.Now().Add(time.Duration(60) * time.Second))
	switch msg := msg.(type) {
	case string:
		return c.WriteMessage(websocket.BinaryMessage, []byte(msg))
	case []byte:
		return c.WriteMessage(websocket.BinaryMessage, msg)
	}
	b, _ := codec.Encode(msg)
	return c.WriteMessage(websocket.BinaryMessage, b)
}
