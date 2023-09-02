package connpool

import (
	"time"

	"github.com/cwloo/gonet/core/base/pool"
	"github.com/cwloo/gonet/core/cb"
)

// 回调池(固定, 阻塞)
var (
	conns = pool.NewConns()
)

func Do(f cb.Functor) {
	conns.Do(f)
}

func GoTimeout(d time.Duration, f cb.Functor, fn cb.Functor) {
	conns.DoTimeout(d, f, fn)
}

func Start() {
	conns.Start()
}

func Stop() {
	conns.Stop()
}
