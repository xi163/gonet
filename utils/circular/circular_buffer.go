package circular

import (
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils/bucket"
)

type Buffer[T any] interface {
	Range(cb func(T) bool)
	Resize(newsize int)
	Empty() bool
	Full() bool
	Reserve() int
	Capacity() int
	Size() int
	PushFront(v T)
	PushBack(v T)
	PopFront()
	PopBack()
	Front() T
	Back() T
	At(i int) T
	Begin() int
	End() int
}

type ring_buffer[T any] struct {
	first, last int
	size, cap   int
	slice       []T
	reset_val   T
	construct   func() T
}

func New[T any](cap int, reset_val T) Buffer[T] {
	s := &ring_buffer[T]{
		first:     0,
		last:      0,
		size:      0,
		reset_val: reset_val,
		cap:       cap,
		slice:     make([]T, cap)} //len:cap cap:cap
	return s
}

func NewWitch[T any](cap int, f func() T) Buffer[T] {
	s := &ring_buffer[T]{
		first:     0,
		last:      0,
		size:      0,
		construct: f,
		cap:       cap,
		slice:     make([]T, cap)}
	return s
}

func (s *ring_buffer[T]) Range(cb func(T) bool) {
	// c := 0
	for i := s.size - 1; i >= 0; i-- {
		if cb(s.At(i)) {
			// logs.Infof("c = %d", c)
			return
		}
		// c++
	}
}

func (s *ring_buffer[T]) Resize(newsize int) {
	for i := 0; i < s.cap; i++ {
		s.slice[i] = s.construct()
	}
	s.size = s.cap
	s.first = 0
	s.last = s.size - 1
	// if newsize > s.size {
	// 	if newsize > s.cap {
	// 		s.setCapacity(newsize)
	// 	}
	// 	s.arrange()
	// 	s.insert(s.end(), newsize-s.size)
	// } else {
	// 	erase(end() - (size() - newsize), e);
	// }
}

func (s *ring_buffer[T]) Empty() bool {
	return s.size == 0
}

func (s *ring_buffer[T]) Full() bool {
	return s.cap == s.size
}

func (s *ring_buffer[T]) Reserve() int {
	return s.cap - s.size
}

func (s *ring_buffer[T]) Capacity() int {
	return s.cap
}

func (s *ring_buffer[T]) Size() int {
	return s.size
}

func (s *ring_buffer[T]) moveBack() {
	// logs.Warnf("cap:%d size:%d begin:%d end:%d", s.cap, s.size, s.first, s.last)
	//Must size > 0
	for i := s.last; i >= s.first; i-- {
		if i+1 >= s.cap {
		} else {
			s.slice[i+1] = s.slice[i]
		}
	}
	if s.last+1 >= s.cap {
		s.size--
	} else {
		s.last++
	}
	s.reset(s.first)
	if s.first < s.last {
		s.first++
	}
	// logs.Warnf("cap:%d size:%d begin:%d end:%d", s.cap, s.size, s.first, s.last)
}

func (s *ring_buffer[T]) PushFront(v T) {
	if s.first > 0 {
		if s.size > 0 {
			s.first--
		}
	} else {
		if s.size > 0 {
			s.moveBack()
			s.first--
		}
	}
	s.slice[s.first] = v
	s.size++
}

func (s *ring_buffer[T]) moveFront() {
	// logs.Warnf("cap:%d size:%d begin:%d end:%d", s.cap, s.size, s.first, s.last)
	//Must size > 0
	for i := s.first; i <= s.last; i++ {
		if i-1 >= 0 {
			s.slice[i-1] = s.slice[i]
		}
	}
	if s.first-1 >= 0 {
		s.first--
	} else {
		s.size--
	}
	s.reset(s.last)
	if s.first < s.last {
		s.last--
	}
	// logs.Warnf("cap:%d size:%d begin:%d end:%d", s.cap, s.size, s.first, s.last)
}

func (s *ring_buffer[T]) PushBack(v T) {
	if s.last+1 < s.cap {
		if s.size > 0 {
			s.last++
		}
	} else {
		if s.size > 0 {
			s.moveFront()
			s.last++
		}
	}
	s.slice[s.last] = v
	s.size++
}

// func (s *ring_buffer[T]) setCapacity(newsize int) {

// }

// func (s *ring_buffer[T]) insert(i int, size int, construct_empty_obj T) {
// 	start := s.first + i
// 	c := size
// 	// for i := s.first; i < size; i++ {
// 	// 	s.slice[pos+i] = construct_emtpy_obj
// 	// }
// }

func (s *ring_buffer[T]) PopFront() {
	if s.size > 0 {
		s.reset(s.first)
		s.size--
		if s.first < s.last {
			s.first++
		}
	}
}

func (s *ring_buffer[T]) PopBack() {
	if s.size > 0 {
		s.reset(s.last)
		s.size--
		if s.first < s.last {
			s.last--
		}
	}
}

func (s *ring_buffer[T]) Begin() int {
	return s.first
}

func (s *ring_buffer[T]) End() int {
	return s.last
}

func (s *ring_buffer[T]) Front() T {
	return s.slice[s.first]
}

func (s *ring_buffer[T]) Back() T {
	return s.slice[s.last]
}

func (s *ring_buffer[T]) At(i int) T {
	return s.slice[s.first+i]
}

func (s *ring_buffer[T]) reset(i int) {
	switch s.construct {
	case nil:
		s.slice[i] = s.reset_val
	default:
		s.slice[i] = s.construct()
	}
}

// func (s *ring_buffer[T]) arrange() {
// 	for i := 0; i < s.size; i++ {
// 		s.slice[i] = s.slice[s.first+i]
// 	}
// }

func Test001() {
	cb := New[int](3, 0)
	logs.Debugf("cap:%d size:%d begin:%d end:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End())
	cb.PushBack(1)
	logs.Debugf("cap:%d size:%d begin:%d end:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End())
	cb.PushBack(2)
	logs.Debugf("cap:%d size:%d begin:%d end:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End())
	cb.PushBack(3)
	logs.Debugf("cap:%d size:%d begin:%d end:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End())
	for i := 0; i < cb.Size(); i++ {
		logs.Debugf("%d", cb.At(i))
	}
	logs.Debugf("-------------------------------------")
	cb.PushBack(4)
	logs.Debugf("cap:%d size:%d begin:%d end:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End())
	for i := 0; i < cb.Size(); i++ {
		logs.Debugf("%d", cb.At(i))
	}
	logs.Debugf("-------------------------------------")
	cb.PushFront(5)
	logs.Debugf("cap:%d size:%d begin:%d end:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End())
	logs.Debugf("cap:%d size:%d", cb.Capacity(), cb.Size())
	for i := 0; i < cb.Size(); i++ {
		logs.Debugf("%d", cb.At(i))
	}
}

func Test002() {
	logs.Debugf("-------------------------------------")
	cb := NewWitch[*bucket.Bucket](3, bucket.NewBucket)
	logs.Debugf("cap:%d size:%d begin:%d end:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End())
	cb.Resize(3)
	logs.Debugf("cap:%d size:%d begin:%d end:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End())
	cb.Back().Add(1)
	cb.Back().Add(2)
	cb.Back().Add(3)
	logs.Debugf("cap:%d size:%d begin:%d end:%d back().size:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End(), cb.Back().Len())
}

func Test003() {
	cb := New[int](4, 0)
	cb.PushBack(1)
	cb.PushBack(2)
	cb.PushBack(3)
	for i := 0; i < cb.Size(); i++ {
		logs.Debugf("%d", cb.At(i))
	}
	logs.Debugf("-------------------------------------")
	cb.PopFront()
	logs.Debugf("cap:%d size:%d begin:%d end:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End())
	for i := 0; i < cb.Size(); i++ {
		logs.Debugf("%d", cb.At(i))
	}
	logs.Debugf("-------------------------------------")
	cb.PushBack(4)
	logs.Debugf("cap:%d size:%d begin:%d end:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End())
	for i := 0; i < cb.Size(); i++ {
		logs.Debugf("%d", cb.At(i))
	}
	logs.Debugf("-------------------------------------")
	cb.PushBack(5)
	logs.Debugf("cap:%d size:%d begin:%d end:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End())
	for i := 0; i < cb.Size(); i++ {
		logs.Debugf("%d", cb.At(i))
	}
}

func Test004() {
	cb := New[int](4, 0)
	cb.PushBack(1)
	cb.PushBack(2)
	cb.PushBack(3)
	for i := 0; i < cb.Size(); i++ {
		logs.Debugf("%d", cb.At(i))
	}
	logs.Debugf("-------------------------------------")
	cb.PopFront()
	logs.Debugf("cap:%d size:%d begin:%d end:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End())
	for i := 0; i < cb.Size(); i++ {
		logs.Debugf("%d", cb.At(i))
	}
	logs.Debugf("-------------------------------------")
	cb.PushFront(4)
	logs.Debugf("cap:%d size:%d begin:%d end:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End())
	for i := 0; i < cb.Size(); i++ {
		logs.Debugf("%d", cb.At(i))
	}
	logs.Debugf("-------------------------------------")
	cb.PushFront(5)
	logs.Debugf("cap:%d size:%d begin:%d end:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End())
	for i := 0; i < cb.Size(); i++ {
		logs.Debugf("%d", cb.At(i))
	}
	logs.Debugf("-------------------------------------")
	cb.PushFront(6)
	logs.Debugf("cap:%d size:%d begin:%d end:%d", cb.Capacity(), cb.Size(), cb.Begin(), cb.End())
	for i := 0; i < cb.Size(); i++ {
		logs.Debugf("%d", cb.At(i))
	}
}
