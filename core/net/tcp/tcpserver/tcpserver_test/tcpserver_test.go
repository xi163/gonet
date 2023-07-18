package tcpserver_test

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cwloo/gonet/core/net/conn"
	"github.com/cwloo/gonet/core/net/tcp/tcpserver"
	"github.com/cwloo/gonet/core/net/transmit"
	wschannel_tst "github.com/cwloo/gonet/core/net/transmit/wschannel/wschannel_tst"
	logs "github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils/conv"
	"github.com/cwloo/gonet/utils/timestamp"
	"github.com/cwloo/gonet/utils/user_context"
	"github.com/cwloo/gonet/utils/user_session"
)

type EchoServer struct {
	addr        string
	numMaxConns int
	n           int32
	server      tcpserver.TCPServer
	users       *user_session.UserToPlatforms
}

func NewServer(addr string) *EchoServer {
	s := &EchoServer{}
	s.addr = addr
	s.numMaxConns = 1000
	s.users = user_session.NewUserToPlatforms()
	s.server = tcpserver.NewTCPServer("EchoServer", s.addr)
	s.server.SetProtocolCallback(s.onProtocol)
	s.server.SetHandshakeCallback(s.onHandshake)
	s.server.SetConnectedCallback(s.onConnected)
	s.server.SetClosedCallback(s.onClosed)
	s.server.SetMessageCallback(s.onMessage)
	s.server.SetHandshakeTimeout(time.Duration(10) * time.Second)
	s.server.SetIdleTimeout(30*time.Second, time.Second)
	s.server.SetReadBufferSize(1024)
	s.server.SetHoldType(conn.KHoldNone)
	logs.Tracef(s.addr)
	return s
}

func (s *EchoServer) run() {
	s.server.ListenTCP()
}

func (s *EchoServer) onProtocol(proto string) transmit.Channel {
	switch proto {
	case "tcp":
		panic("tcp Channel undefine")
	case "ws", "wss":
		return wschannel_tst.NewWSChannel()
	}
	panic(errors.New("no proto setup"))
}

func (s *EchoServer) onHandshake(w http.ResponseWriter, r *http.Request) bool {
	if atomic.LoadInt32(&s.n) >= int32(s.numMaxConns) {
		logs.Errorf("numMaxConns=%v", s.numMaxConns)
		return false
	}
	// query := r.URL.Query()
	logs.Infof("%v", r.URL.String())
	return true
}

func (s *EchoServer) onConnected(peer conn.Session, v ...any) {
	if peer.Connected() {
		num := atomic.AddInt32(&s.n, 1)
		logs.Debugf(" %d [%v] <- [%v]", num, peer.LocalAddr(), peer.RemoteAddr())
		ctx := user_context.NewCtx()
		peer.SetContext("ctx", ctx)
		s.addUserConn(peer)
		peer.Write(conv.StrToByte("server"))
	} else {
		panic("error")
	}
}

func (s *EchoServer) onClosed(peer conn.Session, reason conn.Reason) {
	if peer.Connected() {
		panic("error")
	} else {
		num := atomic.AddInt32(&s.n, -1)
		logs.Tracef(" %d [%v] <- [%v] %v", num, peer.LocalAddr(), peer.RemoteAddr(), reason.Msg)
		s.delUserConn(peer)
		peer.SetContext("ctx", nil)
	}
}

func (s *EchoServer) onMessage(peer conn.Session, msg any, recvTime timestamp.T) {
	logs.Debugf("%v", string(msg.([]byte)))
	// peer.Write(utils.Str2Byte("server"))
}

func (s *EchoServer) addUserConn(peer conn.Session) {
	// ctx := peer.GetContext("ctx").(user_context.Ctx)
	// s.users.AddUserConn(ctx.GetUserId(), ctx.GetPlatformId(), ctx.GetSession(), peer)
}

func (s *EchoServer) delUserConn(peer conn.Session) {
	// ctx := peer.GetContext("ctx").(user_context.Ctx)
	// s.users.DelUserConn(ctx.GetUserId(), ctx.GetPlatformId(), ctx.GetSession())
}

func TestMain(m *testing.M) {
	m.Run()
}

func Test(t *testing.T) {
	t.Run("tcpserver_test:", tcpserver_test)
}

func tcpserver_test(t *testing.T) {
	path, _ := os.Executable()
	dir, exec := filepath.Split(path)
	logs.SetTimezone(logs.MY_CST)
	logs.SetMode(logs.M_STDOUT_FILE)
	logs.SetStyle(logs.F_DETAIL)
	logs.SetLevel(logs.LVL_DEBUG)
	logs.Init(dir+"/logs", exec, 100000000)
	var wg sync.WaitGroup
	wg.Add(1)
	s := NewServer("ws://:7788/")
	s.run()
	wg.Wait()
}
