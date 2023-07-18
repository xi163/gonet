package lq

import (
	"container/list"
	"sync"

	"github.com/cwloo/gonet/core/base/mq"
	"github.com/cwloo/gonet/core/cb"
)

// <summary>
// queue list阻塞队列
// <summary>
type queue struct {
	lock *sync.Mutex
	cond *sync.Cond
	list *list.List
}

func NewQueue(size int) mq.BlockQueue {
	// if size <= 0 {
	// 	panic(errors.New("error"))
	// }
	s := &queue{
		list: list.New(),
		lock: &sync.Mutex{}}
	s.cond = sync.NewCond(s.lock)
	return s
}

func (s *queue) Name() string {
	return "list"
}

func (s *queue) Push(data any) {
	s.lock.Lock()
	s.list.PushBack(data)
	s.cond.Signal()
	s.lock.Unlock()
}

// 一次取一个
func (s *queue) Pop() (data any, exit, empty bool, code int) {
	s.lock.Lock()
	for s.list.Len() == 0 {
		s.cond.Wait()
	}
	if elem := s.list.Front(); elem != nil {
		data = elem.Value
		if data == nil {
			exit = true
		} else if m, ok := data.(*mq.ExitStruct); ok {
			exit = true
			code = m.Code
		}
		s.list.Remove(elem)
		elem = nil
	}
	s.lock.Unlock()
	return
}

// 批量全部取
func (s *queue) Pick() (v []any) {
	s.lock.Lock()
	for s.list.Len() == 0 {
		s.cond.Wait()
	}
	s.swap(func(elem *list.Element) {
		data := elem.Value
		v = append(v, data)
	})
	s.lock.Unlock()
	return
}

// 批量全部取直到遇到nil
func (s *queue) Pick_until() (v []any, exit bool, code int) {
	s.lock.Lock()
	for s.list.Len() == 0 {
		s.cond.Wait()
	}
	exit, code = s.swap_until(func(elem *list.Element) {
		data := elem.Value
		v = append(v, data)
	})
	s.lock.Unlock()
	return
}

func (s *queue) exec_step(handler cb.Processor, args ...any) (exit bool, code int) {
	msg, EXIT, empty, CODE := s.Pop()
	if EXIT {
		exit = EXIT
		code = CODE
	} else if !empty {
		if _, ok := msg.(*mq.WakeupStruct); !ok {
			if handler(msg, args...) {
				exit = true
				return
			}
		}
	}
	return
}

func (s *queue) exec_step_until(handler cb.Processor, args ...any) (exit bool, code int) {
	msg, EXIT, empty, CODE := s.Pop()
	if EXIT {
		exit = EXIT
		code = CODE
	} else if !empty {
		if _, ok := msg.(*mq.WakeupStruct); !ok {
			if handler(msg, args...) {
				exit = true
				return
			}
		}
	}
	return
}

func (s *queue) exec_all(handler cb.Processor, args ...any) (exit bool, code int) {
	msgs := s.Pick()
	for _, msg := range msgs {
		if msg == nil {
			exit = true
		} else if m, ok := msg.(*mq.ExitStruct); ok {
			exit = true
			code = m.Code
		} else if _, ok := msg.(*mq.WakeupStruct); !ok {
			if handler(msg, args...) {
				exit = true
				return
			}
		}
	}
	return
}

func (s *queue) exec_all_until(handler cb.Processor, args ...any) (exit bool, code int) {
	msgs, EXIT, CODE := s.Pick_until()
	for _, msg := range msgs {
		if _, ok := msg.(*mq.WakeupStruct); !ok {
			if handler(msg, args...) {
				exit = true
				return
			}
		}
	}
	exit = EXIT
	code = CODE
	return
}

// 一次取一个或批量全部取
func (s *queue) Exec(step bool, handler cb.Processor, args ...any) (exit bool, code int) {
	if step {
		exit, code = s.exec_step(handler, args...)
	} else {
		exit, code = s.exec_all(handler, args...)
	}
	return
}

// 一次取一个或批量全部取直到遇到nil
func (s *queue) Exec_until(step bool, handler cb.Processor, args ...any) (exit bool, code int) {
	if step {
		exit, code = s.exec_step_until(handler, args...)
	} else {
		exit, code = s.exec_all_until(handler, args...)
	}
	return
}

func (s *queue) Size() int {
	s.lock.Lock()
	c := s.list.Len()
	s.lock.Unlock()
	return c
}

func (s *queue) Range(f func(*list.Element)) {
	var next *list.Element
	for elem := s.list.Front(); elem != nil; elem = next {
		f(elem)
		next = elem.Next()
	}
}

func (s *queue) swap(f func(*list.Element)) {
	var next *list.Element
	for elem := s.list.Front(); elem != nil; elem = next {
		f(elem)
		next = elem.Next()
		s.list.Remove(elem)
		elem = nil
	}
}

func (s *queue) swap_until(f func(*list.Element)) (exit bool, code int) {
	var next *list.Element
	for elem := s.list.Front(); elem != nil; elem = next {
		next = elem.Next()
		if elem.Value == nil {
			exit = true
			s.list.Remove(elem)
			elem = nil
			break
		} else if m, ok := elem.Value.(*mq.ExitStruct); ok {
			exit = true
			code = m.Code
			s.list.Remove(elem)
			elem = nil
			break
		}
		f(elem)
		s.list.Remove(elem)
		elem = nil
	}
	return
}

func (s *queue) clear() {
	var next *list.Element
	for elem := s.list.Front(); elem != nil; elem = next {
		next = elem.Next()
		s.list.Remove(elem)
		elem = nil
	}
}

func (s *queue) Wakeup() {
	s.Push(mq.NewWakeupStruct())
}
