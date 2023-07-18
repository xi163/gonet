package tcpclient

import (
	"time"

	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/core/net/conn"
)

// <summary>
// TCPClient TCP客户端
// <summary>
type TCPClient interface {
	Name() string
	Peers() conn.Sessions
	ConnectTCP(address ...string)
	Reconnect()
	Disconnect()
	Retry() bool
	EnableRetry(retry bool)
	Range(cb func(peer conn.Session))
	SetHoldType(hold conn.HoldType)
	SetDialTimeout(d time.Duration)
	SetIdleTimeout(timeout, d time.Duration)
	SetRetryInterval(d time.Duration)
	SetProtocolCallback(cb cb.OnProtocol)
	SetConnectedCallback(cb cb.OnConnected)
	SetClosedCallback(cb cb.OnClosed)
	SetMessageCallback(cb cb.OnMessage)
	SetWriteCompleteCallback(cb cb.OnWriteComplete)
}
