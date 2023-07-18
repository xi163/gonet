package gopool

import (
	"time"

	"github.com/cwloo/gonet/core/base/pool"
	"github.com/cwloo/gonet/core/cb"
)

// go协程池(动态, 非阻塞)
var (
	gos = pool.NewGos()
)

func Go(f cb.Functor) {
	gos.Go(f)
}

func GoTimeout(d time.Duration, f cb.Functor, fn cb.Functor) {
	gos.GoTimeout(d, f, fn)
}

func Num() int {
	return gos.Num()
}

func ResetNum() {
	gos.ResetNum()
}

func Start() {
	gos.Start()
}

func Stop() {
	gos.Stop()
}
