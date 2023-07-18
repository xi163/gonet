package bucket

import (
	"sync"

	"github.com/cwloo/gonet/core/net/conn"
)

var (
	t = sync.Pool{
		New: func() any {
			return &data{}
		},
	}
)

type OpType uint8

const (
	Push OpType = iota + 10
	Update
)

// <summary>
// Data 定时轮盘桶节点，处理空闲会话
// <summary>
type Data interface {
	OpType() OpType
	Cb() func(args ...any)
	Peer() conn.Session
	Cursor() int32
	Put()
}

type data struct {
	op     OpType
	cursor int32
	peer   conn.Session
	cb     func(args ...any)
}

func NewPushBucket(peer conn.Session, cb func(args ...any)) Data {
	s := t.Get().(*data)
	s.op = Push
	s.peer = peer
	s.cb = cb
	return s
}

func NewUpdateBucket(peer conn.Session, cursor int32, cb func(args ...any)) Data {
	s := t.Get().(*data)
	s.op = Update
	s.cursor = cursor
	s.peer = peer
	s.cb = cb
	return s
}

func (s *data) Peer() conn.Session {
	return s.peer
}

func (s *data) Cursor() int32 {
	return s.cursor
}

func (s *data) Cb() func(args ...any) {
	return s.cb
}

func (s *data) OpType() OpType {
	return s.op
}

func (s *data) Put() {
	t.Put(s)
}
