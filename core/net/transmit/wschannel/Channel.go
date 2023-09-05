package wschannel

import (
	"errors"

	"github.com/cwloo/gonet/core/net/transmit"
	logs "github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils/codec"
	"github.com/cwloo/gonet/utils/conv"

	"github.com/gorilla/websocket"
)

// websocket协议读写解析
type Channel struct {
}

func NewChannel() transmit.Channel {
	return &Channel{}
}

func (s *Channel) OnRecv(conn any) (any, error) {
	c, _ := conn.(*websocket.Conn)
	if c == nil {
		panic(errors.New("error"))
	}
	// c.SetReadLimit(1024)
	msgType, msg, err := c.ReadMessage()
	if msgType == websocket.PingMessage {
		logs.Infof("", "this is a pingMessage")
	}
	//TextMessage/BinaryMessage
	if websocket.TextMessage != msgType {
		panic(errors.New("error"))
	}
	return msg, err
}

func (s *Channel) OnSend(conn any, msg any) error {
	c, _ := conn.(*websocket.Conn)
	if c == nil {
		panic(errors.New("error"))
	}
	// c.SetWriteDeadline(time.Now().Add(time.Duration(60) * time.Second))
	switch msg := msg.(type) {
	case string:
		return c.WriteMessage(websocket.BinaryMessage, conv.StrToByte(msg))
	case []byte:
		return c.WriteMessage(websocket.BinaryMessage, msg)
	}
	b, _ := codec.Encode(msg)
	return c.WriteMessage(websocket.BinaryMessage, b)
}
