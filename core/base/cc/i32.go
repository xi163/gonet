package cc

import "sync/atomic"

// <summary>
// I32 32位ID生成器
// <summary>
type I32 interface {
	New() int32
}

type i32 struct {
	id int32
}

func NewI32() I32 {
	return &i32{}
}

func (s *i32) New() int32 {
	return atomic.AddInt32(&s.id, 1)
}
