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
	connectedPool = sync.Pool{
		New: func() any {
			return &Connected{}
		},
	}
	closingPool = sync.Pool{
		New: func() any {
			return &Closing{}
		},
	}
	closedPool = sync.Pool{
		New: func() any {
			return &Closed{}
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
)

type Type int8

const (
	EVTConnected Type = iota //建立连接事件
	EVTClosing               //通知关闭事件
	EVTClosed                //响应断开事件
	EVTRead                  //网络读取事件
	EVTCustom                //自定义事件
)

// 事件数据
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

// 建立连接事件
type Connected struct {
	Peer    conn.Session
	Handler cb.OnConnected
	Args    []any
}

func (s *Connected) Put() {
	connectedPool.Put(s)
}

func CreateConnected(peer conn.Session, v ...any) *Connected {
	s := connectedPool.Get().(*Connected)
	s.Peer = peer
	s.Args = append(s.Args, v...)
	return s
}

func CreateConnectedWith(handler cb.OnConnected, peer conn.Session, v ...any) *Connected {
	s := connectedPool.Get().(*Connected)
	s.Handler = handler
	s.Peer = peer
	s.Args = append(s.Args, v...)
	return s
}

// 通知关闭事件
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

// 响应断开事件
type Closed struct {
	Peer    conn.Session
	Reason  conn.Reason
	Handler cb.OnClosed
	Args    []any
}

func (s *Closed) Put() {
	closedPool.Put(s)
}

func CreateClosed(peer conn.Session, reason conn.Reason, v ...any) *Closed {
	s := closedPool.Get().(*Closed)
	s.Peer = peer
	s.Reason = reason
	s.Args = append(s.Args, v...)
	return s
}

func CreateClosedWith(handler cb.OnClosed, peer conn.Session, reason conn.Reason, v ...any) *Closed {
	s := closedPool.Get().(*Closed)
	s.Handler = handler
	s.Peer = peer
	s.Reason = reason
	s.Args = append(s.Args, v...)
	return s
}

// 网络读取事件
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

// 自定义事件
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
