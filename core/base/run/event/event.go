package event

import (
	"sync"
	"time"

	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/core/net/conn"
)

var (
	dataPool = sync.Pool{
		New: func() any {
			return &Data{}
		},
	}
	readPool = sync.Pool{
		New: func() any {
			return &Read{}
		},
	}
	customPool = sync.Pool{
		New: func() any {
			return &Custom{}
		},
	}
	closingPool = sync.Pool{
		New: func() any {
			return &Closing{}
		},
	}
)

type Type int8

const (
	EVTRead Type = iota
	EVTCustom
	EVTClosing
)

// <summary>
// Data 事件数据
// <summary>
type Data struct {
	Event  Type
	Object any
	Extra  any
}

func Create(ev Type, obj any, ext any) *Data {
	s := dataPool.Get().(*Data)
	s.Event = ev
	s.Object = obj
	s.Extra = ext
	return s
}

func (s *Data) Put() {
	dataPool.Put(s)
}

// <summary>
// Read 读事件数据
// <summary>
type Read struct {
	Cmd     uint32
	Peer    conn.Session
	Msg     any
	Handler cb.ReadCallback
}

func (s *Read) Put() {
	readPool.Put(s)
}

func CreateRead(cmd uint32, msg any, peer conn.Session) *Read {
	s := readPool.Get().(*Read)
	s.Cmd = cmd
	s.Msg = msg
	s.Peer = peer
	return s
}

func CreateReadWith(handler cb.ReadCallback, cmd uint32, msg any, peer conn.Session) *Read {
	s := readPool.Get().(*Read)
	s.Handler = handler
	s.Cmd = cmd
	s.Msg = msg
	s.Peer = peer
	return s
}

// <summary>
// Custom 自定义事件数据
// <summary>
type Custom struct {
	Cmd     uint32
	Peer    conn.Session
	Msg     any
	Handler cb.CustomCallback
}

func (s *Custom) Put() {
	customPool.Put(s)
}

func CreateCustom(cmd uint32, msg any, peer conn.Session) *Custom {
	s := customPool.Get().(*Custom)
	s.Cmd = cmd
	s.Msg = msg
	s.Peer = peer
	return s
}

func CreateCustomWith(handler cb.CustomCallback, cmd uint32, msg any, peer conn.Session) *Custom {
	s := customPool.Get().(*Custom)
	s.Handler = handler
	s.Cmd = cmd
	s.Msg = msg
	s.Peer = peer
	return s
}

// <summary>
// Closing 通知关闭事件数据
// <summary>
type Closing struct {
	D    time.Duration
	Peer conn.Session
}

func (s *Closing) Put() {
	closingPool.Put(s)
}

func CreateClosing(d time.Duration, peer conn.Session) *Closing {
	s := closingPool.Get().(*Closing)
	s.D = d
	s.Peer = peer
	return s
}
