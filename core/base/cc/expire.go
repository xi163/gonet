package cc

import (
	"sync"
	"time"
)

var (
	e = sync.Pool{
		New: func() any {
			return &expire{}
		},
	}
)

// 过期判断结构
type Expire interface {
	StartTime() time.Time
	Duration() time.Duration
	Before(time time.Time) bool
	After(time time.Time) bool
	Expired(time time.Time) bool
	Put()
}

type expire struct {
	start time.Time
	d     time.Duration
}

func NewExpire(start time.Time, d time.Duration) Expire {
	s := e.Get().(*expire)
	s.start = start
	s.d = d
	return s
}

func (s *expire) StartTime() time.Time {
	return s.start
}

func (s *expire) Duration() time.Duration {
	return s.Duration()
}

func (s *expire) Before(time time.Time) bool {
	return s.start.UTC().Add(s.d).Before(time.UTC())
}

func (s *expire) After(time time.Time) bool {
	return s.start.UTC().Add(s.d).After(time.UTC())
}

func (s *expire) Expired(time time.Time) bool {
	return s.After(time)
}

func (s *expire) Put() {
	e.Put(s)
}
