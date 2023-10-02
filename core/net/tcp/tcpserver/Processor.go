package tcpserver

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/core/net/conn"
	"github.com/cwloo/gonet/core/net/keepalive"
	"github.com/cwloo/gonet/core/net/tcp"
	"github.com/cwloo/gonet/core/net/transmit"
	"github.com/cwloo/gonet/core/net/transmit/tcpchannel"
	"github.com/cwloo/gonet/core/net/transmit/wschannel"
	logs "github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils/timestamp"
	"github.com/gorilla/websocket"
)

// TCP服务端
type Processor struct {
	name            string
	numConnected    int32
	hold            conn.HoldType
	peers           conn.Sessions
	acceptor        tcp.Acceptor
	onConnected     cb.OnConnected
	onClosed        cb.OnClosed
	onMessage       cb.OnMessage
	onWriteComplete cb.OnWriteComplete
}

func NewTCPServer(name string, address ...string) TCPServer {
	s := &Processor{
		name:     name,
		hold:     conn.KHoldNone,
		peers:    conn.NewSessions(),
		acceptor: tcp.NewAcceptor(name, address...)}
	s.acceptor.SetProtocolCallback(s.onProtocol)
	s.acceptor.SetConditionCallback(s.onCondition)
	s.acceptor.SetNewConnectionCallback(s.newConnection)
	s.SetConnectedCallback(s.OnConnected)
	s.SetClosedCallback(s.OnClosed)
	s.SetMessageCallback(s.OnMessage)
	s.SetWriteCompleteCallback(s.OnWriteComplete)
	return s
}

func (s *Processor) Name() string {
	return s.name
}

func (s *Processor) Peers() conn.Sessions {
	if s.peers == nil {
		panic(errors.New("error"))
	}
	return s.peers
}

func (s *Processor) remove(v any) {
	if conn.KHoldTemporary == s.hold {
		s.peers.Remove(v.(conn.Session))
	}
}

func (s *Processor) reset(v any) {
	if conn.KHold == s.hold {
		// peer will be dtcor !
		s.peers.Remove(v.(conn.Session))
	}
}

func (s *Processor) SetHoldType(hold conn.HoldType) {
	s.hold = hold
}

func (s *Processor) Range(cb func(peer conn.Session)) {
	if conn.KHold == s.hold {
		s.peers.Range(cb)
	}
}

func (s *Processor) assertAcceptor() {
	if s.acceptor == nil {
		panic(errors.New("error"))
	}
}

func (s *Processor) SetHandshakeTimeout(d time.Duration) {
	s.assertAcceptor()
	s.acceptor.SetHandshakeTimeout(d)
}

func (s *Processor) SetIdleTimeout(timeout, d time.Duration) {
	s.assertAcceptor()
	s.acceptor.SetIdleTimeout(d)
	keepalive.Init(timeout, d)
}

func (s *Processor) SetReadBufferSize(readBufferSize int) {
	s.assertAcceptor()
	s.acceptor.SetReadBufferSize(readBufferSize)
}

func (s *Processor) ListenAddr() *conn.Address {
	s.assertAcceptor()
	return s.acceptor.Addr()
}

func (s *Processor) SetCertFile(certfile, keyfile string) {
	s.assertAcceptor()
	s.acceptor.SetCertFile(certfile, keyfile)
}

func (s *Processor) SetProtocolCallback(cb cb.OnProtocol) {
	s.assertAcceptor()
	s.acceptor.SetProtocolCallback(cb)
}

func (s *Processor) SetHandshakeCallback(cb cb.OnHandshake) {
	s.assertAcceptor()
	s.acceptor.SetHandshakeCallback(cb)
}

func (s *Processor) SetConditionCallback(cb cb.OnCondition) {
	s.assertAcceptor()
	s.acceptor.SetConditionCallback(cb)
}

func (s *Processor) SetConnectedCallback(cb cb.OnConnected) {
	s.onConnected = cb
}

func (s *Processor) SetClosedCallback(cb cb.OnClosed) {
	s.onClosed = cb
}

func (s *Processor) SetMessageCallback(cb cb.OnMessage) {
	s.onMessage = cb
}

func (s *Processor) SetWriteCompleteCallback(cb cb.OnWriteComplete) {
	s.onWriteComplete = cb
}

func (s *Processor) ListenTCP(address ...string) {
	s.acceptor.ListenTCP(address...)
}

func (s *Processor) Stop() {
	s.acceptor.Stop()
	if conn.KHold == s.hold {
		s.peers.CloseAll()
	}
}

func (s *Processor) onCondition(addr net.Addr) bool {
	return true
}

func (s *Processor) newConnection(c any, channel transmit.Channel, protoName string, v ...any) {
	switch protoName {
	case "tcp":
		if p, ok := c.(net.Conn); ok {
			connID := conn.NewConnID()
			localAddr := p.LocalAddr().String()
			peerAddr := p.RemoteAddr().String()
			peer := tcp.NewTCPConnection(
				connID,
				strings.Join([]string{s.name, "#", localAddr, "<-", peerAddr, "#", strconv.FormatInt(connID, 10)}, ""),
				c,
				conn.KServer,
				channel, localAddr, peerAddr, protoName, s.acceptor.GetIdleTimeout())
			peer.(*tcp.TCPConnection).SetConnectedCallback(s.onConnected)
			peer.(*tcp.TCPConnection).SetClosedCallback(s.onClosed)
			peer.(*tcp.TCPConnection).SetMessageCallback(s.onMessage)
			peer.(*tcp.TCPConnection).SetWriteCompleteCallback(s.onWriteComplete)
			peer.(*tcp.TCPConnection).SetCloseCallback(s.removeConnection)
			peer.(*tcp.TCPConnection).SetErrorCallback(s.onConnectionError)
			peer.(*tcp.TCPConnection).SetEstablishCallback(s.remove)
			peer.(*tcp.TCPConnection).SetDestroyCallback(s.reset)
			peer.(*tcp.TCPConnection).ConnectEstablished(v...)
			// save peer first, otherwise it will be dtcor immediately
			if conn.KHoldNone != s.hold && !s.peers.Add(peer) {
				peer.Close()
			}
		} else {
			logs.Fatalf("error")
		}
	case "ws", "wss":
		if p, ok := c.(*websocket.Conn); ok {
			connID := conn.NewConnID()
			localAddr := p.LocalAddr().String()
			peerAddr := p.RemoteAddr().String()
			peer := tcp.NewTCPConnection(
				connID,
				strings.Join([]string{s.name, "#", localAddr, "<-", peerAddr, "#", strconv.FormatInt(connID, 10)}, ""),
				c,
				conn.KServer,
				channel, localAddr, peerAddr, protoName, s.acceptor.GetIdleTimeout())
			peer.(*tcp.TCPConnection).SetConnectedCallback(s.onConnected)
			peer.(*tcp.TCPConnection).SetClosedCallback(s.onClosed)
			peer.(*tcp.TCPConnection).SetMessageCallback(s.onMessage)
			peer.(*tcp.TCPConnection).SetWriteCompleteCallback(s.onWriteComplete)
			peer.(*tcp.TCPConnection).SetCloseCallback(s.removeConnection)
			peer.(*tcp.TCPConnection).SetErrorCallback(s.onConnectionError)
			peer.(*tcp.TCPConnection).SetEstablishCallback(s.remove)
			peer.(*tcp.TCPConnection).SetDestroyCallback(s.reset)
			peer.(*tcp.TCPConnection).ConnectEstablished(v...)
			// save peer first, otherwise it will be dtcor immediately
			if conn.KHoldNone != s.hold && !s.peers.Add(peer) {
				peer.Close()
			}
		} else {
			logs.Fatalf("error")
		}
	default:
		logs.Fatalf("error")
	}
}

func (s *Processor) onProtocol(proto string) transmit.Channel {
	switch proto {
	case "tcp":
		return tcpchannel.NewChannel()
	case "ws", "wss":
		return wschannel.NewChannel()
	}
	panic("no proto setup")
}

func (s *Processor) OnConnected(peer conn.Session, v ...any) {
	if peer.Connected() {
		num := atomic.AddInt32(&s.numConnected, 1)
		logs.Debugf("%d [%v] <- [%v]", num, peer.LocalAddr(), peer.RemoteAddr())
	} else {
		logs.Fatalf("error")
	}
}

func (s *Processor) OnClosed(peer conn.Session, reason conn.Reason, v ...any) {
	if peer.Connected() {
		logs.Fatalf("error")
	} else {
		num := atomic.AddInt32(&s.numConnected, -1)
		logs.Tracef("%d [%v] <- [%v] %v", num, peer.LocalAddr(), peer.RemoteAddr(), reason.Msg)
	}
}

func (s *Processor) OnMessage(peer conn.Session, msg any, msgType int, recvTime timestamp.T) {
	// logs.Infof("")
}

func (s *Processor) OnWriteComplete(peer conn.Session) {
	// logs.Debugf("")
}

func (s *Processor) removeConnection(peer conn.Session) {
	// s.peers.Remove(peer)
	peer.(*tcp.TCPConnection).ConnectDestroyed()
}

func (s *Processor) onConnectionError(err error) {
	// logs.Errorf("")
}
