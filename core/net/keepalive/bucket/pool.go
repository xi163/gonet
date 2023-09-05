package bucket

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/pipe"
	"github.com/cwloo/gonet/core/base/run/timer_wheel"
	"github.com/cwloo/gonet/core/base/timer"
	"github.com/cwloo/gonet/core/cb"
)

// 定时轮盘池，处理空闲会话(多生产者，多消费者)
type Pool interface {
	D() time.Duration
	Bucketsz() int32
	Next() (pipe pipe.Pipe)
	Start()
	Stop()
	SetProcessor(handler cb.Processor)
	SetTimerCallback(handler timer.TimerCallback)
}

type pool struct {
	name        string
	num, bucket int32
	d           time.Duration
	timerCb     timer.TimerCallback
	handler     cb.Processor
	i32         cc.I32
	slice       []pipe.Pipe
	next        int32
}

func NewPool(name string, bucket int32, d time.Duration, num int32) Pool {
	s := &pool{
		name:   name,
		bucket: bucket,
		d:      d,
		num:    num,
		i32:    cc.NewI32(),
	}
	return s
}

func (s *pool) D() time.Duration {
	return s.d
}

func (s *pool) Bucketsz() int32 {
	return s.bucket
}

func (s *pool) SetProcessor(handler cb.Processor) {
	s.handler = handler
}

func (s *pool) SetTimerCallback(handler timer.TimerCallback) {
	s.timerCb = handler
}

func (s *pool) assertHandler() {
	if s.handler == nil {
		panic(errors.New("error: pool.handler is nil"))
	}
}

func (s *pool) assertTimerCb() {
	if s.handler == nil {
		panic(errors.New("error: pool.timerCb is nil"))
	}
}

func (s *pool) append(pipe pipe.Pipe) {
	s.slice = append(s.slice, pipe)
}

func (s *pool) Next() (pipe pipe.Pipe) {
	if len(s.slice) > 0 {
		// 	pipe = s.slice[s.next]
		// 	s.next++
		// 	if s.next >= int32(len(s.slice)) {
		// 		s.next = 0
		// 	}
		c := atomic.AddInt32(&s.next, 1)
		if c >= int32(len(s.slice)) {
			atomic.StoreInt32(&s.next, -1)
			c = atomic.AddInt32(&s.next, 1)
		}
		pipe = s.slice[c]
	}
	return
}

func (s *pool) expand(size int32, d time.Duration, num int32) {
	for i := int32(0); i < num; i++ {
		id := s.i32.New()
		pipe := s.new_pipe(id, size, d)
		s.append(pipe)
	}
}

func (s *pool) new_pipe(id int32, size int32, d time.Duration) pipe.Pipe {
	s.assertHandler()
	s.assertTimerCb()
	nonblock := true
	tick := true
	// d := time.Second
	runner := timer_wheel.NewProcessor(size, tick, d, s.handler, s.timerCb)
	pipe := pipe.NewPipe(id, "bucket.pipe", 500, nonblock, runner)
	return pipe
}

func (s *pool) Start() {
	if s.bucket > 0 {
		s.expand(s.bucket, s.d, s.num)
	}
}

func (s *pool) Stop() {
	for _, pipe := range s.slice {
		pipe.Close()
	}
	s.slice = []pipe.Pipe{}
}

func (s *pool) onQuit(pipe pipe.Pipe) {
}
