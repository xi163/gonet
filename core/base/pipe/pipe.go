package pipe

import (
	"errors"
	"fmt"
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/mq"
	"github.com/cwloo/gonet/core/base/mq/ch"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
)

// <summary>
// Pipe 管道(单生产者，单消费者)
// <summary>
type Pipe interface {
	ID() int32
	Name() string
	Slot() run.Slot
	Queue() mq.Queue
	Runner() run.Processor
	Do(data any)
	DoTimeout(d time.Duration, data any, cb cb.Functor)
	AssociatedUp()
	AssociatedDown()
	AssociatedCount() int
	Close()
	NotifyClose() bool
}

type pipe struct {
	slot run.Slot
	mq   mq.Queue
	run  run.Processor
	flag cc.AtomFlag
	c    cc.Counter
	cb   func(slot run.Slot)
}

func format(id int32, name string, q mq.Queue, r run.Processor) string {
	return name + fmt.Sprintf(".%v.%v.slot.%v", q.Name(), r.Name(), id)
}

func NewPipe(id int32, name string, size int, nonblock bool, r run.Processor) Pipe {
	s := &pipe{
		mq:   ch.NewChan(1, size, nonblock),
		run:  r,
		flag: cc.NewAtomFlag(),
		c:    cc.NewAtomCounter(),
	}
	s.assertRunner()
	s.run.SetQueue(s.mq)
	s.run.IdleUp() //空闲协程数量递增
	s.slot = run.NewSlot(id, format(id, name, s.mq, s.run), s.onQuit)
	s.slot.Sched(s.run)
	return s
}

func NewPipeWith(id int32, name string, q mq.Queue, r run.Processor) Pipe {
	s := &pipe{
		mq:   q,
		run:  r,
		flag: cc.NewAtomFlag(),
		c:    cc.NewAtomCounter(),
	}
	s.assertRunner()
	s.run.SetQueue(s.mq)
	s.run.IdleUp() //空闲协程数量递增
	s.slot = run.NewSlot(id, format(id, name, s.mq, s.run), s.onQuit)
	s.slot.Sched(s.run)
	return s
}

func NewPipeWithQuit(id int32, name string, q mq.Queue, r run.Processor, onQuit func(slot run.Slot)) Pipe {
	s := &pipe{
		mq:   q,
		run:  r,
		flag: cc.NewAtomFlag(),
		c:    cc.NewAtomCounter(),
		cb:   onQuit,
	}
	s.assertRunner()
	s.run.SetQueue(s.mq)
	s.run.IdleUp() //空闲协程数量递增
	s.slot = run.NewSlot(id, format(id, name, s.mq, s.run), s.onQuit)
	s.slot.Sched(s.run)
	return s
}

func (s *pipe) ID() int32 {
	s.assertSlot()
	return s.slot.ID()
}

func (s *pipe) Name() string {
	s.assertSlot()
	return s.slot.Name()
}

func (s *pipe) AssociatedUp() {
	s.c.Up()
}

func (s *pipe) AssociatedDown() {
	s.c.Down()
}

func (s *pipe) AssociatedCount() int {
	return s.c.Count()
}

func (s *pipe) assertSlot() {
	if s.slot == nil {
		panic(errors.New("error: pipe.slot is nil"))
	}
}

func (s *pipe) assertQueue() {
	if s.mq == nil {
		panic(errors.New("error: pipe.mq is nil"))
	}
}

func (s *pipe) assertRunner() {
	if s.run == nil {
		panic(errors.New("error: pipe.run is nil"))
	}
}

func (s *pipe) Slot() run.Slot {
	s.assertSlot()
	return s.slot
}

func (s *pipe) Queue() mq.Queue {
	s.assertQueue()
	return s.mq
}

func (s *pipe) Runner() run.Processor {
	s.assertRunner()
	return s.run
}

func (s *pipe) Do(data any) {
	if data != nil {
		s.do(data)
	}
}

func (s *pipe) DoTimeout(d time.Duration, data any, f cb.Functor) {
	if data != nil {
		s.do(cb.NewTimeout(time.Now(), d, data))
		// select {
		// case <-time.After(5 * time.Second):
		// 	break
		// }
		task.After(d, f)
	}
}

func (s *pipe) do(data any) {
	s.assertQueue()
	s.mq.Push(data)
}

func (s *pipe) Close() {
	if s.mq != nil && s.flag.TestSet() {
		s.mq.Push(nil)
		s.slot.Wait()
		s.flag.Reset()
	}
}

func (s *pipe) NotifyClose() (ok bool) {
	if s.mq != nil && s.flag.TestSet() {
		ok = true
		s.mq.Push(nil)
		s.flag.Reset()
	}
	return
}

func (s *pipe) onQuit(slot run.Slot) {
	if s.cb != nil {
		s.cb(slot)
	}
	s.slot = nil
	s.mq = nil
	s.run = nil
}
