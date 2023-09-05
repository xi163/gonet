package conn

import (
	"sync"
)

type HoldType uint8

const (
	KHoldNone HoldType = iota
	KHoldTemporary
	KHold
)

// 连接会话容器
type Sessions interface {
	Get(id int64) Session
	Count() int
	Add(peer Session) bool
	Remove(peer Session)
	Range(cb func(peer Session))
	CloseAll()
	Wait()
	Stop()
}

type sessions struct {
	peers map[int64]Session
	n     int
	l     *sync.RWMutex
	c     *sync.Cond
	stop  bool
	done  bool
}

func NewSessions() Sessions {
	s := &sessions{l: &sync.RWMutex{}, peers: map[int64]Session{}}
	s.c = sync.NewCond(s.l)
	return s
}

func (s *sessions) Get(id int64) Session {
	s.l.RLock()
	if peer, ok := s.peers[id]; ok {
		s.l.RUnlock()
		return peer
	}
	s.l.RUnlock()
	return nil
}

func (s *sessions) Count() int {
	s.l.RLock()
	//c = len(s.peers)
	c := s.n
	s.l.RUnlock()
	return c
}

func (s *sessions) Add(peer Session) bool {
	ok := false
	s.l.Lock()
	if !s.stop {
		// logs.Debugf("%v", peer.Name())
		s.peers[peer.ID()] = peer
		s.n++
		ok = true
	}
	s.l.Unlock()
	return ok
}

func (s *sessions) Remove(peer Session) {
	s.l.Lock()
	if _, ok := s.peers[peer.ID()]; ok {
		// logs.Debugf("%v", peer.Name())
		delete(s.peers, peer.ID())
		s.n--
	}
	// if s.stop && len(s.peers) == 0 {
	if s.stop && s.n == 0 {
		s.done = true
		s.c.Signal()
	}
	s.l.Unlock()
}

func (s *sessions) Range(cb func(peer Session)) {
	s.l.RLock()
	for _, peer := range s.peers {
		cb(peer)
	}
	s.l.RUnlock()
}

// s.closeAll -> peer.Close -> s.Remove
func (s *sessions) closeAll(stop bool) {
	s.l.RLock()
	if stop {
		s.stop = true
	}
	for _, peer := range s.peers {
		peer.Close()
	}
	s.l.RUnlock()
}

func (s *sessions) CloseAll() {
	s.closeAll(false)
}

func (s *sessions) Wait() {
	s.l.Lock()
	for !s.done {
		s.c.Wait()
	}
	s.l.Unlock()
}

func (s *sessions) Stop() {
	s.closeAll(true)
}
