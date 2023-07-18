package ordermap

import (
	"container/list"
)

type Pair struct {
	Key any
	Val any
}

type M struct {
	list *list.List
}

func New() *M {
	return &M{list: list.New()}
}

func (s *M) Insert(key any, value any, compare func(a, b any) bool) {
	pos := s.list.Front()
	for ; pos != nil; pos = pos.Next() {
		if !compare(key, pos.Value.(*Pair).Key) {
			data := &Pair{Key: key, Val: value}
			s.list.InsertBefore(data, pos)
			break
		}
	}
	if pos == nil {
		data := &Pair{Key: key, Val: value}
		s.list.PushBack(data)
	}
}

func (s *M) Top() (any, any) {
	if elem := s.list.Front(); elem != nil {
		data := elem.Value.(*Pair)
		return data.Key, data.Val
	}
	return nil, nil
}

func (s *M) Front() *list.Element {
	return s.list.Front()
}

func (s *M) Pop() {
	if elem := s.list.Front(); elem != nil {
		s.list.Remove(elem)
	}
}

func (s *M) Empty() bool {
	return s.list.Len() == 0
}
