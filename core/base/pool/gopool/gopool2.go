package gopool

import (
	"github.com/cwloo/gonet/core/base/pool"
	"github.com/cwloo/gonet/core/cb"
)

var (
	gos2 = pool.NewGos2()
)

func Go2(f cb.Functor) {
	gos2.Go2(f)
}

func Len() (c int) {
	return gos2.Len()
}
