package cc

import "errors"

// chan信号
type Singal interface {
	Signal()
	Signaled() (signal bool)
	Read() <-chan bool
	Close()
}

type singal struct {
	signal chan bool
	closed AtomFlag
}

func NewSingal() *singal {
	s := &singal{
		signal: make(chan bool, 1),
		closed: NewAtomFlag(),
	}
	return s
}

func (s *singal) Signal() {
	if !s.closed.IsSet() {
		//chan满则阻塞等待
		s.signal <- true
		s.Close()
	} else {
		s.closed.Reset()
		s.signal = make(chan bool, 1)
		//chan满则阻塞等待
		s.signal <- true
		s.Close()
	}
}

func (s *singal) Read() <-chan bool {
	if s.signal == nil {
		panic(errors.New("error: singal.Read signal is nil"))
	}
	return s.signal
}

func (s *singal) Close() {
	if s.closed.TestSet() {
		close(s.signal)
	}
}

func (s *singal) Signaled() (signal bool) {
	if !s.closed.IsSet() {
		select {
		case <-s.signal:
			signal = true
			break
		default:
			break
		}
	}
	return
}
