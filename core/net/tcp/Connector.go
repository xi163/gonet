package tcp

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/cwloo/gonet/core/base/pool/connpool"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/core/net/conn"
	"github.com/cwloo/gonet/core/net/transmit"
	logs "github.com/cwloo/gonet/logs"

	"github.com/gorilla/websocket"
)

// <summary>
// Connector TCP连接器
// <summary>
type Connector interface {
	Reconnect()
	Retry() bool
	EnableRetry(bool)
	ServerAddr() string
	ConnectTCP(address ...string)
	SetProtocolCallback(cb cb.OnProtocol)
	SetNewConnectionCallback(cb cb.OnNewConnection)
	SetConnectErrorCallback(cb cb.OnConnectError)
	GetIdleTimeout() time.Duration
	SetIdleTimeout(d time.Duration)
	SetDialTimeout(d time.Duration)
	SetRetryInterval(d time.Duration)
}

// <summary>
// Connector TCP连接器
// <summary>
type connector struct {
	name, tmp       string
	retry           bool
	addr            *conn.Address
	dialTimeout     time.Duration
	d               time.Duration
	idleTimeout     time.Duration
	channel         transmit.Channel
	onProtocol      cb.OnProtocol
	onNewConnection cb.OnNewConnection
	onConnectError  cb.OnConnectError
}

func NewConnector(name string, address ...string) Connector {
	s := &connector{
		tmp:         name,
		dialTimeout: 10 * time.Second,
		idleTimeout: 30 * time.Second,
		d:           time.Second}
	if len(address) > 0 {
		s.addr = conn.ParseAddress(address[0])
	}
	return s
}

func (s *connector) SetDialTimeout(d time.Duration) {
	s.dialTimeout = d
}

func (s *connector) GetIdleTimeout() time.Duration {
	return s.idleTimeout
}

func (s *connector) SetIdleTimeout(d time.Duration) {
	s.idleTimeout = d
}

func (s *connector) SetRetryInterval(d time.Duration) {
	s.d = d
}

func (s *connector) toName() {
	if s.name == "" {
		s.name = s.tmp + "#" + s.addr.Format() + ".connector"
	}
}

func (s *connector) Retry() bool {
	return s.retry
}

func (s *connector) EnableRetry(retry bool) {
	s.retry = retry
}

func (s *connector) ServerAddr() string {
	return s.addr.Addr
}

func (s *connector) SetProtocolCallback(cb cb.OnProtocol) {
	s.onProtocol = cb
}

func (s *connector) SetNewConnectionCallback(cb cb.OnNewConnection) {
	s.onNewConnection = cb
}

func (s *connector) SetConnectErrorCallback(cb cb.OnConnectError) {
	s.onConnectError = cb
}

func (s *connector) assertProtocol() {
	if s.onProtocol == nil {
		panic(errors.New("error"))
	}
}

func (s *connector) assertOnNewConnection() {
	if s.onNewConnection == nil {
		panic(errors.New("error"))
	}
}

func (s *connector) connectTCPTimeout(addr *conn.Address, d time.Duration) error {
	logs.Debugf("%s", addr.Format())
	c, err := net.DialTimeout(addr.Proto, addr.Addr, d)
	if err != nil {
		// logs.Errorf(err.Error())
		s.onConnectError(addr.Proto, err)
		return err
	}
	switch conn.UsePool {
	case true:
		connpool.Do(cb.NewFunctor00(func() {
			s.onNewConnection(c, s.channel, s.addr.Proto)
		}))
	default:
		s.onNewConnection(c, s.channel, s.addr.Proto)
	}
	return nil
}

func (s *connector) connectWSTimeout(addr *conn.Address, d time.Duration) error {
	dialer := websocket.Dialer{}
	dialer.Proxy = http.ProxyFromEnvironment
	dialer.HandshakeTimeout = d
	u := url.URL{Scheme: addr.Proto, Host: addr.Addr, Path: addr.Path}
	logs.Debugf("%s", addr.Format())
	c, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		// logs.Errorf(err.Error())
		s.onConnectError(addr.Proto, err)
		return err
	}
	switch conn.UsePool {
	case true:
		connpool.Do(cb.NewFunctor00(func() {
			s.onNewConnection(c, s.channel, s.addr.Proto)
		}))
	default:
		s.onNewConnection(c, s.channel, s.addr.Proto)
	}
	return nil
}

func (s *connector) ConnectTCP(address ...string) {
	s.assertProtocol()
	s.assertOnNewConnection()
	if len(address) > 0 {
		s.addr = conn.ParseAddress(address[0])
	}
	if s.addr != nil {
		s.toName()
		s.channel = s.onProtocol(s.addr.Proto)
		switch s.addr.Proto {
		case "ws", "wss":
			if s.connectWSTimeout(s.addr, s.dialTimeout) != nil && s.retry {
				time.AfterFunc(s.d, s.reconnect)
			}
		case "tcp":
			if s.connectTCPTimeout(s.addr, s.dialTimeout) != nil && s.retry {
				time.AfterFunc(s.d, s.reconnect)
			}
		}
	}
}

func (s *connector) reconnect() {
	// logs.Debugf("%v %v", s.name, s.addr.Addr)
	switch s.addr.Proto {
	case "tcp":
		if s.connectTCPTimeout(s.addr, s.dialTimeout) != nil && s.retry {
			time.AfterFunc(s.d, s.reconnect)
		}
	case "ws", "wss":
		if s.connectWSTimeout(s.addr, s.dialTimeout) != nil && s.retry {
			time.AfterFunc(s.d, s.reconnect)
		}
	}
}

func (s *connector) Reconnect() {
	time.AfterFunc(s.d, s.reconnect)
}
