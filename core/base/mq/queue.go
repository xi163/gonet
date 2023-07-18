package mq

import (
	"sync"

	"github.com/cwloo/gonet/core/cb"
)

type New func(v ...any) Queue

// <summary>
// Queue 消息队列
// <summary>
type Queue interface {
	Name() string
	// 入队列
	Push(data any)
	// 出队列 一次取一个
	Pop() (data any, exit, empty bool, code int)
	// 掏空队列 批量全部取
	Pick() (v []any)
	// 掏空队列 批量全部取直到遇到nil
	Pick_until() (v []any, exit bool, code int)
	// 一次取一个或批量全部取
	Exec(step bool, handler cb.Processor, args ...any) (exit bool, code int)
	// 一次取一个或批量全部取直到遇到nil
	Exec_until(step bool, handler cb.Processor, args ...any) (exit bool, code int)
	// 空闲队列大小
	Size() int
}

// <summary>
// BlockQueue slice/list阻塞队列
// <summary>
type BlockQueue interface {
	Queue
	Wakeup()
}

var (
	w = sync.Pool{
		New: func() any {
			return &WakeupStruct{}
		},
	}
)

type WakeupStruct struct{}

func NewWakeupStruct() *WakeupStruct {
	return w.Get().(*WakeupStruct)
}

func (s *WakeupStruct) Put() {
	w.Put(s)
}

type ExitStruct struct {
	Code int
}
