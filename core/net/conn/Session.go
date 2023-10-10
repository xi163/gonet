package conn

import "time"

type ReasonID uint8

const (
	KNoError           ReasonID = ReasonID(0)
	KPeerClosed        ReasonID = ReasonID(1) //对端关闭连接
	KSelfClosed        ReasonID = ReasonID(2) //本端关闭连接
	KSelfClosedDelay   ReasonID = ReasonID(3) //本端延时关闭
	KSelfClosedExpired ReasonID = ReasonID(4) //过期关闭对端
)

type Reason struct {
	Id  ReasonID
	Msg string
}

var (
	UsePool            = true
	ENoError           = Reason{KNoError, "NoError"}
	EPeerClosed        = Reason{KPeerClosed, "peer closed connection"}
	ESelfClosed        = Reason{KSelfClosed, "self closed connection"}
	ESelfClosedDelay   = Reason{KSelfClosedDelay, "self closed connection delay"}
	ESelfClosedExpired = Reason{KSelfClosedExpired, "self closed expired connection"}
	Reasons            = []Reason{ENoError, EPeerClosed, ESelfClosed, ESelfClosedDelay, ESelfClosedExpired}
)

type State uint8

const (
	KDisconnected State = State(0)
	KConnected    State = State(1)
)

type Type uint8

const (
	KClient Type = Type(0)
	KServer Type = Type(1)
)

// 连接会话
type Session interface {
	ID() int64
	Name() string
	ProtoName() string
	Type() Type
	Connected() bool
	LocalAddr() string
	RemoteAddr() string
	RemoteRegion() Region
	SetContext(key any, val any) (old any)
	GetContext(key any) any
	SetContextLocker(key any, val any) (old any)
	GetContextLocker(key any) any
	Write(msg any)
	WriteText(msg any)
	Close()
	CloseAfter(d time.Duration)
	CloseExpired()
	Put()
}
