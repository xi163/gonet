package pool

import (
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/pipe"
	run_gos "github.com/cwloo/gonet/core/base/run/gos"
	"github.com/cwloo/gonet/core/cb"
	logs "github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils/pool"
	"github.com/cwloo/gonet/utils/safe"
)

// go协程池
type Gos2 interface {
	Go2(f cb.Functor)
	Len() (c int)
}

type gos2 struct {
	i32  cc.I32
	pool pool.Pool
}

func NewGos2() Gos2 {
	s := &gos2{
		i32: cc.NewI32(),
	}
	s.pool = *pool.NewPoolWith(s.new)
	return s
}

func (s *gos2) Go2(f cb.Functor) {
	p, _ := s.Get()
	p.Do(cb.NewFunctor10(func(args any) {
		safe.Call2(f.Call)
		f.Put()
		s.Put(args.(pipe.Pipe))
	}, p))
}

func (s *gos2) Len() (c int) {
	return s.pool.Len()
}

func (s *gos2) new(cb func(error, ...any), v ...any) (p any, e error) {
	id := s.i32.New()
	nonblock := true
	runner := run_gos.NewProcessor(s.handler)
	p = pipe.NewPipe(id, "gos.pipe", 1, nonblock, runner)
	cb(nil)
	return
}

func (s *gos2) Get() (p pipe.Pipe, e error) {
	v, err := s.pool.Get()
	e = err
	switch err {
	case nil:
		p = v.(pipe.Pipe)
	default:
		logs.Errorf(err.Error())
	}
	return
}

func (s *gos2) Put(pipe pipe.Pipe) {
	s.pool.Put(pipe)
}

func (s *gos2) Close(reset func(pipe.Pipe)) {
	s.pool.Reset(func(value any) {
		reset(value.(pipe.Pipe))
		value.(pipe.Pipe).Close()
	}, false)
}

func (s *gos2) handler(msg any, args ...any) bool {
	switch msg := msg.(type) {
	case cb.Functor:
		safe.Call2(msg.Call)
		msg.Put()
	case cb.Timeout:
		if !msg.Expire().Expired(time.Now()) {
			switch data := msg.Data().(type) {
			case cb.Functor:
				data.CallWith(msg.Expire())
				data.Put()
			}
		}
		msg.Put()
	}
	return false
}
