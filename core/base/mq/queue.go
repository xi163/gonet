package mq

import (
	"sync"

	"github.com/cwloo/gonet/core/cb"
)

type New func(v ...any) Queue

// 消息队列
type Queue interface {
	Name() string
	Push(data any)
	Pop() (data any, exit, empty bool, code int)
	Pick() (v []any)
	Pick_until() (v []any, exit bool, code int)
	Exec(step bool, handler cb.Processor, args ...any) (exit bool, code int)
	Exec_until(step bool, handler cb.Processor, args ...any) (exit bool, code int)
	Size() int
}

// slice/list阻塞队列
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
