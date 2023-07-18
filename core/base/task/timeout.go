package task

import (
	"errors"
	"runtime"
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/mq"
	"github.com/cwloo/gonet/core/base/mq/ch"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/core/base/run/timeout"
	"github.com/cwloo/gonet/core/base/timer"
	"github.com/cwloo/gonet/core/cb"
)

// 超时任务池(固定, 非阻塞)
var (
	timeouts = NewTimeout()
)

func After(d time.Duration, cb cb.Functor) {
	timeouts.After(d, cb)
}

// <summary>
// Timeout 超时任务池(单生产者，多消费者)
// <summary>
type Timeout interface {
	After(d time.Duration, f cb.Functor)
}

type timeoutTask struct {
	t Task
	c cc.Counter
}

func NewTimeout() Timeout {
	cpu := runtime.NumCPU()
	fixed := true    //固定协程数量
	nonblock := true //非阻塞
	tick := false    //开启tick检查
	d := time.Second //tick间隔时间
	s := &timeoutTask{c: cc.NewAtomCounter()}
	runner := timeout.NewProcessor(tick, d, s.handler, s.onTimer)
	s.t = NewTask("timeout.task", cpu, cpu, fixed, nonblock, runner)
	s.t.SetNew(mq.New(s.New))
	s.t.SetProcessor(cb.Processor(s.handler))
	s.t.SetOverload(run.Overload(s.overload))
	s.t.SetGcCondition(run.GcCondition(s.gcCondition))
	s.t.Start()
	return s
}

func (s *timeoutTask) New(v ...any) (q mq.Queue) {
	if t, ok := ch.NewChan(v[0].(int), v[1].(int), v[2].(bool)).(mq.Queue); ok {
		q = t
		return
	}
	panic(errors.New("new mq error"))
}

// 过载判断
func (s *timeoutTask) overload(r run.Processor) (n int, b bool) {
	if q, ok := r.Queue().(ch.Queue); ok {
		n = 1 + q.Length() + q.Size()
		if n > r.IdleCount() {
			// b = true
		}
	} else {
		n = 1 + r.Queue().Size()
		if n > r.IdleCount() {
			// b = true
		}
	}
	return
}

// GC垃圾回收条件检查
func (s *timeoutTask) gcCondition(r run.Processor, args run.Args) (n int, b bool) {
	if q, ok := r.Queue().(ch.Queue); ok {
		n = q.Length() + q.Size()
		if r.IdleCount() > 0 {
			// b = true
		}
	} else {
		n = r.Queue().Size()
		if r.IdleCount() > 0 {
			// b = true
		}
	}
	return
}

func (s *timeoutTask) onTimer(timerID uint32, dt int32, args ...any) bool {
	if len(args) == 0 {
		panic(errors.New("timeoutTask.args 0"))
	}
	if args[0] == nil {
		panic(errors.New("timeoutTask.args[0] is nil"))
	}
	switch args[0].(type) {
	case cb.Functor:
		f, _ := args[0].(cb.Functor)
		f.Call()
		f.Put()
		break
	}
	return true
}

func (s *timeoutTask) handler(msg any, args ...any) bool {
	s.c.Up()
	if len(args) < 1 {
		panic(errors.New("args.size"))
	}
	proc, ok := args[0].(run.Proc)
	if !ok {
		panic(errors.New("arg[0]"))
	}
	arg, ok := proc.Args().(*timeout.Args)
	if !ok {
		panic(errors.New(""))
	}
	switch msg.(type) {
	case timer.Data:
		data, _ := msg.(timer.Data)
		switch data.OpType() {
		case timer.RunAfter:
			timerId := arg.RunAfter(data.Delay(), data.Args()...)
			data.Cb()(timerId)
			break
		case timer.RunAfterWith:
			timerId := arg.RunAfterWith(data.Delay(), data.TimerCallback(), data.Args()...)
			data.Cb()(timerId)
			break
		case timer.RunEvery:
			timerId := arg.RunEvery(data.Delay(), data.Interval(), data.Args()...)
			data.Cb()(timerId)
			break
		case timer.RunEveryWith:
			timerId := arg.RunEveryWith(data.Delay(), data.Interval(), data.TimerCallback(), data.Args()...)
			data.Cb()(timerId)
			break
		case timer.RemoveTimer:
			arg.RemoveTimer(data.TimerId())
			break
		case timer.RemoveTimers:
			arg.RemoveTimers()
			break
		}
		data.Put()
		break
	}
	return false
}

func (s *timeoutTask) Num() int {
	return s.c.Count()
}

func (s *timeoutTask) ResetNum() {
	s.c.Reset()
}

func (s *timeoutTask) After(d time.Duration, cb cb.Functor) {
	s.t.Runner().IdleDown()                        //空闲协程数量递减
	second := int32(int64(d) / int64(time.Second)) //秒
	millisec := second * 1000                      //毫秒
	s.t.Do(timer.NewAfter(millisec, func(args ...any) {
	}, cb))
}
