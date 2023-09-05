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

// 连接池(固定, 阻塞)
type Conns interface {
	Start()
	Stop()
	Do(f cb.Functor)
	DoTimeout(d time.Duration, f cb.Functor, cb cb.Functor)
}

type conns struct {
	t task.Task
}

func NewConns() Conns {
	cpu := runtime.NumCPU()
	fixed := true
	nonblock := false
	s := &conns{
		t: task.NewGos("conns.task", cpu, cpu, fixed, nonblock, nil),
	}
	s.t.SetNew(mq.New(s.New))
	s.t.SetProcessor(cb.Processor(s.handler))
	return s
}

func (s *conns) New(v ...any) (q mq.Queue) {
	if t, ok := ch.NewChan(v[0].(int), v[1].(bool)).(mq.Queue); ok {
		q = t
		return
	}
	panic(errors.New("new mq error"))
}

func (s *conns) handler(msg any, args ...any) bool {
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

func (s *conns) Do(f cb.Functor) {
	s.t.Do(f)
}

func (s *conns) DoTimeout(d time.Duration, f cb.Functor, cb cb.Functor) {
	s.t.DoTimeout(d, f, cb)
}

func (s *conns) Start() {
	s.t.Start()
}

func (s *conns) Stop() {
	s.t.Stop()
}
