package tcpserver

import (
	"time"

	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/core/net/conn"
)

// <summary>
// TCPServer TCP服务端
// <summary>
type TCPServer interface {
	Name() string
	Peers() conn.Sessions
	ListenAddr() *conn.Address
	ListenTCP(address ...string)
	Stop()
	Range(cb func(peer conn.Session))
	SetHoldType(holdType conn.HoldType)
	SetProtocolCallback(cb cb.OnProtocol)
	SetHandshakeCallback(cb cb.OnHandshake)
	SetConditionCallback(cb cb.OnCondition)
	SetConnectedCallback(cb cb.OnConnected)
	SetClosedCallback(cb cb.OnClosed)
	SetMessageCallback(cb cb.OnMessage)
	SetWriteCompleteCallback(cb cb.OnWriteComplete)
	SetCertFile(certfile, keyfile string)
	SetHandshakeTimeout(d time.Duration)
	SetIdleTimeout(timeout, d time.Duration)
	SetReadBufferSize(readBufferSize int)
}
