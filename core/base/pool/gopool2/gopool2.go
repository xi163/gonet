package gopool2

import (
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/pipe"
	"github.com/cwloo/gonet/core/base/run/gos"
	"github.com/cwloo/gonet/core/cb"
	logs "github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils/pool"
	"github.com/cwloo/gonet/utils/safe"
)

var (
	pool_ = NewGos()
)

func Go(f func()) {
	pool_.Go(f)
}

func Put(pipe pipe.Pipe) {
	pool_.Put(pipe)
}

type Gos interface {
	Len() (c int)
	Go(f func())
	Put(pipe pipe.Pipe)
}

type gos_ struct {
	i32  cc.I32
	pool pool.Pool
}

func NewGos() Gos {
	s := &gos_{
		i32: cc.NewI32(),
	}
	s.pool = *pool.NewPoolWith(s.new)
	return s
}

func (s *gos_) Go(f func()) {
	p, _ := s.Get()
	p.Do(cb.NewFunctor10(func(args any) {
		f()
		s.Put(args.(pipe.Pipe))
	}, p))
}

func (s *gos_) Len() (c int) {
	return s.pool.Len()
}

func (s *gos_) new(cb func(error, ...any), v ...any) (p any, e error) {
	id := s.i32.New()
	nonblock := true
	runner := gos.NewProcessor(s.handler)
	p = pipe.NewPipe(id, "gos.pipe", 1, nonblock, runner)
	return
}

func (s *gos_) Get() (p pipe.Pipe, e error) {
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

func (s *gos_) Put(pipe pipe.Pipe) {
	s.pool.Put(pipe)
}

func (s *gos_) Close(reset func(pipe.Pipe)) {
	s.pool.Reset(func(value any) {
		reset(value.(pipe.Pipe))
		value.(pipe.Pipe).Close()
	}, false)
}

func (s *gos_) handler(msg any, args ...any) bool {
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
