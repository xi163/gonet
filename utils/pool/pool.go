package pool

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/cwloo/gonet/logs"
)

// <summary>
// Pool
// <summary>
type Pool struct {
	vec *list.List
	new func(func(error, ...any), ...any) (any, error)
	l   *sync.RWMutex
}

func NewPool() *Pool {
	return &Pool{vec: list.New(), l: &sync.RWMutex{}}
}

func NewPoolWith(new func(func(error, ...any), ...any) (any, error)) *Pool {
	return &Pool{vec: list.New(), new: new, l: &sync.RWMutex{}}
}

func (s *Pool) SetNew(new func(func(error, ...any), ...any) (any, error)) {
	s.new = new
}

func (s *Pool) Len() (c int) {
	s.l.RLock()
	c = s.vec.Len()
	s.l.RUnlock()
	return
}

func (s *Pool) Range(cb func(value any)) {
	s.l.RLock()
	for elem := s.vec.Front(); elem != nil; elem = elem.Next() {
		cb(elem.Value)
	}
	s.l.RUnlock()
}

func (s *Pool) Update(cb func(any, func(error, ...any)) error) {
	s.l.Lock()
	var next *list.Element
	for elem := s.vec.Front(); elem != nil; elem = next {
		next = elem.Next()
		err := cb(&elem.Value, s.onUpdate)
		switch err {
		case nil:
		default:
			s.vec.Remove(elem)
		}
	}
	s.l.Unlock()
}

func (s *Pool) get(cb func(int, ...any), v ...any) (any, error) {
	s.l.Lock()
	switch s.vec.Len() {
	case 0:
		s.l.Unlock()
		return s.new(s.onNew, v...)
	default:
		elem := s.vec.Front()
		value := elem.Value
		s.vec.Remove(elem)
		s.l.Unlock()
		cb(s.Len(), v...)
		return value, nil
	}
}

func (s *Pool) Get(v ...any) (any, error) {
	switch s.Len() {
	case 0:
		return s.new(s.onNew, v...)
	default:
		return s.get(s.onGet, v...)
	}
}

func (s *Pool) Put(value any) {
	s.l.Lock()
	s.vec.PushBack(value)
	s.l.Unlock()
}

func (s *Pool) Reset(cb func(value any), direct bool) {
	s.l.Lock()
	switch direct {
	case true:
		s.vec.Init()
	default:
		var next *list.Element
		for elem := s.vec.Front(); elem != nil; elem = next {
			next = elem.Next()
			cb(elem.Value)
			s.vec.Remove(elem)
		}
	}
	s.l.Unlock()
}

func (s *Pool) onUpdate(err error, v ...any) {
	switch err {
	case nil:
		switch len(v) {
		case 0:
			logs.Errorf("ok")
		default:
			logs.Errorf("%v ok", fmt.Sprintf("%v", v))
		}
	default:
		switch len(v) {
		case 0:
			logs.Errorf("%v", err.Error())
		default:
			logs.Errorf("%v %v", fmt.Sprintf("%v", v), err.Error())
		}
	}
}

func (s *Pool) onNew(err error, v ...any) {
	switch err {
	case nil:
		switch len(v) {
		case 0:
			logs.Debugf("ok len=%v", s.Len())
		default:
			logs.Debugf("%v ok len=%v", fmt.Sprintf("%v", v), s.Len())
		}
	default:
		switch len(v) {
		case 0:
			logs.Errorf("%v", err.Error())
		default:
			logs.Errorf("%v %v", fmt.Sprintf("%v", v), err.Error())
		}
	}
}

func (s *Pool) onGet(length int, v ...any) {
	switch len(v) {
	case 0:
		logs.Errorf("ok len=%v", length)
	default:
		logs.Errorf("%v ok len=%v", fmt.Sprintf("%v", v), length)
	}
}
