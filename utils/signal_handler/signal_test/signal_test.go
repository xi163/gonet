package signal_test

import (
	"testing"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils/signal_handler"
)

func TestMain(m *testing.M) {
	m.Run()
}

func Test(t *testing.T) {
	t.Run("signal_handler_test:", signal_handler_test)
}

func signal_handler_test(t *testing.T) {
	signal_handler.RegisterStop(func() {
		logs.Debugf("Stopping..")
	})
	signal_handler.Wait()
}
