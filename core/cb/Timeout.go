package cb

import (
	"sync"
	"time"

	"github.com/cwloo/gonet/core/base/cc"
)

var (
	t = sync.Pool{
		New: func() any {
			return &timeout{}
		},
	}
)

// 超时请求结构
type Timeout interface {
	Data() any
	Expire() cc.Expire
	Put()
}

type timeout struct {
	expire cc.Expire
	data   any
}

func NewTimeout(start time.Time, d time.Duration, data any) Timeout {
	s := t.Get().(*timeout)
	s.expire = cc.NewExpire(start, d)
	s.data = data
	return s
}

func (s *timeout) Data() any {
	return s.data
}

func (s *timeout) Expire() cc.Expire {
	return s.expire
}

func (s *timeout) Put() {
	t.Put(s)
}
