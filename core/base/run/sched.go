package run

import (
	"errors"
	"sync"
	"time"

	"github.com/cwloo/gonet/core/base/cc"
)

// 邮槽
type Slot interface {
	ID() int32
	Name() string
	Dead() bool
	Duration() time.Duration
	Reset(d time.Duration)
	Add(args ...any)
	Sched(r Processor) Proc
	Proc() Proc
	Wait()
	Stop()
}

type slot struct {
	id     int32
	name   string
	proc   Proc
	lock   *sync.Mutex
	cond   *sync.Cond
	flag   [2]cc.AtomFlag
	onQuit func(slot Slot)
}

func NewSlot(id int32, name string, onQuit func(slot Slot)) Slot {
	if onQuit == nil {
		panic(errors.New("onQuit is nil"))
	}
	s := &slot{
		id:     id,
		name:   name,
		lock:   &sync.Mutex{},
		flag:   [2]cc.AtomFlag{cc.NewAtomFlag(), cc.NewAtomFlag()},
		onQuit: onQuit,
	}
	s.cond = sync.NewCond(s.lock)
	return s
}

func (s *slot) ID() int32 {
	return s.id
}

func (s *slot) Name() string {
	return s.name
}

func (s *slot) Dead() bool {
	return s.proc == nil
}

func (s *slot) assertProc() {
	if s.proc == nil {
		panic(errors.New("error: slot.proc is nil"))
	}
}

func (s *slot) Proc() Proc {
	s.assertProc()
	return s.proc
}

func (s *slot) Duration() time.Duration {
	return s.Proc().Args().Duration()
}

func (s *slot) Reset(d time.Duration) {
	s.Proc().Args().Reset(d)
}

func (s *slot) Add(args ...any) {
	s.Proc().Args().Add(args...)
}

func (s *slot) Sched(r Processor) Proc {
	if s.proc == nil && s.flag[0].TestSet() {
		go s.run(r)
		s.wait()
		s.flag[0].Reset()
	}
	return s.proc
}

func (s *slot) wait() {
	s.lock.Lock()
	for s.proc == nil {
		s.cond.Wait()
	}
	s.lock.Unlock()
}

func (s *slot) run(r Processor) {
	proc := NewProc(s.name, r)
	s.lock.Lock()
	s.proc = proc
	s.cond.Signal()
	s.lock.Unlock()
	s.proc.Run()
	s.onQuit(s)
	s.lock.Lock()
	s.proc = nil
	s.cond.Signal()
	s.lock.Unlock()
}

func (s *slot) wait_stop() {
	s.lock.Lock()
	for s.proc != nil {
		s.cond.Wait()
	}
	s.lock.Unlock()
}

func (s *slot) Wait() {
	s.wait_stop()
}

func (s *slot) Stop() {
	if s.proc != nil && s.flag[1].TestSet() {
		s.proc.Quit()
		s.wait_stop()
		s.flag[1].Reset()
	}
}
