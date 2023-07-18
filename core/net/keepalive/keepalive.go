package keepalive

import (
	"sync"

	"github.com/cwloo/gonet/core/base/pipe"
	"github.com/cwloo/gonet/core/net/conn"
	"github.com/cwloo/gonet/core/net/keepalive/bucket"
)

var (
	t = sync.Pool{
		New: func() any {
			return &keepalive{}
		},
	}
)

// <summary>
// Bucket 定时轮盘池，处理空闲会话(多生产者，多消费者)
// <summary>
type Buckets interface {
	Push(peer conn.Session)
	Update(peer conn.Session)
	Put()
}

type keepalive struct {
	cursor int32
	pipe   pipe.Pipe
}

func NewBuckets() Buckets {
	s := t.Get().(*keepalive)
	s.pipe = buckets.Next()
	switch s.pipe {
	case nil:
	default:
		s.pipe.AssociatedUp()
	}
	return s
}

func (s *keepalive) Push(peer conn.Session) {
	switch s.pipe {
	case nil:
	default:
		s.pipe.Do(
			bucket.NewPushBucket(
				peer,
				func(args ...any) {
					if len(args) > 0 {
						if cursor, ok := args[0].(int32); ok {
							if cursor > 0 {
								s.cursor = cursor
							}
						}
					}
				}))
	}
}

func (s *keepalive) Update(peer conn.Session) {
	switch s.pipe {
	case nil:
	default:
		s.pipe.Do(
			bucket.NewUpdateBucket(
				peer,
				s.cursor,
				func(args ...any) {
					if len(args) > 0 {
						if cursor, ok := args[0].(int32); ok {
							if cursor > 0 {
								s.cursor = cursor
							}
						}
					}
				}))
	}
}

func (s *keepalive) Put() {
	switch s.pipe {
	case nil:
	default:
		s.pipe.AssociatedDown()
	}
	t.Put(s)
}
