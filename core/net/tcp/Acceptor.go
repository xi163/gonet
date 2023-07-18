package tcp

import (
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/pool/connpool"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/core/net/conn"
	"github.com/cwloo/gonet/core/net/transmit"
	logs "github.com/cwloo/gonet/logs"

	"github.com/gorilla/websocket"
)

// <summary>
// Acceptor TCP接受器
// <summary>
type Acceptor interface {
	Addr() *conn.Address
	ListenTCP(address ...string)
	Stop()
	GetIdleTimeout() time.Duration
	SetCertFile(certfile, keyfile string)
	SetProtocolCallback(cb cb.OnProtocol)
	SetHandshakeCallback(cb cb.OnHandshake)
	SetConditionCallback(cb cb.OnCondition)
	SetNewConnectionCallback(cb cb.OnNewConnection)
	SetHandshakeTimeout(d time.Duration)
	SetIdleTimeout(d time.Duration)
	SetReadBufferSize(readBufferSize int)
}

// <summary>
// acceptor TCP接受器
// <summary>
type acceptor struct {
	certfile, keyfile string
	name              string
	started           bool
	listening         bool
	serving           bool
	stopping          cc.Singal
	closing           [2]cc.AtomFlag
	flag              [2]cc.AtomFlag
	lock              *sync.Mutex
	cond              *sync.Cond
	addr              *conn.Address
	upgrader          *websocket.Upgrader
	server            *http.Server
	listener          net.Listener
	channel           transmit.Channel
	onProtocol        cb.OnProtocol
	onHandshake       cb.OnHandshake
	onCondition       cb.OnCondition
	onNewConnection   cb.OnNewConnection
	handshakeTimeout  time.Duration
	idleTimeout       time.Duration
	readBufferSize    int
}

func NewAcceptor(name string, address ...string) Acceptor {
	s := &acceptor{
		name:     name,
		lock:     &sync.Mutex{},
		stopping: cc.NewSingal(),
		closing:  [2]cc.AtomFlag{cc.NewAtomFlag(), cc.NewAtomFlag()},
		flag:     [2]cc.AtomFlag{cc.NewAtomFlag(), cc.NewAtomFlag()},
	}
	if len(address) > 0 {
		s.addr = conn.ParseAddress(address[0])
	}
	s.cond = sync.NewCond(s.lock)
	return s
}

func (s *acceptor) SetHandshakeTimeout(d time.Duration) {
	s.handshakeTimeout = d
}

func (s *acceptor) GetIdleTimeout() time.Duration {
	return s.idleTimeout
}

func (s *acceptor) SetIdleTimeout(d time.Duration) {
	s.idleTimeout = d
}

func (s *acceptor) SetReadBufferSize(readBufferSize int) {
	s.readBufferSize = readBufferSize
}

func (s *acceptor) toName() {
	s.name = s.name + "#" + s.addr.Format() + ".acceptor"
}

func (s *acceptor) Addr() *conn.Address {
	return s.addr
}

func (s *acceptor) SetCertFile(certfile, keyfile string) {
	s.certfile = certfile
	s.keyfile = keyfile
}

func (s *acceptor) SetProtocolCallback(cb cb.OnProtocol) {
	s.onProtocol = cb
}

func (s *acceptor) SetHandshakeCallback(cb cb.OnHandshake) {
	s.onHandshake = cb
}

func (s *acceptor) SetConditionCallback(cb cb.OnCondition) {
	s.onCondition = cb
}

func (s *acceptor) SetNewConnectionCallback(cb cb.OnNewConnection) {
	s.onNewConnection = cb
}

func (s *acceptor) assertProtocol() {
	if s.onProtocol == nil {
		panic(errors.New("error"))
	}
}

func (s *acceptor) assertOnNewConnection() {
	if s.onNewConnection == nil {
		panic(errors.New("error"))
	}
}

// func (s *acceptor) assertOnCondition() {
// 	if s.onCondition == nil {
// 		panic(errors.New("error"))
// 	}
// }

func (s *acceptor) ListenTCP(address ...string) {
	s.assertProtocol()
	// s.assertOnCondition()
	s.assertOnNewConnection()
	if len(address) > 0 {
		s.addr = conn.ParseAddress(address[0])
	}
	if s.addr != nil && !s.started && s.flag[0].TestSet() {
		s.close()
		s.listenTCP()
		s.flag[0].Reset()
	}
}

func (s *acceptor) Stop() {
	if s.started && s.flag[1].TestSet() {
		if s.server != nil {
			s.stop_serve()
		} else {
			s.stop_accept()
		}
		s.flag[1].Reset()
	}
}

func (s *acceptor) listenTCP() {
	// logs.Warnf("addr=%v", s.addr.Addr)
	s.toName()
	listener, err := net.Listen("tcp", s.addr.Addr)
	if err != nil {
		logs.Errorf(err.Error())
		return
	}
	s.listener = listener
	s.listening = true
	s.channel = s.onProtocol(s.addr.Proto)
	logs.Debugf("%s", s.addr.Format())
	switch s.addr.Proto {
	case "ws", "wss":
		s.upgradeAndServe(s.addr)
	case "tcp":
		go s.accept()
		s.wait()
	}
}

func (s *acceptor) accept() {
	s.signal()
	var delay time.Duration
	// defer s.close()
	for !s.is_stopping() {
		c, err := s.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if delay == 0 {
					delay = 5 * time.Millisecond
				} else {
					delay *= 2
				}
				if max := 1 * time.Second; delay > max {
					delay = max
				}
				logs.Errorf("error: %v; retrying in %v", err, delay)
				time.Sleep(delay)
				continue
			}
			logs.Errorf(err.Error())
			return
		}
		switch conn.UsePool {
		case true:
			connpool.Do(cb.NewFunctor00(func() {
				if s.onCondition != nil && !s.onCondition(c.RemoteAddr()) {
					c.Close()
				} else if s.onNewConnection != nil {
					s.onNewConnection(c, s.channel, s.addr.Proto)
				} else {
					c.Close()
				}
			}))
		default:
			if s.onCondition != nil && !s.onCondition(c.RemoteAddr()) {
				c.Close()
			} else if s.onNewConnection != nil {
				s.onNewConnection(c, s.channel, s.addr.Proto)
			} else {
				c.Close()
			}
		}
	}
	s.cleanup()
}

//	&http.Request{
//		Method: http.MethodGet,
//		Header: http.Header{
//			"Upgrade":               []string{"websocket"},
//			"Connection":            []string{"upgrade"},
//			"Sec-Websocket-Key":     []string{"dGhlIHNhbXBsZSBub25jZQ=="},
//			"Sec-Websocket-Version": []string{"13"},
//		}}
func (s acceptor) upgradeAndServe(addr *conn.Address) {
	s.upgrader = &websocket.Upgrader{
		HandshakeTimeout: s.handshakeTimeout,
		ReadBufferSize:   s.readBufferSize,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
	mux := http.NewServeMux()
	mux.HandleFunc(addr.Path, func(w http.ResponseWriter, r *http.Request) {
		if s.onHandshake != nil && !s.onHandshake(w, r) {
			return
		}
		c, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		switch conn.UsePool {
		case true:
			connpool.Do(cb.NewFunctor00(func() {
				if s.onCondition != nil && !s.onCondition(c.RemoteAddr()) {
					c.Close()
				} else if s.onNewConnection != nil {
					s.onNewConnection(c, s.channel, s.addr.Proto, w, r)
				} else {
					c.Close()
				}
			}))
		default:
			if s.onCondition != nil && !s.onCondition(c.RemoteAddr()) {
				c.Close()
			} else if s.onNewConnection != nil {
				s.onNewConnection(c, s.channel, s.addr.Proto, w, r)
			} else {
				c.Close()
			}
		}
	})
	s.server = &http.Server{
		Addr:    addr.Addr,
		Handler: mux}
	// s.server = &http.Server{
	// 	Addr:        addr.Addr,
	// 	Handler:     mux,
	// 	IdleTimeout: s.idleTimeout}
	s.serving = true
	go s.serve()
	s.wait()
}

func (s *acceptor) serve() {
	s.signal()
	// defer s.close()
	if s.certfile != "" && s.keyfile != "" {
		err := s.server.ServeTLS(s.listener, s.certfile, s.keyfile)
		if err != nil {
			logs.Errorf(err.Error())
		}
	} else {
		err := s.server.Serve(s.listener)
		if err != nil {
			logs.Errorf(err.Error())
		}
	}
	s.cleanup()
}

func (s *acceptor) cleanup() {
	//先退出accept/serve再关闭listener
	s.close()
	s.server = nil
	s.listener = nil
	s.started = false
}

func (s *acceptor) wait() {
	s.lock.Lock()
	for !s.started {
		s.cond.Wait()
	}
	s.lock.Unlock()
}

func (s *acceptor) signal() {
	s.lock.Lock()
	s.started = true
	s.cond.Signal()
	s.lock.Unlock()
}

func (s *acceptor) is_stopping() bool {
	return s.stopping.Signaled()
}

func (s *acceptor) stop_accept() {
	s.stopping.Signal()
}

func (s *acceptor) stop_serve() {
	s.close_serve()
}

func (s *acceptor) close() {
	if s.listening && s.closing[0].TestSet() {
		s.listener.Close()
		s.listening = false
		s.closing[0].Reset()
	}
}

func (s *acceptor) close_serve() {
	if s.serving && s.closing[1].TestSet() {
		s.server.Close()
		s.serving = false
		s.closing[1].Reset()
	}
}
