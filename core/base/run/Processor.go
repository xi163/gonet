package run

import (
	"github.com/cwloo/gonet/core/base/mq"

	"github.com/cwloo/gonet/core/cb"
)

type EndType uint8

const (
	QUIT EndType = iota + 1
	GC
	STOP
)

type GcCondition func(r Processor, args Args) (int, bool)

type Overload func(r Processor) (int, bool)

// <summary>
// Processor 执行消息队列
// <summary>
type Processor interface {
	Name() string
	Queue() mq.Queue
	SetQueue(q mq.Queue)
	SetProcessor(handler cb.Processor)
	NewArgs(proc Proc) Args
	Run(proc Proc)
	Wait()
}
