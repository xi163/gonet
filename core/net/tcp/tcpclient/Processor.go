package tcpclient

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
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

// TCP客户端
type Processor struct {
	name            string
	hold            conn.HoldType
	peers           conn.Sessions
	connector       tcp.Connector
	onConnected     cb.OnConnected
	onClosed        cb.OnClosed
	onMessage       cb.OnMessage
	onWriteComplete cb.OnWriteComplete
}

func NewTCPClient(name string, address ...string) TCPClient {
	s := &Processor{
		name:      name,
		hold:      conn.KHoldNone,
		peers:     conn.NewSessions(),
		connector: tcp.NewConnector(name, address...)}
	s.connector.SetDialTimeout(10 * time.Second)
	s.connector.SetIdleTimeout(30 * time.Second)
	s.connector.SetRetryInterval(time.Second)
	s.connector.SetProtocolCallback(s.onProtocol)
	s.connector.SetNewConnectionCallback(s.newConnection)
	s.connector.SetConnectErrorCallback(s.onConnectError)
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
		panic("error")
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

func (s *Processor) assertConnector() {
	if s.connector == nil {
		panic(errors.New("error"))
	}
}

func (s *Processor) Retry() bool {
	s.assertConnector()
	return s.connector.Retry()
}

func (s *Processor) EnableRetry(retry bool) {
	s.assertConnector()
	s.connector.EnableRetry(retry)
}

func (s *Processor) SetDialTimeout(d time.Duration) {
	s.assertConnector()
	s.connector.SetDialTimeout(d)
}

func (s *Processor) SetIdleTimeout(timeout, d time.Duration) {
	s.assertConnector()
	s.connector.SetIdleTimeout(d)
	keepalive.Init(timeout, d)
}

func (s *Processor) SetRetryInterval(d time.Duration) {
	s.assertConnector()
	s.connector.SetRetryInterval(d)
}

func (s *Processor) SetProtocolCallback(cb cb.OnProtocol) {
	s.assertConnector()
	s.connector.SetProtocolCallback(cb)
}

func (s *Processor) SetConnectErrorCallback(cb cb.OnConnectError) {
	s.assertConnector()
	s.connector.SetConnectErrorCallback(cb)
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

func (s *Processor) newConnection(c any, channel transmit.Channel, protoName string, v ...any) {
	switch protoName {
	case "tcp":
		if p, ok := c.(net.Conn); ok {
			connID := conn.NewConnID()
			localAddr := p.LocalAddr().String()
			peerAddr := p.RemoteAddr().String()
			peer := tcp.NewTCPConnection(
				connID,
				strings.Join([]string{s.name, "#", localAddr, "->", peerAddr, "#", strconv.FormatInt(connID, 10)}, ""),
				c,
				conn.KClient,
				channel, localAddr, peerAddr, protoName, s.connector.GetIdleTimeout())
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
				strings.Join([]string{s.name, "#", localAddr, "->", peerAddr, "#", strconv.FormatInt(connID, 10)}, ""),
				c,
				conn.KClient,
				channel, localAddr, peerAddr, protoName, s.connector.GetIdleTimeout())
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

func (s *Processor) ConnectTCP(header http.Header, address ...string) {
	s.assertConnector()
	s.connector.ConnectTCP(header, address...)
}

func (s *Processor) onConnectError(proto string, err error) {
}

func (s *Processor) OnConnected(peer conn.Session, v ...any) {
	if peer.Connected() {
		logs.Debugf("[%v] -> [%v]", peer.LocalAddr(), peer.RemoteAddr())
	} else {
		logs.Fatalf("error")
	}
}

func (s *Processor) OnClosed(peer conn.Session, reason conn.Reason, v ...any) {
	if peer.Connected() {
		logs.Fatalf("error")
	} else {
		logs.Tracef("[%v] -> [%v] %v", peer.LocalAddr(), peer.RemoteAddr(), reason.Msg)
	}
}

func (s *Processor) OnMessage(peer conn.Session, msg any, recvTime timestamp.T) {
	// logs.Infof("")
}

func (s *Processor) OnWriteComplete(peer conn.Session) {
	// logs.Debugf("")
}

func (s *Processor) removeConnection(peer conn.Session) {
	peer.(*tcp.TCPConnection).ConnectDestroyed()
	s.assertConnector()
	if s.connector.Retry() {
		s.connector.Reconnect()
	}
}

func (s *Processor) onConnectionError(err error) {
	// logs.Errorf("")
}

func (s *Processor) Reconnect() {
	s.assertConnector()
	s.connector.Reconnect()
}

func (s *Processor) Disconnect() {
	s.assertConnector()
	s.connector.EnableRetry(false)
	if conn.KHold == s.hold {
		s.peers.CloseAll()
	}
}
