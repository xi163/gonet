package task

import (
	"errors"
	"runtime"
	"time"

	"github.com/cwloo/gonet/core/base/mq"
	"github.com/cwloo/gonet/core/base/mq/ch"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/core/base/run/timeout"
	"github.com/cwloo/gonet/core/base/timer"
	"github.com/cwloo/gonet/core/cb"
)

// 定时任务池(固定, 非阻塞)
var (
	timeouts = NewTimeout()
)

func After(d time.Duration, cb cb.Functor) {
	timeouts.After(d, cb)
}

// 定时任务池(单生产者，多消费者)
type Timeout interface {
	After(d time.Duration, f cb.Functor)
}

type timeoutTask struct {
	t Task
}

func NewTimeout() Timeout {
	cpu := runtime.NumCPU()
	fixed := true
	nonblock := true
	tick := false
	d := time.Second
	s := &timeoutTask{}
	runner := timeout.NewProcessor(tick, d, s.handler, s.onTimer)
	s.t = NewTask("timeout.task", cpu, cpu, fixed, nonblock, runner)
	s.t.SetNew(mq.New(s.New))
	s.t.SetProcessor(cb.Processor(s.handler))
	s.t.Start()
	return s
}

func (s *timeoutTask) New(v ...any) (q mq.Queue) {
	if t, ok := ch.NewChan(v[0].(int), v[1].(bool)).(mq.Queue); ok {
		q = t
		return
	}
	panic(errors.New("new mq error"))
}

// func (s *timeoutTask) overload(r run.Processor) (n int, b bool) {
// 	if q, ok := r.Queue().(ch.Queue); ok {
// 		n = 1 + q.Length() + q.Size()
// 		if n > r.IdleCount() {
// 			// b = true
// 		}
// 	} else {
// 		n = 1 + r.Queue().Size()
// 		if n > r.IdleCount() {
// 			// b = true
// 		}
// 	}
// 	return
// }

// func (s *timeoutTask) gcCondition(r run.Processor, args run.Args) (n int, b bool) {
// 	if q, ok := r.Queue().(ch.Queue); ok {
// 		n = q.Length() + q.Size()
// 		if r.IdleCount() > 0 {
// 			// b = true
// 		}
// 	} else {
// 		n = r.Queue().Size()
// 		if r.IdleCount() > 0 {
// 			// b = true
// 		}
// 	}
// 	return
// }

func (s *timeoutTask) onTimer(timerID uint32, dt int32, args ...any) bool {
	if len(args) == 0 {
		panic(errors.New("timeoutTask.args 0"))
	}
	if args[0] == nil {
		panic(errors.New("timeoutTask.args[0] is nil"))
	}
	switch f := args[0].(type) {
	case cb.Functor:
		f.Call()
		f.Put()
	}
	return true
}

func (s *timeoutTask) handler(msg any, args ...any) bool {
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
	switch msg := msg.(type) {
	case timer.Data:
		switch msg.OpType() {
		case timer.RunAfter:
			timerId := arg.RunAfter(msg.Delay(), msg.Args()...)
			msg.Cb()(timerId)
		case timer.RunAfterWith:
			timerId := arg.RunAfterWith(msg.Delay(), msg.TimerCallback(), msg.Args()...)
			msg.Cb()(timerId)
		case timer.RunEvery:
			timerId := arg.RunEvery(msg.Delay(), msg.Interval(), msg.Args()...)
			msg.Cb()(timerId)
		case timer.RunEveryWith:
			timerId := arg.RunEveryWith(msg.Delay(), msg.Interval(), msg.TimerCallback(), msg.Args()...)
			msg.Cb()(timerId)
		case timer.RemoveTimer:
			arg.RemoveTimer(msg.TimerId())
		case timer.RemoveTimers:
			arg.RemoveTimers()
		}
		msg.Put()
	}
	return false
}

func (s *timeoutTask) After(d time.Duration, cb cb.Functor) {
	second := int32(int64(d) / int64(time.Second)) //秒
	millisec := second * 1000                      //毫秒
	s.t.Do(timer.NewAfter(millisec, func(args ...any) {
	}, cb))
}
