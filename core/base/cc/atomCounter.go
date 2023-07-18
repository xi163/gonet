package cc

import (
	"sync"
	"sync/atomic"
)

// <summary>
// atomCounter 原子锁计数器
// <summary>
type atomCounter struct {
	c    int32
	lock *sync.Mutex
	cond *sync.Cond
}

func NewAtomCounter() Counter {
	s := &atomCounter{
		lock: &sync.Mutex{},
	}
	s.cond = sync.NewCond(s.lock)
	return s
}

func (s *atomCounter) Up() {
	atomic.AddInt32(&s.c, 1)
	// c := atomic.AddInt32(&s.c, 1)
	// logs.Debugf("%d", c)
}

func (s *atomCounter) Down() {
	c := atomic.AddInt32(&s.c, -1)
	if c == 0 {
		s.cond.Signal()
	}
	// logs.Debugf("%d", c)
}

func (s *atomCounter) Wait() {
	s.lock.Lock()
	for atomic.LoadInt32(&s.c) != 0 {
		s.cond.Wait()
	}
	s.lock.Unlock()
}

func (s *atomCounter) Count() int {
	return int(atomic.LoadInt32(&s.c))
}

func (s *atomCounter) Reset() {
	atomic.StoreInt32(&s.c, 0)
}
