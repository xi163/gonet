package mailbox

import (
	"errors"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/mq"
	"github.com/cwloo/gonet/core/base/mq/ch"
	"github.com/cwloo/gonet/core/base/pipe"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/core/base/run/cell"
	"github.com/cwloo/gonet/core/base/run/event"
	"github.com/cwloo/gonet/core/base/run/workers"
	"github.com/cwloo/gonet/core/base/timer"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/utils/safe"
)

// <summary>
// Pipes 基于pipe的邮槽管理器接口/管道池(多生产者，多消费者)
// <summary>
type Pipes interface {
	Add(d time.Duration, creator cell.WorkerCreator, size, num int)
	AddOne(d time.Duration, creator cell.WorkerCreator, size int) pipe.Pipe
	Range(cb func(pipe.Pipe, int))
	Next() (pipe pipe.Pipe)
	Start()
	Wait()
	Stop()
	Num() int
	ResetNum()
}

type pipes struct {
	name   string
	i32    cc.I32
	slice  []pipe.Pipe
	signal cc.SysSignal
	next   int32
	c      cc.Counter
}

func NewPipes(name string) Pipes {
	s := &pipes{
		name:   name,
		i32:    cc.NewI32(),
		c:      cc.NewAtomCounter(),
		signal: cc.NewSysSignal(),
	}
	return s
}

func (s *pipes) Add(d time.Duration, creator cell.WorkerCreator, size, num int) {
	for i := 0; i < num; i++ {
		id := s.i32.New()
		pipe := s.new_pipe(id, d, creator)
		s.append(pipe)
	}
}

func (s *pipes) AddOne(d time.Duration, creator cell.WorkerCreator, size int) pipe.Pipe {
	id := s.i32.New()
	pipe := s.new_pipe(id, d, creator)
	s.append(pipe)
	return pipe
}

func (s *pipes) new_pipe(id int32, d time.Duration, creator cell.WorkerCreator) pipe.Pipe {
	cpu := runtime.NumCPU()
	cpu = 1
	nonblock := true //非阻塞
	tick := false    //开启tick检查
	// d := time.Second //tick间隔时间
	runner := workers.NewProcessor(tick, d, s.handler, s.onTimer, creator)
	// runner := workers.NewProcessor(d, s.handler, nil, creator)
	pipe := pipe.NewPipe(id, "worker.pipe", cpu, nonblock, runner)
	return pipe
}

func (s *pipes) New(v ...any) (q mq.Queue) {
	if t, ok := ch.NewChan(v[0].(int), v[1].(int), v[2].(bool)).(mq.Queue); ok {
		q = t
		return
	}
	panic(errors.New("new mq error"))
}

func (s *pipes) onTimer(timerID uint32, dt int32, args ...any) bool {
	if len(args) == 0 {
		panic(errors.New("pipes.args 0"))
	}
	if args[0] == nil {
		panic(errors.New("pipes.args[0] is nil"))
	}
	switch args[0].(type) {
	case cb.Functor:
		f, _ := args[0].(cb.Functor)
		f.Call()
		break
	}
	return true
}

func (s *pipes) handler(msg any, args ...any) bool {
	s.c.Up()
	if len(args) < 2 {
		panic(errors.New("args.size"))
	}
	proc, ok := args[0].(run.Proc)
	if !ok {
		panic(errors.New("arg[0]"))
	}
	arg, ok := proc.Args().(*workers.Args)
	if !ok {
		panic(errors.New(""))
	}
	worker, ok := args[1].(cell.Worker)
	if !ok {
		panic(errors.New("args[1]"))
	}
	proc.AssertThis()
	switch msg.(type) {
	case *event.Data:
		data, _ := msg.(*event.Data)
		proc.ResetDispatcher()
		switch data.Event {
		case event.EVTRead:
			// 网络读事件
			ev, _ := data.Object.(*event.Read)
			if ev.Handler != nil {
				ev.Handler(ev.Cmd, ev.Msg, ev.Peer)
			} else {
				worker.(cell.NetWorker).OnRead(ev.Cmd, ev.Msg, ev.Peer)
			}
			break
		case event.EVTCustom:
			// 自定义事件
			ev, _ := data.Object.(*event.Custom)
			if ev.Handler != nil {
				ev.Handler(ev.Cmd, ev.Msg, ev.Peer)
			} else {
				worker.(cell.NetWorker).OnCustom(ev.Cmd, ev.Msg, ev.Peer)
			}
			break
		case event.EVTClosing:
			// 通知关闭事件
			ev, _ := data.Object.(*event.Closing)
			if ev.D > 0 {
				ev.Peer.CloseAfter(ev.D)
			} else {
				ev.Peer.Close()
			}
			break
		}
		if proc.Dispatcher() != nil {
			proc.Dispatcher().Do(data)
		} else {
			s.recycle(data)
		}
		break
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
	// logger.Debugf("NumProcessed:%v", s.Num())
	return false
}

func (s *pipes) recycle(data *event.Data) {
	switch data.Event {
	case event.EVTRead:
		ev, _ := data.Object.(*event.Read)
		ev.Put()
		break
	case event.EVTCustom:
		ev, _ := data.Object.(*event.Custom)
		ev.Put()
		break
	case event.EVTClosing:
		ev, _ := data.Object.(*event.Closing)
		ev.Put()
		break
	}
	data.Put()
	// runtime.GC()
}

func (s *pipes) Range(cb func(pipe.Pipe, int)) {
	safe.Call(func() {
		for i, pipe := range s.slice {
			cb(pipe, i)
		}
	})
}

func (s *pipes) append(pipe pipe.Pipe) {
	s.slice = append(s.slice, pipe)
}

func (s *pipes) Next() (pipe pipe.Pipe) {
	if len(s.slice) > 0 {
		// 	pipe = s.slice[s.next]
		// 	s.next++
		// 	if s.next >= int32(len(s.slice)) {
		// 		s.next = 0
		// 	}
		c := atomic.AddInt32(&s.next, 1)
		if c >= int32(len(s.slice)) {
			atomic.StoreInt32(&s.next, 0)
		}
		mod := c % int32(len(s.slice))
		pipe = s.slice[mod]
	} else {
		panic(errors.New("pipes.pipes is empty"))
	}
	return
}

func (s *pipes) Start() {
	if len(s.slice) > 0 {
		s.signal.Start(s.clear)
	}
}

// 手动清理
func (s *pipes) clear() {
	for _, pipe := range s.slice {
		pipe.Close()
	}
	s.slice = []pipe.Pipe{}
}

// 等待退出
func (s *pipes) Wait() {
	s.signal.WaitSignal()
}

// 主动退出
func (s *pipes) Stop() {
	s.signal.Stop()
}

func (s *pipes) onQuit(pipe pipe.Pipe) {
}

func (s *pipes) Num() int {
	return s.c.Count()
}

func (s *pipes) ResetNum() {
	s.c.Reset()
}
