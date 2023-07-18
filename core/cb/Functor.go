package cb

import (
	"sync"

	"github.com/cwloo/gonet/core/base/cc"
)

var (
	f00 = sync.Pool{
		New: func() any {
			return &Functor00{}
		},
	}
	f10 = sync.Pool{
		New: func() any {
			return &Functor10{}
		},
	}
	f20 = sync.Pool{
		New: func() any {
			return &Functor20{}
		},
	}
	f01 = sync.Pool{
		New: func() any {
			return &Functor01{}
		},
	}
	f11 = sync.Pool{
		New: func() any {
			return &Functor11{}
		},
	}
	f21 = sync.Pool{
		New: func() any {
			return &Functor21{}
		},
	}
)

// <summary>
// Functor 回调函数
// <summary>
type Functor interface {
	Call() (any, error)
	CallWith(expire cc.Expire) (any, error)
	Put()
}

type Functor00 struct {
	f func()
}

func NewFunctor00(f func()) Functor {
	s := f00.Get().(*Functor00)
	s.f = f
	return s
}

func (s *Functor00) Call() (v any, err error) {
	s.f()
	return
}

func (s *Functor00) CallWith(expire cc.Expire) (v any, err error) {
	s.f()
	return
}

func (s *Functor00) Put() {
	f00.Put(s)
}

type Functor10 struct {
	f    func(args any)
	args any
}

// func(v any){}, []any{a, b, c})
func NewFunctor10(f func(args any), args any) Functor {
	s := f10.Get().(*Functor10)
	s.f = f
	s.args = args
	return s
}

func (s *Functor10) Call() (v any, err error) {
	s.f(s.args)
	return
}

func (s *Functor10) CallWith(expire cc.Expire) (v any, err error) {
	s.f(s.args)
	return
}

func (s *Functor10) Put() {
	f10.Put(s)
}

type Functor20 struct {
	f    func(args ...any)
	args []any
}

// func(v ...any){}, a, b, c)
func NewFunctor20(f func(args ...any), args ...any) Functor {
	s := f20.Get().(*Functor20)
	s.f = f
	s.args = append(s.args, args...)
	return s
}

func (s *Functor20) Call() (v any, err error) {
	s.f(s.args...)
	return
}

func (s *Functor20) CallWith(expire cc.Expire) (v any, err error) {
	s.f(s.args...)
	return
}

func (s *Functor20) Put() {
	f20.Put(s)
}

type Functor01 struct {
	f func() (any, error)
}

func NewFunctor01(f func() (any, error)) Functor {
	s := f01.Get().(*Functor01)
	s.f = f
	return s
}

func (s *Functor01) Call() (any, error) {
	return s.f()
}

func (s *Functor01) CallWith(expire cc.Expire) (any, error) {
	return s.f()
}

func (s *Functor01) Put() {
	f01.Put(s)
}

type Functor11 struct {
	f    func(args any) (any, error)
	args any
}

func NewFunctor11(f func(args any) (any, error), args any) Functor {
	s := f11.Get().(*Functor11)
	s.f = f
	s.args = args
	return s
}

func (s *Functor11) Call() (any, error) {
	return s.f(s.args)
}

func (s *Functor11) CallWith(expire cc.Expire) (any, error) {
	return s.f(s.args)
}

func (s *Functor11) Put() {
	f11.Put(s)
}

type Functor21 struct {
	f    func(args ...any) (any, error)
	args []any
}

func NewFunctor21(f func(args ...any) (any, error), args ...any) Functor {
	s := f21.Get().(*Functor21)
	s.f = f
	s.args = append(s.args, args...)
	return s
}

func (s *Functor21) Call() (any, error) {
	return s.f(s.args...)
}

func (s *Functor21) CallWith(expire cc.Expire) (any, error) {
	return s.f(s.args...)
}

func (s *Functor21) Put() {
	f21.Put(s)
}
