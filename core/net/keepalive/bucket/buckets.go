package bucket

import (
	"errors"
	"runtime"
	"time"

	"github.com/cwloo/gonet/core/base/pipe"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/core/base/run/timer_wheel"
	"github.com/cwloo/gonet/core/base/timer"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/core/net/conn"
)

// 定时轮盘池，处理空闲会话(多生产者，多消费者)
type Buckets interface {
	Size() int32
	Interval() int32
	Next() (pipe pipe.Pipe)
	Start()
	Stop()
}

type buckets struct {
	p Pool
}

func NewBuckets(bucketsz int32, d time.Duration) Buckets {
	if bucketsz <= 0 {
		return &buckets{}
	}
	cpu := runtime.NumCPU()
	//d := time.Second
	s := &buckets{
		p: NewPool("buckets.pool", bucketsz, d, int32(cpu)),
	}
	s.p.SetProcessor(cb.Processor(s.handler))
	s.p.SetTimerCallback(timer.TimerCallback(s.onTimer))
	s.p.Start()
	return s
}

func (s *buckets) Size() int32 {
	switch s.p {
	case nil:
		return 0
	default:
		return s.p.Bucketsz()
	}
}

func (s *buckets) Interval() int32 {
	switch s.p {
	case nil:
		return 0
	default:
		return int32(int64(s.p.D()) / int64(time.Second))
	}
}

func (s *buckets) onTimer(timerID uint32, dt int32, args ...any) bool {
	if len(args) == 0 {
		panic(errors.New("buckets.args 0"))
	}
	if args[0] == nil {
		panic(errors.New("buckets.args[0] is nil"))
	}
	switch arg := args[0].(type) {
	case *timer_wheel.Args:
		if arg.GetUsing() {
			// logs.Warnf("tick:%d dt:%d %v", s.Interval(), dt/1000, reflect.TypeOf(args).Name())
			msgs := arg.PopBucket(s.Interval())
			for _, msg := range msgs {
				if peer, ok := msg.(conn.Session); ok {
					peer.CloseExpired()
				}
			}
		}
	case run.Args:
	default:
	}
	return true
}

func (s *buckets) handler(msg any, args ...any) bool {
	if len(args) < 1 {
		panic(errors.New("args.size"))
	}
	proc, ok := args[0].(run.Proc)
	if !ok {
		panic(errors.New("arg[0]"))
	}
	arg, ok := proc.Args().(*timer_wheel.Args)
	if !ok {
		panic(errors.New(""))
	}
	switch data := msg.(type) {
	case Data:
		switch data.OpType() {
		case Push:
			arg.SetUsing(true)
			cursor := arg.PushBucket(data.Peer(), s.Size())
			data.Cb()(cursor)
		case Update:
			cursor := arg.UpdateBucket(data.Peer(), data.Cursor(), s.Size())
			data.Cb()(cursor)
		}
		data.Put()
	case timer.Data:
		switch data.OpType() {
		case timer.RunAfter:
			timerId := arg.RunAfter(data.Delay(), data.Args()...)
			data.Cb()(timerId)
		case timer.RunAfterWith:
			timerId := arg.RunAfterWith(data.Delay(), data.TimerCallback(), data.Args()...)
			data.Cb()(timerId)
		case timer.RunEvery:
			timerId := arg.RunEvery(data.Delay(), data.Interval(), data.Args()...)
			data.Cb()(timerId)
		case timer.RunEveryWith:
			timerId := arg.RunEveryWith(data.Delay(), data.Interval(), data.TimerCallback(), data.Args()...)
			data.Cb()(timerId)
		case timer.RemoveTimer:
			arg.RemoveTimer(data.TimerId())
		case timer.RemoveTimers:
			arg.RemoveTimers()
		}
		data.Put()
	}
	return false
}

func (s *buckets) Start() {
	switch s.p {
	case nil:
	default:
		s.p.Start()
	}
}

func (s *buckets) Stop() {
	switch s.p {
	case nil:
	default:
		s.p.Stop()
	}
}

func (s *buckets) Next() (pipe pipe.Pipe) {
	switch s.p {
	case nil:
	default:
		pipe = s.p.Next()
	}
	return
}
