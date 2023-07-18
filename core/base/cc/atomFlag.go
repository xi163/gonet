package cc

import (
	"sync/atomic"
)

// <summary>
// AtomFlag 原子操作
// <summary>
type AtomFlag interface {
	Set()
	Reset()
	TestSet() bool
	TestReset() bool
	IsSet() bool
	IsReset() bool
}

type atomFlag struct {
	c int32
}

const (
	Reseted int32 = iota
	Seted
)

func NewAtomFlag() AtomFlag {
	s := &atomFlag{}
	s.Reset()
	return s
}

func (s *atomFlag) Set() {
	atomic.StoreInt32(&s.c, Seted)
}

func (s *atomFlag) Reset() {
	atomic.StoreInt32(&s.c, Reseted)
}

func (s *atomFlag) TestSet() bool {
	// return Reseted == atomic.SwapInt32(&s.c, Seted)
	return atomic.CompareAndSwapInt32(&s.c, Reseted, Seted)
}

func (s *atomFlag) TestReset() bool {
	return atomic.CompareAndSwapInt32(&s.c, Seted, Reseted)
}

func (s *atomFlag) IsSet() bool {
	return atomic.LoadInt32(&s.c) == Seted
}

func (s *atomFlag) IsReset() bool {
	return atomic.LoadInt32(&s.c) == Reseted
}
