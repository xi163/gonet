package bucket

import "sync"

var (
	pool = sync.Pool{
		New: func() any {
			return NewBucket()
		},
	}
)

type Bucket struct {
	m map[any]bool
	l *sync.Mutex
}

func NewBucket() *Bucket {
	return &Bucket{
		m: map[any]bool{},
		l: &sync.Mutex{}}
}

func (s *Bucket) Add(val any) {
	s.l.Lock()
	s.m[val] = true
	s.l.Unlock()
}

func (s *Bucket) Len() int {
	s.l.Lock()
	len := len(s.m)
	s.l.Unlock()
	return len
}

func (s *Bucket) Remove(val any) bool {
	s.l.Lock()
	if _, ok := s.m[val]; ok {
		delete(s.m, val)
		s.l.Unlock()
		return true
	}
	s.l.Unlock()
	return false
}

func (s *Bucket) Pop() (v []any) {
	s.l.Lock()
	for val := range s.m {
		v = append(v, val)
	}
	if len(s.m) > 0 {
		s.m = map[any]bool{}
	}
	s.l.Unlock()
	return
}
