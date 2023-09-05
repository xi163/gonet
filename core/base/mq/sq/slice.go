package sq

import (
	_ "net/http/pprof"
	"sync"

	"github.com/cwloo/gonet/core/base/mq"
	"github.com/cwloo/gonet/core/cb"
)

// slice非阻塞队列
type slice struct {
	lock  *sync.Mutex
	slice []any
}

func NewSlice(size int) mq.Queue {
	// if size <= 0 {
	// 	panic(errors.New("error"))
	// }
	s := &slice{
		lock: &sync.Mutex{},
	}
	return s
}

func (s *slice) Name() string {
	return "slice"
}

func (s *slice) Push(data any) {
	s.lock.Lock()
	s.slice = append(s.slice, data)
	s.lock.Unlock()
}

// 一次取一个
func (s *slice) Pop() (data any, exit, empty bool, code int) {
	s.lock.Lock()
	length := len(s.slice)
	if length > 0 {
		data = s.slice[0]
		if data == nil {
			exit = true
		} else if m, ok := data.(*mq.ExitStruct); ok {
			exit = true
			code = m.Code
		}
		if length > 1 {
			s.slice = s.slice[1:]
		} else {
			s.slice = []any{}
		}
	} else {
		empty = true
	}
	s.lock.Unlock()
	return
}

// 批量全部取
func (s *slice) Pick() (v []any) {
	s.lock.Lock()
	if len(s.slice) > 0 {
		v = s.slice[:]
		s.slice = []any{}
	}
	s.lock.Unlock()
	return
}

// 批量全部取直到遇到nil
func (s *slice) Pick_until() (v []any, exit bool, code int) {
	s.lock.Lock()
	// for i, data := range s.slice {
	// 	if data == nil {
	// 		s.slice = append(s.slice[:i], s.slice[i+1:]...)
	// 		exit = true
	// 		break
	// 	} else {
	// 		v = append(v, data)
	// 		s.slice = append(s.slice[:i], s.slice[i+1:]...)
	// 	}
	// }
	for i := 0; i < len(s.slice); i++ {
		if s.slice[i] == nil {
			s.slice = append(s.slice[:i], s.slice[i+1:]...)
			exit = true
			break
		} else if m, ok := s.slice[i].(*mq.ExitStruct); ok {
			exit = true
			code = m.Code
			break
		} else {
			v = append(v, s.slice[i])
			s.slice = append(s.slice[:i], s.slice[i+1:]...)
		}
	}
	s.lock.Unlock()
	return
}

func (s *slice) exec_step(handler cb.Processor, args ...any) (exit bool, code int) {
	msg, EXIT, empty, CODE := s.Pop()
	if EXIT {
		exit = EXIT
		code = CODE
	} else if !empty {
		if handler(msg, args...) {
			exit = true
			return
		}
	}
	return
}

func (s *slice) exec_step_until(handler cb.Processor, args ...any) (exit bool, code int) {
	msg, EXIT, empty, CODE := s.Pop()
	if EXIT {
		exit = EXIT
		code = CODE
	} else if !empty {
		if handler(msg, args...) {
			exit = true
			return
		}
	}
	return
}

func (s *slice) exec_all(handler cb.Processor, args ...any) (exit bool, code int) {
	msgs := s.Pick()
	for _, msg := range msgs {
		if msg == nil {
			exit = true
		} else if m, ok := msg.(*mq.ExitStruct); ok {
			exit = true
			code = m.Code
		} else {
			if handler(msg, args...) {
				exit = true
				return
			}
		}
	}
	return
}

func (s *slice) exec_all_until(handler cb.Processor, args ...any) (exit bool, code int) {
	msgs, EXIT, CODE := s.Pick_until()
	for _, msg := range msgs {
		if handler(msg, args...) {
			exit = true
			return
		}
	}
	exit = EXIT
	code = CODE
	return
}

// 一次取一个或批量全部取
func (s *slice) Exec(step bool, handler cb.Processor, args ...any) (exit bool, code int) {
	if step {
		exit, code = s.exec_step(handler, args...)
	} else {
		exit, code = s.exec_all(handler, args...)
	}
	return
}

// 一次取一个或批量全部取直到遇到nil
func (s *slice) Exec_until(step bool, handler cb.Processor, args ...any) (exit bool, code int) {
	if step {
		exit, code = s.exec_step_until(handler, args...)
	} else {
		exit, code = s.exec_all_until(handler, args...)
	}
	return
}

func (s *slice) Size() int {
	s.lock.Lock()
	c := len(s.slice)
	s.lock.Unlock()
	return c
}

func (s *slice) Reset() {
}
