package logs

import (
	"errors"
	"runtime"

	// "github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/mq"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/core/cb"
)

// 执行消息队列
type Processor struct {
	run.Processor
	mq      mq.BlockQueue
	handler cb.Processor
}

func NewProcessor(handler cb.Processor) run.Processor {
	return &Processor{
		handler: handler,
	}
}

func NewProcessorWith(q mq.BlockQueue, handler cb.Processor) run.Processor {
	return &Processor{
		mq:      q,
		handler: handler,
	}
}

func (s *Processor) SetProcessor(handler cb.Processor) {
	s.handler = handler
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
	if s.mq == nil {
		panic(errors.New("error: logs.Processor.mq is nil"))
	}
	if s.handler == nil {
		panic(errors.New("error: logs.Processor.handler is nil"))
	}
	if proc.Args() == nil {
		panic(errors.New("error: logs.Processor.args is nil"))
	}
	// arg := proc.Args().(*Args)
	// tickerGC := run.NewTrigger(10 * time.Second)
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
	// tickerGC.Stop()
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

func (s *Processor) Wait() {
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
