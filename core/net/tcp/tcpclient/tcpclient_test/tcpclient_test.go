package tcpclienttest_test

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/cwloo/gonet/core/net/conn"
	"github.com/cwloo/gonet/core/net/tcp/tcpclient"
	"github.com/cwloo/gonet/core/net/transmit"
	wschannel_tst "github.com/cwloo/gonet/core/net/transmit/wschannel/wschannel_tst"
	logs "github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils/timestamp"
	"github.com/cwloo/gonet/utils/user_context"
)

type EchoClient struct {
	addr   string
	client tcpclient.TCPClient
}

func NewClient(addr string) *EchoClient {
	s := &EchoClient{}
	s.addr = addr
	s.client = tcpclient.NewTCPClient("EchoClient", s.addr)
	s.client.SetProtocolCallback(s.onProtocol)
	s.client.SetConnectedCallback(s.onConnected)
	s.client.SetClosedCallback(s.onClosed)
	s.client.SetMessageCallback(s.onMessage)
	s.client.SetDialTimeout(10 * time.Second)
	s.client.SetIdleTimeout(30*time.Second, time.Second)
	s.client.EnableRetry(true)
	s.client.SetRetryInterval(10 * time.Second)
	s.client.SetHoldType(conn.KHoldNone)
	logs.Tracef(s.addr)
	return s
}

func (s *EchoClient) connect() {
	s.client.ConnectTCP(nil, s.addr)
}

func (s *EchoClient) onProtocol(proto string) transmit.Channel {
	switch proto {
	case "tcp":
		panic("tcp Channel undefine")
	case "ws", "wss":
		return wschannel_tst.NewWSChannel()
	}
	panic(errors.New("no proto setup"))
}

func (s *EchoClient) onConnected(peer conn.Session, v ...any) {
	if peer.Connected() {
		logs.Infof("[%v] -> [%v]", peer.LocalAddr(), peer.RemoteAddr())
		ctx := user_context.NewCtx()
		peer.SetContext("ctx", ctx)
		// peer.Write(utils.Str2Byte("client"))
	} else {
		panic("error")
	}
}

func (s *EchoClient) onClosed(peer conn.Session, reason conn.Reason, v ...any) {
	if peer.Connected() {
		panic("error")
	} else {
		logs.Tracef("[%v] -> [%v] %v", peer.LocalAddr(), peer.RemoteAddr(), reason.Msg)
		peer.SetContext("ctx", nil)
	}
}

func (s *EchoClient) onMessage(peer conn.Session, msg any, msgType int, recvTime timestamp.T) {
	logs.Debugf("%v", string(msg.([]byte)))
	// peer.Write(utils.Str2Byte("client"))
}

func TestMain(m *testing.M) {
	m.Run()
}

func Test(t *testing.T) {
	t.Run("tcpclient_test:", tcpclient_test)
}

func tcpclient_test(t *testing.T) {
	path, _ := os.Executable()
	dir, exec := filepath.Split(path)
	logs.SetTimezone(logs.MY_CST)
	logs.SetMode(logs.M_STDOUT_FILE)
	logs.SetStyle(logs.F_DETAIL)
	logs.SetLevel(logs.LVL_DEBUG)
	logs.Init(dir+"/logs", exec, 100000000)
	var wg sync.WaitGroup
	wg.Add(1)
	s := NewClient("ws://192.168.1.104:7788/")
	s.connect()
	wg.Wait()
}
