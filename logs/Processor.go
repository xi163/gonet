package logs

import (
	"errors"
	"runtime"
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/mq"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/core/cb"
)

// <summary>
// Processor 执行消息队列
// <summary>
type Processor struct {
	run.Processor
	mq          mq.BlockQueue
	counter     cc.Counter
	idleCounter cc.Counter
	handler     cb.Processor
	gcCondition run.GcCondition
}

func NewProcessor(handler cb.Processor) run.Processor {
	s := &Processor{
		handler:     handler,
		counter:     cc.NewAtomCounter(),
		idleCounter: cc.NewAtomCounter(),
	}
	return s
}

func NewProcessorWith(q mq.BlockQueue, handler cb.Processor) run.Processor {
	s := &Processor{
		mq:          q,
		handler:     handler,
		counter:     cc.NewAtomCounter(),
		idleCounter: cc.NewAtomCounter(),
	}
	return s
}

func (s *Processor) SetProcessor(handler cb.Processor) {
	s.handler = handler
}

func (s *Processor) SetGcCondition(handler run.GcCondition) {
	s.gcCondition = handler
}

func (s *Processor) Name() string {
	return "logs.Processor"
}

func (s *Processor) assertQueue() {
	if s.mq == nil {
		panic(errors.New("logs.Processor.mq is nil"))
	}
}

func (s *Processor) Queue() mq.Queue {
	s.assertQueue()
	return s.mq
}

func (s *Processor) SetQueue(q mq.Queue) {
	if mq, ok := q.(mq.BlockQueue); ok {
		s.mq = mq
	} else {
		panic(errors.New("need mq.BlockQueue"))
	}
}

func (s *Processor) NewArgs(proc run.Proc) run.Args {
	return newArgs(proc)
}

func (s *Processor) Run(proc run.Proc) {
	// logs.Debugf("%s started...", proc.Name())
	if s.mq == nil {
		panic(errors.New("error: logs.Processor.mq is nil"))
	}
	if s.handler == nil {
		panic(errors.New("error: logs.Processor.handler is nil"))
	}
	// if s.gcCondition == nil {
	// 	panic(errors.New("error: logs.Processor.gcCondition is nil"))
	// }
	if proc.Args() == nil {
		panic(errors.New("error: logs.Processor.args is nil"))
	}
	// arg := proc.Args().(*Args)
	tickerGC := run.NewTrigger(10 * time.Second)
	// s.counter.Up()
	// s.idleCounter.Up()
	flag := run.STOP
	i, t := 0, 200
EXIT:
	for {
		if i > t {
			i = 0
			runtime.GC()
			// runtime.Gosched()
		}
		i++
		exit, _ := s.mq.Exec(false, s.handler, proc)
		if exit {
			break EXIT
		}
	}
	tickerGC.Stop()
	// s.idleCounter.Down()
	// s.counter.Down()
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
	// s.counter.Wait()
}

func SafeCall(
	f func(bool, cb.Processor, ...any) (exit bool, code int),
	b bool,
	handler cb.Processor,
	args ...any) (err error) {
	defer Catch()
	f(b, handler, args...)
	return
}
