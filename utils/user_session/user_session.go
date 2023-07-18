package user_session

import (
	"sync"

	"github.com/cwloo/gonet/core/net/conn"
	"github.com/cwloo/gonet/logs"
)

// <summary>
// [session]=conn
// <summary>
type SessionToConn struct {
	l *sync.RWMutex
	m map[string]conn.Session
}

func NewSessionToConn() *SessionToConn {
	return &SessionToConn{m: map[string]conn.Session{}, l: &sync.RWMutex{}}
}

func (s *SessionToConn) Len() (c int) {
	s.l.RLock()
	c = len(s.m)
	s.l.RUnlock()
	return
}

func (s *SessionToConn) Get(sessionId string) (peer conn.Session, ok bool) {
	s.l.RLock()
	peer, ok = s.m[sessionId]
	s.l.RUnlock()
	return
}

func (s *SessionToConn) Do(sessionId string, cb func(conn.Session)) {
	var peer conn.Session
	s.l.RLock()
	if c, ok := s.m[sessionId]; ok {
		peer = c
		s.l.RUnlock()
		goto end
	}
	s.l.RUnlock()
	return
end:
	cb(peer)
}

func (s *SessionToConn) Add(sessionId string, peer conn.Session) (old conn.Session) {
	s.l.Lock()
	if c, ok := s.m[sessionId]; ok {
		old = c
	}
	s.m[sessionId] = peer
	s.l.Unlock()
	return
}

func (s *SessionToConn) Remove(sessionId string) (peer conn.Session) {
	_, ok := s.Get(sessionId)
	switch ok {
	case true:
		peer, _ = s.remove(sessionId)
	default:
	}
	return
}

func (s *SessionToConn) remove(sessionId string) (peer conn.Session, ok bool) {
	n := 0
	s.l.Lock()
	peer, ok = s.m[sessionId]
	switch ok {
	case true:
		delete(s.m, sessionId)
		n = len(s.m)
		s.l.Unlock()
		goto OK
	}
	s.l.Unlock()
	return
OK:
	logs.Errorf("%v size=%v", sessionId, n)
	return
}

func (s *SessionToConn) Range(cb func(string, conn.Session)) {
	s.l.RLock()
	for sessionId, peer := range s.m {
		cb(sessionId, peer)
	}
	s.l.RUnlock()
}

// <summary>
// [platformid][session]=conn
// <summary>
type PlatformToSessions struct {
	l *sync.RWMutex
	m map[int]*SessionToConn
}

func NewPlatformToSessions() *PlatformToSessions {
	return &PlatformToSessions{m: map[int]*SessionToConn{}, l: &sync.RWMutex{}}
}

func (s *PlatformToSessions) Len() int {
	s.l.RLock()
	c := len(s.m)
	s.l.RUnlock()
	return c
}

func (s *PlatformToSessions) Get(platformId int) (c *SessionToConn, ok bool) {
	s.l.RLock()
	c, ok = s.m[platformId]
	s.l.RUnlock()
	return
}

func (s *PlatformToSessions) Do(platformId int, cb func(*SessionToConn)) {
	var sessions *SessionToConn
	s.l.RLock()
	if c, ok := s.m[platformId]; ok {
		sessions = c
		s.l.RUnlock()
		goto end
	}
	s.l.RUnlock()
	return
end:
	cb(sessions)
}

func (s *PlatformToSessions) Add(platformId int, sessions *SessionToConn) (old *SessionToConn) {
	s.l.Lock()
	if c, ok := s.m[platformId]; ok {
		old = c
	}
	s.m[platformId] = sessions
	s.l.Unlock()
	return
}

func (s *PlatformToSessions) Remove(platformId int) (sessions *SessionToConn) {
	_, ok := s.Get(platformId)
	switch ok {
	case true:
		sessions, _ = s.remove(platformId)
	default:
	}
	return
}

func (s *PlatformToSessions) remove(platformId int) (sessions *SessionToConn, ok bool) {
	n := 0
	s.l.Lock()
	sessions, ok = s.m[platformId]
	switch ok {
	case true:
		delete(s.m, platformId)
		n = len(s.m)
		s.l.Unlock()
		goto OK
	}
	s.l.Unlock()
	return
OK:
	logs.Errorf("%v size=%v", platformId, n)
	return
}

func (s *PlatformToSessions) Range(cb func(int, *SessionToConn)) {
	s.l.RLock()
	for platformId, sessions := range s.m {
		cb(platformId, sessions)
	}
	s.l.RUnlock()
}

// <summary>
// [userid][platformid][session]=conn
// <summary>
type UserToPlatforms struct {
	l *sync.RWMutex
	m map[string]*PlatformToSessions
}

func NewUserToPlatforms() *UserToPlatforms {
	return &UserToPlatforms{m: map[string]*PlatformToSessions{}, l: &sync.RWMutex{}}
}

func (s *UserToPlatforms) NumOfLoads() int {
	n := 0
	s.Range(func(userId string, platforms *PlatformToSessions) {
		platforms.Range(func(platformId int, sessions *SessionToConn) {
			n += sessions.Len()
			// sessions.Range(func(sessionId string, _ conn.Session) {
			// })
		})
	})
	return n
}

func (s *UserToPlatforms) Len() int {
	s.l.RLock()
	c := len(s.m)
	s.l.RUnlock()
	return c
}

func (s *UserToPlatforms) Get(userId string) *PlatformToSessions {
	s.l.RLock()
	if c, ok := s.m[userId]; ok {
		s.l.RUnlock()
		return c
	}
	s.l.RUnlock()
	return nil
}

func (s *UserToPlatforms) Do(userId string, cb func(*PlatformToSessions)) {
	var platforms *PlatformToSessions
	s.l.RLock()
	if c, ok := s.m[userId]; ok {
		platforms = c
		s.l.RUnlock()
		goto end
	}
	s.l.RUnlock()
	return
end:
	cb(platforms)
}

func (s *UserToPlatforms) Add(userId string, platforms *PlatformToSessions) (old *PlatformToSessions) {
	s.l.Lock()
	if c, ok := s.m[userId]; ok {
		old = c
	}
	s.m[userId] = platforms
	s.l.Unlock()
	return
}

func (s *UserToPlatforms) Remove(userId string) (platforms *PlatformToSessions) {
	s.l.Lock()
	if c, ok := s.m[userId]; ok {
		platforms = c
		delete(s.m, userId)
	}
	s.l.Unlock()
	return
}

func (s *UserToPlatforms) Range(cb func(string, *PlatformToSessions)) {
	s.l.RLock()
	for userId, platforms := range s.m {
		cb(userId, platforms)
	}
	s.l.RUnlock()
}

func (s *UserToPlatforms) AddUserConn(userId string, platformId int, sessionId string, peer conn.Session) {
	logs.Warnf("userId=%v platformId=%v sessionId=%v", userId, platformId, sessionId)
	logs.Infof("------------------------------- before -------------------------------")
	s.Range(func(userId string, platforms *PlatformToSessions) {
		platforms.Range(func(platformId int, sessions *SessionToConn) {
			sessions.Range(func(sessionId string, _ conn.Session) {
				logs.Infof("userId=%v platformId=%v sessionId=%v", userId, platformId, sessionId)
			})
		})
	})
	if platforms := s.Get(userId); platforms != nil {
		if sessions, _ := platforms.Get(platformId); sessions != nil {
			if sessions.Add(sessionId, peer) != nil {
				panic("error")
			}
		} else {
			sessions := NewSessionToConn()
			if sessions.Add(sessionId, peer) != nil {
				panic("error")
			}
			if platforms.Add(platformId, sessions) != nil {
				panic("error")
			}
		}
	} else {
		platforms := NewPlatformToSessions()
		sessions := NewSessionToConn()
		sessions.Add(sessionId, peer)
		platforms.Add(platformId, sessions)
		if s.Add(userId, platforms) != nil {
			panic("error")
		}
	}
	logs.Infof("------------------------------- after -------------------------------")
	s.Range(func(userId string, platforms *PlatformToSessions) {
		platforms.Range(func(platformId int, sessions *SessionToConn) {
			sessions.Range(func(sessionId string, _ conn.Session) {
				logs.Infof("userId=%v platformId=%v sessionId=%v", userId, platformId, sessionId)
			})
		})
	})
}

func (s *UserToPlatforms) DelUserConn(userId string, platformId int, sessionId string) (peer conn.Session) {
	logs.Warnf("userId=%v platformId=%v sessionId=%v", userId, platformId, sessionId)
	logs.Infof("------------------------------- before -------------------------------")
	s.Range(func(userId string, platforms *PlatformToSessions) {
		platforms.Range(func(platformId int, sessions *SessionToConn) {
			sessions.Range(func(sessionId string, _ conn.Session) {
				logs.Infof("userId=%v platformId=%v sessionId=%v", userId, platformId, sessionId)
			})
		})
	})
	if platforms := s.Get(userId); platforms != nil {
		if sessions, _ := platforms.Get(platformId); sessions != nil {
			peer = sessions.Remove(sessionId)
			if sessions.Len() == 0 {
				platforms.Remove(platformId)
				if platforms.Len() == 0 {
					s.Remove(userId)
					// logs.Infof(" len=%v", s.Len())
				}
			}
		} else {
			// logs.Errorf("get userId=%v platformId=%v failed", userId, platformId)
		}
	} else {
		// logs.Errorf("get userId=%v failed", userId)
	}
	logs.Infof("------------------------------- after -------------------------------")
	s.Range(func(userId string, platforms *PlatformToSessions) {
		platforms.Range(func(platformId int, sessions *SessionToConn) {
			sessions.Range(func(sessionId string, _ conn.Session) {
				logs.Infof("userId=%v platformId=%v sessionId=%v", userId, platformId, sessionId)
			})
		})
	})
	return
}

func (s *UserToPlatforms) GetUserConns(userId string) *PlatformToSessions {
	return s.Get(userId)
}

func (s *UserToPlatforms) GetUserPlatformConns(userId string, platformId int) (c *SessionToConn) {
	if platforms := s.Get(userId); platforms != nil {
		c, _ = platforms.Get(platformId)
	}
	return
}

func (s *UserToPlatforms) GetUserConn(userId string, platformId int, sessionId string) (peer conn.Session) {
	if platforms := s.Get(userId); platforms != nil {
		if sessions, _ := platforms.Get(platformId); sessions != nil {
			peer, _ = sessions.Get(sessionId)
		}
	}
	return
}

// func (s *UserToPlatforms) PrintUserConn(name string) {
// 	logs.Infof("-------------------------------%v-------------------------------", name)
// 	s.Range(func(userId string, platforms *PlatformToSessions) {
// 		platforms.Range(func(platformId int, sessions *SessionToConn) {
// 			sessions.Range(func(sessionId string, _ conn.Session) {
// 				logs.Infof2("userId=%v platformId=%v sessionId=%v", userId, platformId, sessionId)
// 			})
// 		})
// 	})
// }
