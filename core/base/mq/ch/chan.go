package ch

import (
	"errors"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/mq"
	"github.com/cwloo/gonet/core/base/mq/lq"
	"github.com/cwloo/gonet/core/cb"
)

// <summary>
// Chan 消息队列
// <summary>
type Chan struct {
	nonblock bool
	size     int32
	mq       chan any
	wakeup   chan bool
	pendings mq.Queue
	closed   [2]cc.AtomFlag
}

func NewChan(size int, nonblock bool) Queue {
	if size <= 0 {
		panic(errors.New("NewChan error: size"))
	}
	s := &Chan{
		nonblock: nonblock,
		size:     int32(size),
		mq:       make(chan any, size),
		wakeup:   make(chan bool, 1),
		pendings: lq.NewList(1000),
		closed:   [2]cc.AtomFlag{cc.NewAtomFlag(), cc.NewAtomFlag()},
	}
	return s
}

func (s *Chan) Name() string {
	return "chan"
}

func (s *Chan) Push(data any) {
	if s.nonblock {
		s.push_nonblock(data)
	} else {
		s.push_blocking(data)
	}
}

func (s *Chan) push_blocking(data any) {
	if data == nil {
		if !s.closed[0].IsSet() {
			select {
			//chan满则阻塞等待
			case s.mq <- data:
				break
			}
			// s.close_msq()
		}
	} else {
		if !s.closed[0].IsSet() {
			select {
			//chan满则阻塞等待
			case s.mq <- data:
				break
			}
		}
	}
}

func (s *Chan) push_nonblock(data any) {
	if data == nil {
		if !s.closed[0].IsSet() {
			select {
			//屏蔽default，chan满则阻塞等待
			//启用default，chan满则执行default语句
			case s.mq <- data:
				break
			default:
				s.push_pending(data)
				break
			}
			// s.close_msq()
		}
	} else {
		if !s.closed[0].IsSet() {
			select {
			//chan满则执行default语句
			case s.mq <- data:
				break
			default:
				s.push_pending(data)
				break
			}
		}
	}
}

func (s *Chan) signal() {
	if len(s.wakeup) == cap(s.wakeup) {
	}
	select {
	case s.wakeup <- true:
		break
	default:
		break
	}
}

func (s *Chan) push_pending(data any) {
	s.pendings.Push(data)
	s.signal() //可能会阻塞，不能放在locker里面
}

// 一次取一个
func (s *Chan) Pop() (data any, exit, empty bool, code int) {
	data, exit, empty, code = s.pendings.Pop()
	return
}

// 批量全部取
func (s *Chan) Pick() (v []any) {
	v = s.pendings.Pick()
	return
}

// 批量全部取直到遇到nil
func (s *Chan) Pick_until() (v []any, exit bool, code int) {
	v, exit, code = s.pendings.Pick_until()
	return
}

// 一次取一个或批量全部取
func (s *Chan) Exec(step bool, handler cb.Processor, args ...any) (exit bool, code int) {
	exit, code = s.pendings.Exec(step, handler, args...)
	return
}

// 一次取一个或批量全部取直到遇到nil
func (s *Chan) Exec_until(step bool, handler cb.Processor, args ...any) (exit bool, code int) {
	exit, code = s.pendings.Exec_until(step, handler, args...)
	return
}

func (s *Chan) Full() bool {
	return len(s.mq) == cap(s.mq)
}

func (s *Chan) Length() int {
	return len(s.mq)
}

func (s *Chan) AssertEmpty() {
	if len(s.mq) != 0 {
		panic(errors.New("error"))
	}
}

func (s *Chan) Size() int {
	return s.pendings.Size()
}

func (s *Chan) Busing() bool {
	return len(s.mq) == cap(s.mq) || len(s.mq)+s.Size() > cap(s.mq)
}

func (s *Chan) Read() <-chan any {
	return s.mq
}

func (s *Chan) Signal() <-chan bool {
	return s.wakeup
}

func (s *Chan) Reset() {
	s.close_signal()
	s.close_msq()
}

func (s *Chan) close_msq() {
	if s.closed[0].TestSet() {
		close(s.mq)
	}
}

func (s *Chan) close_signal() {
	if s.closed[1].TestSet() {
		close(s.wakeup)
	}
}
