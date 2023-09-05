package cc

import (
	"sync"
)

// 普通锁计数器
type counter struct {
	c    int32
	lock *sync.Mutex
	cond *sync.Cond
}

func NewCounter() Counter {
	s := &counter{
		lock: &sync.Mutex{},
	}
	s.cond = sync.NewCond(s.lock)
	return s
}

func (s *counter) Up() {
	s.lock.Lock()
	s.c++
	s.lock.Unlock()
}

func (s *counter) Down() {
	s.lock.Lock()
	s.c--
	if s.c == 0 {
		s.cond.Signal()
	}
	s.lock.Unlock()
}

func (s *counter) Wait() {
	s.lock.Lock()
	for s.c != 0 {
		s.cond.Wait()
	}
	s.lock.Unlock()
}

func (s *counter) Count() int {
	s.lock.Lock()
	c := int(s.c)
	s.lock.Unlock()
	return c
}

func (s *counter) Reset() {
	s.lock.Lock()
	s.c = 0
	s.lock.Unlock()
}
