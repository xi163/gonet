package transmit

import (
	"io"
	"net"
)

// 发送websocket消息
type Messagetruct struct {
	Type int
	Msg  any
}

// 消息传输接口(流协议读写解析)
type Channel interface {
	// 接收数据
	OnRecv(conn any) (int, any, error)
	// 发送数据
	OnSend(conn any, msg any, msgType int) error
}

func IsEOFOrReadError(err error) bool {
	if err == io.EOF {
		return true
	}
	ne, ok := err.(*net.OpError)
	return ok && ne.Op == "read"
}

func IsEOFOrWriteError(err error) bool {
	if err == io.EOF {
		return true
	}
	ne, ok := err.(*net.OpError)
	return ok && ne.Op == "write"
}
