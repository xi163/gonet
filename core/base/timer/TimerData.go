package timer

import (
	"sync"
)

var (
	t = sync.Pool{
		New: func() any {
			return &data{}
		},
	}
)

type OpType uint8

const (
	RunAfter OpType = iota + 10
	RunAfterWith
	RunEvery
	RunEveryWith
	RemoveTimer
	RemoveTimers
)

// 定时器节点数据
type Data interface {
	OpType() OpType
	TimerId() uint32
	Delay() int32
	Interval() int32
	Args() []any
	TimerCallback() TimerCallback
	Cb() func(args ...any)
	Put()
}

type data struct {
	op              OpType
	timerID         uint32
	delay, interval int32
	args            []any
	handler         TimerCallback
	cb              func(args ...any)
}

func NewAfter(delay int32, cb func(args ...any), args ...any) Data {
	s := t.Get().(*data)
	s.op = RunAfter
	s.delay = delay
	s.args = args
	// s.args = append(s.args, args...)
	s.cb = cb
	return s
}

func NewAfterWith(delay int32, handler TimerCallback, cb func(args ...any), args ...any) Data {
	s := t.Get().(*data)
	s.op = RunAfterWith
	s.delay = delay
	s.handler = handler
	s.args = args
	// s.args = append(s.args, args...)
	s.cb = cb
	return s
}

func NewEvery(delay, interval int32, cb func(args ...any), args ...any) Data {
	s := t.Get().(*data)
	s.op = RunEvery
	s.delay = delay
	s.interval = interval
	s.args = args
	// s.args = append(s.args, args...)
	s.cb = cb
	return s
}

func NewEveryWith(delay, interval int32, handler TimerCallback, cb func(args ...any), args ...any) Data {
	s := t.Get().(*data)
	s.op = RunEveryWith
	s.delay = delay
	s.interval = interval
	s.handler = handler
	s.args = args
	// s.args = append(s.args, args...)
	s.cb = cb
	return s
}

func NewRemove(timerID uint32) Data {
	s := t.Get().(*data)
	s.op = RemoveTimer
	s.timerID = timerID
	return s
}

func NewRemoveAll() Data {
	s := t.Get().(*data)
	s.op = RemoveTimers
	return s
}

func (s *data) OpType() OpType {
	return s.op
}

func (s *data) TimerId() uint32 {
	return s.timerID
}

func (s *data) Delay() int32 {
	return s.delay
}

func (s *data) Interval() int32 {
	return s.interval
}

func (s *data) Args() []any {
	return s.args
}

func (s *data) TimerCallback() TimerCallback {
	return s.handler
}

func (s *data) Cb() func(args ...any) {
	return s.cb
}

func (s *data) Put() {
	t.Put(s)
}
