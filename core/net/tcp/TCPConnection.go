package tcp

import (
	"errors"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/mq"
	"github.com/cwloo/gonet/core/base/mq/lq"
	"github.com/cwloo/gonet/core/base/pool/gopool2"
	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/core/net/conn"
	"github.com/cwloo/gonet/core/net/keepalive"
	"github.com/cwloo/gonet/core/net/transmit"
	logs "github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils/safe"
	"github.com/cwloo/gonet/utils/timestamp"

	"github.com/gorilla/websocket"
)

var (
	connPool = sync.Pool{
		New: func() any {
			return &TCPConnection{}
		},
	}
)

// <summary>
// TCPConnection TCP连接会话
// <summary>
type TCPConnection struct {
	id                int64
	name              string
	localAddr         string
	remoteAddr        string
	protoName         string
	conn              any
	context           map[any]any
	connType          conn.Type
	mq                mq.BlockQueue
	channel           transmit.Channel
	wg                sync.WaitGroup
	closed            bool
	closing           cc.AtomFlag
	flag              cc.AtomFlag
	state             conn.State
	reason            conn.ReasonTD
	buckets           keepalive.Buckets
	onConnected       cb.OnConnected
	onClosed          cb.OnClosed
	onMessage         cb.OnMessage
	onWriteComplete   cb.OnWriteComplete
	closeCallback     cb.CloseCallback
	errorCallback     cb.ErrorCallback
	establishCallback func(v any)
	destroyCallback   func(v any)
}

func NewTCPConnection(id int64, name string, c any, connType conn.Type, channel transmit.Channel, localAddr, remoteAddr, protoName string, d time.Duration) conn.Session {
	peer := connPool.Get().(*TCPConnection)
	peer.id = id
	peer.name = name
	peer.conn = c
	peer.localAddr = localAddr
	peer.remoteAddr = remoteAddr
	peer.protoName = protoName
	peer.state = conn.KDisconnected
	peer.reason = conn.KNoError
	peer.connType = connType
	peer.context = map[any]any{}
	peer.mq = lq.NewQueue(0)
	peer.channel = channel
	peer.closing = cc.NewAtomFlag()
	peer.flag = cc.NewAtomFlag()
	peer.buckets = keepalive.NewBuckets()
	return peer
}

func (s *TCPConnection) Put() {
	connPool.Put(s)
}

func (s *TCPConnection) ID() int64 {
	return s.id
}

func (s *TCPConnection) Name() string {
	return s.name
}

func (s *TCPConnection) setState(state conn.State) {
	s.state = state
}

func (s *TCPConnection) setReason(reason conn.ReasonTD) {
	s.reason = reason
}

func (s *TCPConnection) Connected() bool {
	return s.state == conn.KConnected
}

func (s *TCPConnection) assertConn() {
	if s.conn == nil {
		panic(errors.New("error"))
	}
}

func (s *TCPConnection) assertChannel() {
	if s.channel == nil {
		panic(errors.New("error"))
	}
}

// func (s *TCPConnection) checkConn() {
// 	if s.conn == nil {
// 		return
// 	}
// 	if _, ok := s.conn.(net.Conn); ok {
// 	} else if _, ok := s.conn.(*websocket.Conn); ok {
// 	} else {
// 		panic(errors.New("error"))
// 	}
// }

func (s *TCPConnection) LocalAddr() string {
	return s.localAddr
}

func (s *TCPConnection) RemoteAddr() string {
	return s.remoteAddr
}

func (s *TCPConnection) ProtoName() string {
	return s.protoName
}

func (s *TCPConnection) Type() conn.Type {
	return s.connType
}

func (s *TCPConnection) SetContext(key any, val any) (old any) {
	if val != nil {
		if val, ok := s.context[key]; ok {
			old = val
		}
		s.context[key] = val
	} else if val, ok := s.context[key]; ok {
		old = val
		delete(s.context, key)
	}
	return
}

func (s *TCPConnection) GetContext(key any) any {
	if val, ok := s.context[key]; ok {
		return val
	}
	return nil
}

func (s *TCPConnection) SetConnectedCallback(cb cb.OnConnected) {
	s.onConnected = cb
}

func (s *TCPConnection) SetClosedCallback(cb cb.OnClosed) {
	s.onClosed = cb
}

func (s *TCPConnection) SetMessageCallback(cb cb.OnMessage) {
	s.onMessage = cb
}

func (s *TCPConnection) SetWriteCompleteCallback(cb cb.OnWriteComplete) {
	s.onWriteComplete = cb
}

func (s *TCPConnection) SetCloseCallback(cb cb.CloseCallback) {
	s.closeCallback = cb
}

func (s *TCPConnection) SetErrorCallback(cb cb.ErrorCallback) {
	s.errorCallback = cb
}

func (s *TCPConnection) SetEstablishCallback(cb func(v any)) {
	s.establishCallback = cb
}

func (s *TCPConnection) SetDestroyCallback(cb func(v any)) {
	s.destroyCallback = cb
}

func (s *TCPConnection) ConnectEstablished(v ...any) {
	s.wg.Add(1)
	s.SetContext("ext", v)
	// go s.readLoop()
	// go s.writeLoop()
	gopool2.Go(s.readLoop)
	gopool2.Go(s.writeLoop)
}

func (s *TCPConnection) connectEstablished(v ...any) {
	if s.id == 0 {
		panic(errors.New("error"))
	}
	s.setState(conn.KConnected)
	s.buckets.Push(s)
	if s.onConnected != nil {
		s.onConnected(s, v...)
	}
	if s.establishCallback != nil {
		s.establishCallback(s)
	}
}

func (s *TCPConnection) ConnectDestroyed() {
	if s.id == 0 {
		panic(errors.New("error"))
	}
	s.setState(conn.KDisconnected)
	s.buckets.Put()
	if s.onClosed != nil {
		s.onClosed(s, conn.Reasons[s.reason])
	}
	if s.destroyCallback != nil {
		s.destroyCallback(s)
	}
}

// 读协程
// onConnected/onMessage/onClosed三个回调同属一个协程内执行
func (s *TCPConnection) readLoop() {
	defer safe.Catch()
	s.assertChannel()
	s.assertConn()
	s.connectEstablished(s.GetContext("ext").([]any)...)
	s.SetContext("ext", nil)
	i, t := 0, 200
LOOP:
	for {
		if i > t {
			i = 0
			runtime.GC()
			// runtime.Gosched()
		}
		i++
		msg, err := s.channel.OnRecv(s.conn)
		if err != nil {
			// logs.Errorf("%v", err)
			// if !IsEOFOrReadError(err) {
			// 	if s.errorCallback != nil {
			// 		s.errorCallback(err)
			// 	}
			// }
			switch s.reason {
			case conn.KSelfClosed:
				// logs.Infof("self closed connection.")
				break LOOP
			case conn.KSelfClosedExpired:
				// logs.Infof("self closed expired connection.")
				break LOOP
			case conn.KSelfClosedDelay:
				// logs.Infof("self closed connection delay.")
				break LOOP
			default:
				// logs.Infof("peer closed connection.")
				s.setReason(conn.KPeerClosed)
				s.mq.Push(&mq.ExitStruct{Code: int(conn.KPeerClosed)})
				if s.errorCallback != nil {
					s.errorCallback(err)
				}
				break LOOP
			}
		} else if msg == nil {
			panic(errors.New("error"))
		} else if s.onMessage != nil {
			s.buckets.Update(s)
			s.onMessage(s, msg, timestamp.Now())
		} else {
			panic(errors.New("error"))
		}
	}
	// 关闭执行流程
	// TCPConnection.Close/CloseAfter ->
	// TCPConnection.close ->
	// TCPConnection.closeCallback ->
	// TCPServer.removeConnection ->
	// TCPConnection.ConnectDestroyed ->
	// TCPServer.onClosed ->
	// Sessions.Remove()
	if s.closeCallback != nil {
		s.closeCallback(s)
	}
	s.wg.Wait()
	// logs.Infof("exit.")
	s.conn = nil
	s.Put()
}

// 写协程
// 先关闭写(Write)再关闭读(onConnected/onMessage/onClosed), onClosed里面写(Write)无效!
func (s *TCPConnection) writeLoop() {
	defer safe.Catch()
	s.assertChannel()
	s.assertConn()
	i, t := 0, 200
LOOP:
	for {
		if i > t {
			i = 0
			runtime.GC()
			// runtime.Gosched()
		}
		i++
		// select {
		// default:
		exit, code := s.mq.Exec(false, func(msg any, args ...any) (exit bool) {
			// if ctx := s.GetContext("ctx").(user_context.Ctx); ctx != nil {
			// 	logs.Debugf("[%v:%v] write =>", ctx.GetUserId(), ctx.GetSession())
			// }
			err := s.channel.OnSend(s.conn, msg)
			if err != nil {
				logs.Errorf("%v", err)
				// if !transmit.IsEOFOrWriteError(err) {
				// 	if s.errorCallback != nil {
				// 		s.errorCallback(err)
				// 	}
				// }
			} else if s.onWriteComplete != nil {
				s.onWriteComplete(s)
			}
			return
		})
		if exit {
			if s.reason == conn.KPeerClosed {
				// logs.Infof("peer closed connection.")
			} else {
				switch code {
				case int(conn.KSelfClosed):
					// logs.Infof("self closed connection.")
					s.setReason(conn.KSelfClosed)
				case int(conn.KSelfClosedExpired):
					// logs.Infof("self closed expired connection.")
					s.setReason(conn.KSelfClosedExpired)
				case int(conn.KSelfClosedDelay):
					// logs.Infof("self closed connection delay.")
					s.setReason(conn.KSelfClosedDelay)
				default:
					panic("error")
				}
			}
			break LOOP
		}
		// }
	}
	s.close()
	// logs.Infof("exit.")
	s.wg.Done()
}

// 写数据
func (s *TCPConnection) Write(msg any) {
	if msg == nil {
		return
	}
	if s.Connected() {
		s.mq.Push(msg)
	}
}

// 关闭连接
func (s *TCPConnection) Close() {
	if s.conn == nil {
		return
	}
	if !s.closed && s.flag.TestSet() {
		s.notifyClose(int(conn.KSelfClosed))
	}
}

// 过期关闭对端
func (s *TCPConnection) CloseExpired() {
	if s.conn == nil {
		return
	}
	if !s.closed && s.flag.TestSet() {
		s.notifyClose(int(conn.KSelfClosedExpired))
	}
}

// 延迟关闭连接
func (s *TCPConnection) CloseAfter(d time.Duration) {
	if s.conn == nil {
		return
	}
	if !s.closed && s.flag.TestSet() {
		task.After(d, cb.NewFunctor00(func() {
			s.notifyClose(int(conn.KSelfClosedDelay))
		}))
	}
}

// 通知关闭连接
func (s *TCPConnection) notifyClose(flag int) {
	switch flag {
	case int(conn.KSelfClosedDelay):
		// logs.Infof("delay close.")
		s.mq.Push(&mq.ExitStruct{Code: int(conn.KSelfClosedDelay)})
	case int(conn.KSelfClosedExpired):
		// logs.Infof("close expired.")
		s.mq.Push(&mq.ExitStruct{Code: int(conn.KSelfClosedExpired)})
	case int(conn.KSelfClosed):
		s.mq.Push(&mq.ExitStruct{Code: int(conn.KSelfClosed)})
	default:
		panic("error")
	}
}

// 关闭执行流程
// TCPConnection.Close/CloseAfter ->
// TCPConnection.close ->
// TCPConnection.closeCallback ->
// TCPServer.removeConnection ->
// TCPConnection.ConnectDestroyed ->
// TCPServer.onClosed ->
// Sessions.Remove()
func (s *TCPConnection) close() {
	if s.conn == nil {
		return
	}
	// logs.Debugf("%v", s.name)
	if !s.closed && s.closing.TestSet() {
		if c, ok := s.conn.(net.Conn); ok {
			err := c.Close()
			if err != nil {
				logs.Errorf("%v", err)
			}
		} else if c, ok := s.conn.(*websocket.Conn); ok {
			err := c.Close()
			if err != nil {
				logs.Errorf("%v", err)
			}
		} else {
			panic(errors.New("error"))
		}
		s.closed = true
		s.closing.Reset()
		s.flag.Reset()
	}
}
