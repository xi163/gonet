package pool

import (
	"errors"
	"runtime"
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/mq"
	"github.com/cwloo/gonet/core/base/mq/ch"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/utils/safe"
)

// <summary>
// Conns 连接池(固定, 阻塞)
// <summary>
type Conns interface {
	Start()
	Stop()
	Do(f cb.Functor)
	DoTimeout(d time.Duration, f cb.Functor, cb cb.Functor)
	Num() int
	ResetNum()
}

type conns struct {
	t task.Task
	c cc.Counter
}

func NewConns() Conns {
	cpu := runtime.NumCPU()
	cpu = 1
	fixed := true    //固定协程数量
	nonblock := true //阻塞
	s := &conns{
		c: cc.NewAtomCounter(),
		t: task.NewGos("conns.task", cpu, 2*cpu, fixed, nonblock, nil),
	}
	s.t.SetNew(mq.New(s.New))
	s.t.SetProcessor(cb.Processor(s.handler))
	s.t.SetOverload(run.Overload(s.overload))
	s.t.SetGcCondition(run.GcCondition(s.gcCondition))
	return s
}

func (s *conns) New(v ...any) (q mq.Queue) {
	if t, ok := ch.NewChan(v[0].(int), v[1].(int), v[2].(bool)).(mq.Queue); ok {
		q = t
		return
	}
	panic(errors.New("new mq error"))
}

// 过载判断
func (s *conns) overload(r run.Processor) (n int, b bool) {
	if q, ok := r.Queue().(ch.Queue); ok {
		n = 1 + q.Length() + q.Size()
		if n > r.IdleCount() {
			b = true
		}
	} else {
		n = 1 + r.Queue().Size()
		if n > r.IdleCount() {
			b = true
		}
	}
	return
}

// GC垃圾回收条件检查
func (s *conns) gcCondition(r run.Processor, args run.Args) (n int, b bool) {
	if q, ok := r.Queue().(ch.Queue); ok {
		n = q.Length() + q.Size()
		if r.IdleCount() > 0 {
			b = true
		}
	} else {
		n = r.Queue().Size()
		if r.IdleCount() > 0 {
			b = true
		}
	}
	return
}

func (s *conns) handler(msg any, args ...any) bool {
	s.c.Up()
	switch msg.(type) {
	case *cb.Functor00, *cb.Functor10, *cb.Functor20, *cb.Functor01, *cb.Functor11, *cb.Functor21:
		data, _ := msg.(cb.Functor)
		safe.Call2(data.Call)
		data.Put()
		break
	case cb.Timeout:
		timeout, _ := msg.(cb.Timeout)
		if !timeout.Expire().Expired(time.Now()) {
			switch timeout.Data().(type) {
			case *cb.Functor00, *cb.Functor10, *cb.Functor20, *cb.Functor01, *cb.Functor11, *cb.Functor21:
				data, _ := msg.(cb.Functor)
				// safe.Call2(data.Call)
				data.CallWith(timeout.Expire())
				data.Put()
				break
			}
		}
		timeout.Put()
		break
	}
	// logs.Debugf("NumProcessed:%v goroutines.idles:%v goroutines.total:%v", s.Num(), s.t.Runner().IdleCount(), s.t.Runner().Count())
	return false
}

func (s *conns) Num() int {
	return s.c.Count()
}

func (s *conns) ResetNum() {
	s.c.Reset()
}

func (s *conns) Do(f cb.Functor) {
	s.t.Runner().IdleDown() //空闲协程数量递减
	s.t.Do(f)
}

func (s *conns) DoTimeout(d time.Duration, f cb.Functor, cb cb.Functor) {
	s.t.Runner().IdleDown() //空闲协程数量递减
	s.t.DoTimeout(d, f, cb)
}

func (s *conns) Start() {
	s.t.Start()
}

func (s *conns) Stop() {
	s.t.Stop()
}
