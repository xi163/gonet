package callpool

import (
	"time"

	"github.com/cwloo/gonet/core/base/pool"
	"github.com/cwloo/gonet/core/cb"
)

// 回调池(固定, 非阻塞)
var (
	calls = pool.NewCalls()
)

func Call(f cb.Functor) {
	calls.Call(f)
}

func GoTimeout(d time.Duration, f cb.Functor, fn cb.Functor) {
	calls.CallTimeout(d, f, fn)
}

func Start() {
	calls.Start()
}

func Stop() {
	calls.Stop()
}
