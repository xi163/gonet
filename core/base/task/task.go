package task

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/mq"
	"github.com/cwloo/gonet/core/base/mq/ch"
	"github.com/cwloo/gonet/core/base/mq/lq"
	"github.com/cwloo/gonet/core/base/mq/sq"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/core/base/watcher"
	"github.com/cwloo/gonet/core/cb"
)

// 任务池(单生产者，多消费者)
type Task interface {
	Fixed() bool
	Nonblock() bool
	Queue() mq.Queue
	Runner() run.Processor
	Do(data any)
	DoTimeout(d time.Duration, data any, cb cb.Functor)
	Start()
	Stop()
	SetNew(handler mq.New)
	SetProcessor(handler cb.Processor)
}

type task struct {
	name     string
	init     int
	size     int
	fixed    bool
	nonblock bool
	i32      cc.I32
	lock     *sync.Mutex
	slots    map[int32]run.Slot
	mq       mq.Queue
	run      run.Processor
	flag     [2]cc.AtomFlag
	watcher  watcher.Watcher
	New      mq.New
}

func NewTask(name string, init, size int, fixed, nonblock bool, r run.Processor) Task {
	s := &task{
		name:     name,
		init:     init,
		size:     size,
		fixed:    fixed,
		nonblock: nonblock,
		i32:      cc.NewI32(),
		lock:     &sync.Mutex{},
		slots:    map[int32]run.Slot{},
		watcher:  watcher.NewWatcher(name, lq.NewQueue(0)),
		run:      r,
		flag: [2]cc.AtomFlag{
			cc.NewAtomFlag(),
			cc.NewAtomFlag()},
	}
	// ptr := utils.IF(New == nil, s.New_chmq, New)
	// unsafe.Pointer(&ptr)
	// s.New = utils.IF(New == nil, s.New_chmq, New).(mq.New)
	// s.overload = utils.IF(overload == nil, s.Overload, overload).(run.Overload)
	return s
}

func (s *task) assertQueue() {
	if s.mq == nil {
		panic(errors.New("error: task.mq is nil"))
	}
}

func (s *task) assertRunner() {
	if s.run == nil {
		panic(errors.New("error: task.run is nil"))
	}
}

func (s *task) SetNew(handler mq.New) {
	if handler == nil {
		panic(errors.New("error: task.SetNew is nil"))
	}
	s.New = handler
}

func (s *task) SetProcessor(handler cb.Processor) {
	s.assertRunner()
	s.run.SetProcessor(handler)
}

func (s *task) Queue() mq.Queue {
	s.assertQueue()
	return s.mq
}

func (s *task) Runner() run.Processor {
	s.assertRunner()
	return s.run
}

func (s *task) Do(data any) {
	if data != nil {
		s.do(data)
	}
}

func (s *task) DoTimeout(d time.Duration, data any, f cb.Functor) {
	if data != nil {
		s.do(cb.NewTimeout(time.Now(), d, data))
		// select {
		// case <-time.After(5 * time.Second):
		// 	break
		// }
		timeouts.After(d, f)
	}
}

func (s *task) do(data any) {
	s.watcher.Start(s.remove)
	s.start()
	s.ensure()
	s.assertQueue()
	s.mq.Push(data)
}

func (s *task) overload(r run.Processor) (n int, b bool) {
	if q, ok := r.Queue().(ch.Queue); ok {
		n = q.Length() + q.Size()
		if n == 0 {
			b = true
		}
	} else {
		n = r.Queue().Size()
		if n == 0 {
			b = true
		}
	}
	return
}

func (s *task) ensure() {
	if !s.Fixed() {
		if n, ok := s.overload(s.run); ok {
			s.expand(2 * n)
		}
	}
}

func (s task) Fixed() bool {
	return s.fixed
}

func (s task) Nonblock() bool {
	return s.nonblock
}

func (s *task) expand(num int) {
	for i := 0; i < num; i++ {
		id := s.i32.New()
		slot := s.new_slot(id)
		s.append(slot)
		slot.Sched(s.run)
	}
}

func (s *task) new_slot(id int32) run.Slot {
	q := s.mq.Name()
	r := s.run.Name()
	name := s.name + fmt.Sprintf(".%v.%v.slot.%v", q, r, id)
	return run.NewSlot(id, name, s.onQuit)
}

func (s *task) append(slot run.Slot) {
	s.lock.Lock()
	s.slots[slot.ID()] = slot
	s.lock.Unlock()
}

func (s *task) start() {
	if s.mq == nil && s.flag[0].TestSet() {
		if s.init > 0 {
			s.mq = s.New(s.size, s.nonblock)
			s.run.SetQueue(s.mq)
			s.expand(s.init)
		} else {
			s.mq = s.New(s.size, s.nonblock)
			s.run.SetQueue(s.mq)
			s.expand(1)
		}
		s.flag[0].Reset()
	}
}

func (s *task) stop() {
	if s.mq != nil && s.flag[1].TestSet() {
		s.mq.Push(nil)
		// for _, slot := range s.slots {
		// 	slot.Wait()
		// }
		s.run.Wait()
		s.slots = map[int32]run.Slot{}
		s.mq = nil
		s.flag[1].Reset()
	}
}

func (s *task) New_chmq(v ...any) (q mq.Queue) {
	if t, ok := ch.NewChan(v[0].(int), v[1].(bool)).(mq.Queue); ok {
		q = t
		return
	}
	panic(errors.New("task.New_chmq error"))
}

func (s *task) New_slicemq(v ...any) (q mq.Queue) {
	q = sq.NewQueue(v[0].(int))
	panic(errors.New("task.New_slicemq error"))
}

func (s *task) New_listmq(v ...any) (q mq.Queue) {
	q = lq.NewQueue(v[0].(int))
	panic(errors.New("task.New_listmq error"))
}

func (s *task) Start() {
	s.watcher.Start(s.remove)
	s.start()
}

func (s *task) remove(ids ...any) (exit bool) {
	s.lock.Lock()
	for _, data := range ids {
		if data == nil {
			exit = true
			break
		}
		if id, ok := data.(int32); ok {
			delete(s.slots, id)
			// logs.Debugf("slot.%d left:%d...", id, len(s.slots))
		}
	}
	s.lock.Unlock()
	return
}

func (s *task) Stop() {
	s.watcher.Stop()
	s.stop()
}

func (s *task) onQuit(slot run.Slot) {
	// logs.Debugf("slot.%d left:%d...", slot.ID(), len(s.slots))
	s.watcher.Push(slot.ID())
}
