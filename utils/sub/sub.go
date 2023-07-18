package sub

import (
	"os"
	"sync"

	"github.com/cwloo/gonet/logs"
)

var (
	mgr = NewManager()
)

func Start(name string, args []string, succ func(int, ...any), cb func(*os.ProcessState, ...any), v ...any) (int, bool) {
	return mgr.Start(name, args, succ, cb, v...)
}

func Range(cb func(int, ...any)) {
	mgr.Range(cb)
}

func Kill(pid int) error {
	return mgr.Kill(pid)
}

func KillAll() {
	mgr.KillAll()
}

func WaitAll() {
	mgr.WaitAll()
}

// <summary>
// Sub
// <summary>
type Sub struct {
	p    *os.Process
	args []any
}

// <summary>
// manager
// <summary>
type manager interface {
	Start(name string, args []string, succ func(int, ...any), cb func(*os.ProcessState, ...any), v ...any) (int, bool)
	Range(cb func(int, ...any))
	Kill(pid int) error
	KillAll()
	WaitAll()
}

// <summary>
// pid
// <summary>
type pid struct {
	m  map[int]*Sub
	l  *sync.RWMutex
	wg sync.WaitGroup
}

func NewManager() manager {
	s := &pid{m: map[int]*Sub{},
		l:  &sync.RWMutex{},
		wg: sync.WaitGroup{}}
	return s
}

func (s *pid) Start(name string, args []string, succ func(int, ...any), cb func(*os.ProcessState, ...any), v ...any) (id int, ok bool) {
	attr := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	p, err := os.StartProcess(name, args, attr)
	if err != nil {
		logs.Errorf("%v", err)
		return
	}
	succ(p.Pid, v...)
	s.add(p, v...)
	go s.monitor(p, cb)
	id = p.Pid
	ok = true
	return
}

func (s *pid) add(p *os.Process, v ...any) {
	s.l.Lock()
	_, ok := s.remove(false, p)
	switch ok {
	case true:
		s.l.Unlock()
		goto ERR
	default:
		s.m[p.Pid] = &Sub{p: p, args: v}
		s.l.Unlock()
		goto OK
	}
ERR:
	logs.Fatalf("error")
	return
OK:
	s.wg.Add(1)
}

func (s *pid) remove(lock bool, p *os.Process) (sub *Sub, ok bool) {
	switch lock {
	case true:
		s.l.Lock()
		sub, ok = s.m[p.Pid]
		switch ok {
		case true:
			delete(s.m, p.Pid)
		}
		s.l.Unlock()
	default:
		sub, ok = s.m[p.Pid]
		switch ok {
		case true:
			delete(s.m, p.Pid)
		}
	}
	return
}

func (s *pid) Range(cb func(int, ...any)) {
	s.l.RLock()
	for pid, sub := range s.m {
		cb(pid, sub.args...)
	}
	s.l.RUnlock()
}

func (s *pid) monitor(p *os.Process, cb func(*os.ProcessState, ...any)) {
	sta, err := p.Wait()
	if err != nil {
		logs.Errorf(err.Error())
	}
	if p.Pid != sta.Pid() {
		logs.Fatalf("%v %v", p.Pid, sta.Pid())
	}
	if sta.Success() {
		logs.Debugf("%v exit(%v) succ = %v", sta.Pid(), sta.ExitCode(), sta.String())
	} else {
		logs.Errorf("%v exit(%v) failed = %v", sta.Pid(), sta.ExitCode(), sta.String())
	}
	sub, ok := s.remove(true, p)
	switch ok {
	case true:
		s.wg.Done()
		cb(sta, sub.args...)
	default:
		logs.Fatalf("error")
	}
}

func (s *pid) Kill(pid int) (err error) {
	s.l.RLock()
	p, ok := s.m[pid]
	switch ok {
	case true:
		switch pid == p.p.Pid {
		case true:
			err = p.p.Kill()
			switch err == nil {
			case true:
			default:
				s.l.RUnlock()
				goto ERR
			}
			s.l.RUnlock()
			goto OK
		}
	}
	s.l.RUnlock()
	return
ERR:
	logs.Errorf(err.Error())
	return
OK:
	logs.Errorf("%v", pid)
	return
}

func (s *pid) KillAll() {
	s.l.RLock()
	for _, p := range s.m {
		err := p.p.Kill()
		switch err == nil {
		case true:
		default:
			logs.Errorf(err.Error())
		}
	}
	s.l.RUnlock()
}

func (s *pid) WaitAll() {
	s.wg.Wait()
}
