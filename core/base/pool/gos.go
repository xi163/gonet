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

// go协程池(动态, 非阻塞)
type Gos interface {
	Start()
	Stop()
	Go(f cb.Functor)
	GoTimeout(d time.Duration, f cb.Functor, cb cb.Functor)
}

type gos struct {
	t task.Task
}

func NewGos() Gos {
	cpu := runtime.NumCPU()
	fixed := false
	nonblock := true
	s := &gos{
		t: task.NewGos("gos.task", cpu, cpu, fixed, nonblock, nil),
	}
	s.t.SetNew(mq.New(s.New))
	s.t.SetProcessor(cb.Processor(s.handler))
	return s
}

func (s *gos) New(v ...any) (q mq.Queue) {
	if t, ok := ch.NewChan(v[0].(int), v[1].(bool)).(mq.Queue); ok {
		q = t
		return
	}
	panic(errors.New("new mq error"))
}

func (s *gos) handler(msg any, args ...any) bool {
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

func (s *gos) Go(f cb.Functor) {
	s.t.Do(f)
}

func (s *gos) GoTimeout(d time.Duration, f cb.Functor, cb cb.Functor) {
	s.t.DoTimeout(d, f, cb)
}

func (s *gos) Start() {
	s.t.Start()
}

func (s *gos) Stop() {
	s.t.Stop()
}
