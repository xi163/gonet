package transmit

import (
	"io"
	"net"
)

// 消息传输接口(流协议读写解析)
type Channel interface {
	// 接收数据
	OnRecv(conn any) (any, error)
	// 发送数据
	OnSend(conn any, msg any) error
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
