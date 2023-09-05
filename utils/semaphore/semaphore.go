package semaphore

import "sync"

// 互斥访问控制
type Sem struct {
	w     *sync.Mutex
	l     *sync.Mutex
	c     *sync.Cond
	avail int64
	n     int64 //n个并发访问
}

func New(n int64) *Sem {
	s := &Sem{n: n, avail: n, l: &sync.Mutex{}, w: &sync.Mutex{}}
	s.c = sync.NewCond(s.l)
	return s
}

func (s *Sem) Enter() {
wait:
	s.wait()
	s.w.Lock()
	if s.avail > 0 {
		s.avail--
		s.w.Unlock()
	} else {
		s.w.Unlock()
		goto wait
	}
}

func (s *Sem) Leave() {
	s.w.Lock()
	if s.avail < s.n {
		s.avail++
		if s.avail == 1 {
			s.c.Signal()
		}
	}
	s.w.Unlock()
}

func (s *Sem) wait() {
	s.l.Lock()
	for s.avail == 0 {
		s.c.Wait()
	}
	s.l.Unlock()
}

type FreeSem struct {
	l     *sync.Mutex
	avail int64
	n     int64
}

func NewFreeSem(n int64) *FreeSem {
	s := &FreeSem{n: n, avail: n, l: &sync.Mutex{}}
	return s
}

func (s *FreeSem) Enter() (bv bool) {
	s.l.Lock()
	if s.avail > 0 {
		s.avail--
		bv = true
	}
	s.l.Unlock()
	return
}

func (s *FreeSem) Leave() {
	s.l.Lock()
	if s.avail < s.n {
		s.avail++
	}
	s.l.Unlock()
}

var sem = New(10)
var ix = 10

func TestSemaphore() {
	for i := 0; i < 100; i++ {
		go func() {
			for {
				sem.Enter()
				ix--
				println("1======= ", ix)
				ix++
				println("2======= ", ix)
				sem.Leave()
			}
		}()
	}
}

func OnInputTestSemaphore(str string) int {
	switch str {
	case "w":
		{
			for i := 0; i < 30; i++ {
				sem.Leave()
			}
			return 0
		}
	}
	return 0
}
