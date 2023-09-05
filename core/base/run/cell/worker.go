package cell

import (
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/core/net/conn"
)

// 业务处理单元
type Worker interface {
	OnInit()
	OnTimer(timerID uint32, dt int32, args ...any) bool
}

// WorkerCreator
type WorkerCreator interface {
	Create(proc run.Proc, args ...any) Worker
}

// NetWorker
type NetWorker interface {
	Worker
	OnConnected(peer conn.Session, v ...any)
	OnClosed(peer conn.Session, reason conn.Reason, v ...any)
	OnRead(cmd uint32, msg any, peer conn.Session)
	OnCustom(cmd uint32, msg any, peer conn.Session)
}

// NetWorkerCreator
type NetWorkerCreator interface {
	WorkerCreator
}
