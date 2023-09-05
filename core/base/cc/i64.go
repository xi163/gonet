package cc

import "sync/atomic"

// 64位ID生成器
type I64 interface {
	New() int64
}

type i64 struct {
	id int64
}

func NewI64() I64 {
	return &i64{}
}

func (s *i64) New() int64 {
	return atomic.AddInt64(&s.id, 1)
}
