package pool

import (
	"errors"
	"runtime"
	"time"

	"github.com/cwloo/gonet/core/base/mq"
	"github.com/cwloo/gonet/core/base/mq/ch"
	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/utils/safe"
)

// <summary>
// Calls 回调池(固定, 非阻塞)
// <summary>
type Calls interface {
	Start()
	Stop()
	Call(f cb.Functor)
	CallTimeout(d time.Duration, f cb.Functor, cb cb.Functor)
}

type calls struct {
	t task.Task
}

func NewCalls() Calls {
	cpu := runtime.NumCPU()
	fixed := true
	nonblock := true
	s := &calls{
		t: task.NewGos("calls.task", cpu, cpu, fixed, nonblock, nil),
	}
	s.t.SetNew(mq.New(s.New))
	s.t.SetProcessor(cb.Processor(s.handler))
	return s
}

func (s *calls) New(v ...any) (q mq.Queue) {
	if t, ok := ch.NewChan(v[0].(int), v[1].(bool)).(mq.Queue); ok {
		q = t
		return
	}
	panic(errors.New("new mq error"))
}

func (s *calls) handler(msg any, args ...any) bool {
	switch msg := msg.(type) {
	case cb.Functor:
		safe.Call2(msg.Call)
		msg.Put()
	case cb.Timeout:
		if !msg.Expire().Expired(time.Now()) {
			switch data := msg.Data().(type) {
			case cb.Functor:
				data.CallWith(msg.Expire())
				data.Put()
			}
		}
		msg.Put()
	}
	return false
}

func (s *calls) Call(f cb.Functor) {
	s.t.Do(f)
}

func (s *calls) CallTimeout(d time.Duration, f cb.Functor, cb cb.Functor) {
	s.t.DoTimeout(d, f, cb)
}

func (s *calls) Start() {
	s.t.Start()
}

func (s *calls) Stop() {
	s.t.Stop()
}
