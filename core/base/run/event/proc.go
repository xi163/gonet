package event

import (
	"time"

	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/core/net/conn"
)

// <summary>
// Proc 事件处理单元
// <summary>
type Proc interface {
	Post(data *Data)
	PostConnected(peer conn.Session, v ...any)
	PostConnectedWith(handler cb.OnConnected, peer conn.Session, v ...any)
	PostClosing(d time.Duration, peer conn.Session)
	PostClosed(peer conn.Session, reason conn.Reason, v ...any)
	PostClosedWith(handler cb.OnClosed, peer conn.Session, reason conn.Reason, v ...any)
	PostRead(cmd uint32, msg any, peer conn.Session)
	PostReadWith(handler cb.ReadCallback, cmd uint32, msg any, peer conn.Session)
	PostCustom(cmd uint32, msg any, peer conn.Session)
	PostCustomWith(handler cb.CustomCallback, cmd uint32, msg any, peer conn.Session)
}
