package tickers

import (
	"errors"
	"runtime"
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/mq"
	"github.com/cwloo/gonet/core/base/mq/ch"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/core/cb"
)

// <summary>
// Processor 执行消息队列
// <summary>
type Processor struct {
	run.Processor
	mq          ch.Queue
	counter     cc.Counter
	idleCounter cc.Counter
	tick        bool
	d           time.Duration
	handler     cb.Processor
	gcCondition run.GcCondition
}

func NewProcessor(tick bool, d time.Duration, handler cb.Processor) run.Processor {
	s := &Processor{
		tick:        tick,
		d:           d,
		handler:     handler,
		counter:     cc.NewAtomCounter(),
		idleCounter: cc.NewAtomCounter(),
	}
	s.SetGcCondition(s.GcCondition)
	return s
}

func NewProcessorWith(q ch.Queue, tick bool, d time.Duration, handler cb.Processor) run.Processor {
	s := &Processor{
		tick:        tick,
		mq:          q,
		d:           d,
		handler:     handler,
		counter:     cc.NewAtomCounter(),
		idleCounter: cc.NewAtomCounter(),
	}
	s.SetGcCondition(s.GcCondition)
	return s
}

func (s *Processor) SetProcessor(handler cb.Processor) {
	s.handler = handler
}

func (s *Processor) SetGcCondition(handler run.GcCondition) {
	s.gcCondition = handler
}

func (s *Processor) Name() string {
	return "tickers.Processor"
}

func (s *Processor) assertQueue() {
	if s.mq == nil {
		panic(errors.New("tickers.Processor.mq is nil"))
	}
}

func (s *Processor) Queue() mq.Queue {
	s.assertQueue()
	return s.mq
}

func (s *Processor) SetQueue(q mq.Queue) {
	if mq, ok := q.(ch.Queue); ok {
		s.mq = mq
	} else {
		panic(errors.New("need ch.Queue"))
	}
}

func (s *Processor) NewArgs(proc run.Proc) run.Args {
	return newArgs(proc, s.d)
}

func (s *Processor) startTicker(arg run.Args, args ...any) {
	if s.tick && s.d > 0 {
		second := int32(int64(s.d) / int64(time.Second)) //秒
		millisec := second * 1000                        //毫秒
		arg.RunEvery(0, millisec, args...)
	}
}

func (s *Processor) Run(proc run.Proc) {
	// logs.Debugf("%s started...", proc.Name())
	if s.mq == nil {
		panic(errors.New("error: tickers.Processor.mq is nil"))
	}
	if s.handler == nil {
		panic(errors.New("error: tickers.Processor.handler is nil"))
	}
	if s.gcCondition == nil {
		panic(errors.New("error: tickers.Processor.gcCondition is nil"))
	}
	if proc.Args() == nil {
		panic(errors.New("error: tickers.Processor.args is nil"))
	}
	ticker := proc.Args().(*Args).ticker
	if ticker == nil {
		panic(errors.New("error: tickers.Processor.ticker is nil"))
	}
	trigger := proc.Args().(*Args).Trigger()
	if trigger == nil {
		panic(errors.New("error: tickers.Processor.trigger is nil"))
	}
	timer := proc.Args().(*Args).timer
	if timer == nil {
		panic(errors.New("error: tickers.Processor.timer is nil"))
	}
	// timeout := proc.Args().(*Args).timeout
	// if timeout == nil {
	// 	panic(errors.New("error: tickers.Processor.timeout is nil"))
	// }
	// timerv2 := proc.Args().(*Args).timerv2
	// if timerv2 == nil {
	// 	panic(errors.New("error: tickers.Processor.timerv2 is nil"))
	// }
	timercb := proc.Args().(*Args).TimerCallback()
	if timercb == nil {
		panic(errors.New("error: tickers.Processor.timerCb is nil"))
	}
	s.counter.Up()
	// s.idleCounter.Up()
	arg := proc.Args().(*Args)
	s.startTicker(proc.Args(), proc.Args())
	tickerGC := run.NewTrigger(10 * time.Second)
	// tickerTimeout := run.NewTrigger(5 * time.Second)
	flag := run.STOP
	i, t := 0, 200
EXIT:
	// for !arg.stopping.Signaled() {
	for {
		if i > t {
			i = 0
			runtime.GC()
			// runtime.Gosched()
		}
		i++
		select {
		case <-arg.stopping.Read():
			if s.Count() == 1 {
				s.mq.AssertEmpty()
				s.flush(arg, proc)
			}
			flag = run.QUIT
			break EXIT
		case <-trigger:
			s.begin(arg)
			timer.Poll(proc.Tid(), timercb)
			s.end(arg)
			break
		// case c, _ := <-timerv2.Do():
		// 	s.begin(arg)
		// 	safe.Call(c.Call)
		// 	s.end(arg)
		// 	break
		case msg, ok := <-s.mq.Read():
			if ok {
				if msg == nil {
					// logs.Debugf("%v exit", proc.Name())
					// panic(errors.New("error: msg is nil"))
					s.mq.Reset()
				} else {
					s.begin(arg)
					s.handler(msg, arg)
					exit, _ := s.mq.Exec_until(false, s.handler, proc)
					if exit {
						s.mq.Reset()
						s.end(arg)
						break
					}
					s.end(arg)
				}
			} else {
				if msg == nil {
					s.mq.AssertEmpty()
					s.flush(arg, proc)
					flag = run.STOP
					break EXIT
				} else {
					panic(errors.New("error: channel closed"))
				}
			}
			break
		case <-s.mq.Signal():
			s.begin(arg)
			exit, _ := s.mq.Exec_until(false, s.handler, proc)
			if exit {
				s.mq.Reset()
				s.end(arg)
				break
			}
			s.end(arg)
			break
		// case <-tickerTimeout.Trigger():
		// 	s.begin(arg)
		// 	timeout.Poll(proc.Tid(), timercb)
		// 	s.end(arg)
		// 	break
		case <-tickerGC.Trigger():
			if s.Gc(proc.Args()) {
				if s.Count() == 1 {
					s.mq.AssertEmpty()
					s.flush(arg, proc)
				}
				flag = run.GC
				break EXIT
			} else {
				// logs.Debugf("mq.len:%v mq.size:%v goroutines.idles:%v goroutines.total:%v", s.mq.Length(), s.mq.Size(), s.IdleCount(), s.Count())
			}
			break
			// case <-time.After(5 * time.Second):
			// 	logs.Debugf("mq.len:%v mq.size:%v", s.mq.Length(), s.mq.Size())
			// 	break
			//轮询时默认case会导致CPU负载非常高，应该禁用
			//default:
		}
	}
	timer.RemoveTimers()
	ticker.Stop()
	tickerGC.Stop()
	s.idleCounter.Down()
	s.counter.Down()
	s.trace(proc.Name(), flag)
}

func (s *Processor) trace(name string, flag run.EndType) {
	switch flag {
	case run.QUIT:
		// logs.Debugf("*** QUIT *** %v mq.len:%v mq.size:%v goroutines.idles:%v goroutines.total:%v", name, s.mq.Length(), s.mq.Size(), s.IdleCount(), s.Count())
		break
	case run.GC:
		// logs.Debugf("*** GC *** %v mq.len:%v mq.size:%v goroutines.idles:%v goroutines.total:%v", name, s.mq.Length(), s.mq.Size(), s.IdleCount(), s.Count())
		break
	case run.STOP:
		// logs.Debugf("*** STOP *** %v mq.len:%v mq.size:%v goroutines.idles:%v goroutines.total:%v", name, s.mq.Length(), s.mq.Size(), s.IdleCount(), s.Count())
		break
	default:
		panic(errors.New(""))
	}
}

func (s *Processor) IdleUp() {
	s.idleCounter.Up()
}

func (s *Processor) IdleDown() {
	s.idleCounter.Down()
}

func (s *Processor) begin(arg run.Args) {
	arg.SetState(true)
	// s.idleCounter.Down()
}

func (s *Processor) end(arg run.Args) {
	arg.SetState(false)
	s.idleCounter.Up()
}

func (s *Processor) flush(arg run.Args, v ...any) {
	// s.begin(arg)
	if s.counter.Count() > 1 {
		SafeCall(s.mq.Exec, true, s.handler, v...)
	} else {
		SafeCall(s.mq.Exec, false, s.handler, v...)
	}
	// s.end(arg)
}

func (s *Processor) Count() int {
	return s.counter.Count()
}

func (s *Processor) IdleCount() int {
	return s.idleCounter.Count()
}

func (s *Processor) Wait() {
	s.counter.Wait()
}

// GC垃圾回收
func (s *Processor) Gc(args run.Args) (b bool) {
	if s.gcCondition == nil {
		return
	}
	if _, ok := s.gcCondition(s, args); ok {
		b = ok
	}
	return
}

// GC垃圾回收条件检查
func (s Processor) GcCondition(r run.Processor, args run.Args) (n int, b bool) {
	// 标记为啥mq等于nil ???
	if s.mq == nil {
		return
	}
	if q, ok := s.Queue().(ch.Queue); ok {
		n = q.Length() + q.Size()
		if s.IdleCount() > 0 {
			b = true
		}
	} else {
		n = s.Queue().Size()
		if s.IdleCount() > 0 {
			b = true
		}
	}
	return
}

func SafeCall(
	f func(bool, cb.Processor, ...any) (exit bool, code int),
	b bool,
	handler cb.Processor,
	args ...any) (err error) {
	defer run.Catch()
	f(b, handler, args...)
	return
}
