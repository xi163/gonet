package watcher

import (
	"sync"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/mq"
)

// 监视器、看门狗
type Watcher interface {
	Push(data any)
	Start(handler func(v ...any) (exit bool))
	Wait()
	Stop()
}

type watcher struct {
	name     string
	watching bool
	lock     *sync.Mutex
	cond     *sync.Cond
	mq       mq.BlockQueue
	flag     [2]cc.AtomFlag
}

func NewWatcher(name string, q mq.BlockQueue) Watcher {
	s := &watcher{
		name: name,
		lock: &sync.Mutex{},
		mq:   q,
		flag: [2]cc.AtomFlag{
			cc.NewAtomFlag(),
			cc.NewAtomFlag()},
	}
	s.cond = sync.NewCond(s.lock)
	return s
}

func (s *watcher) Push(data any) {
	s.mq.Push(data)
}

func (s *watcher) Start(handler func(v ...any) (exit bool)) {
	if !s.watching && s.flag[0].TestSet() {
		go s.watch(handler)
		s.wait()
		s.flag[0].Reset()
	}
}

func (s *watcher) wait() {
	s.lock.Lock()
	for !s.watching {
		s.cond.Wait()
	}
	s.lock.Unlock()
}

func (s *watcher) Wait() {
	s.lock.Lock()
	for s.watching {
		s.cond.Wait()
	}
	s.lock.Unlock()
}

func (s *watcher) Stop() {
	if s.watching && s.flag[1].TestSet() {
		s.mq.Push(nil)
		s.Wait()
		s.flag[1].Reset()
	}
}

func (s *watcher) watch(handler func(v ...any) (exit bool)) {
	s.lock.Lock()
	s.watching = true
	s.cond.Signal()
	s.lock.Unlock()
	for {
		msgs := s.mq.Pick()
		exit := handler(msgs...)
		if exit {
			break
		}
	}
	s.lock.Lock()
	s.watching = false
	s.cond.Signal()
	s.lock.Unlock()
}
