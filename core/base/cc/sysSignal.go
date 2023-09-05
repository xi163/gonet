package cc

import (
	"os"
	"os/signal"
	"sync"
)

// 系统中断信号
type SysSignal interface {
	Start(handler func())
	Wait()
	WaitSignal()
	Stop()
}

// 系统中断信号
type sysSignal struct {
	watching bool
	lock     *sync.Mutex
	cond     *sync.Cond
	ch       chan os.Signal
	done     chan os.Signal
	flag     [2]AtomFlag
	handler  func()
}

func NewSysSignal() SysSignal {
	s := &sysSignal{
		lock: &sync.Mutex{},
		flag: [2]AtomFlag{
			NewAtomFlag(),
			NewAtomFlag()},
	}
	s.cond = sync.NewCond(s.lock)
	return s
}

func (s *sysSignal) Start(handler func()) {
	if !s.watching && s.flag[0].TestSet() {
		s.handler = handler
		s.ch = make(chan os.Signal)
		s.done = make(chan os.Signal)
		signal.Notify(s.ch, os.Interrupt, os.Kill)
		go s.watch()
		s.wait()
		s.flag[0].Reset()
	}
}

func (s *sysSignal) wait() {
	s.lock.Lock()
	for !s.watching {
		s.cond.Wait()
	}
	s.lock.Unlock()
}

func (s *sysSignal) Wait() {
	s.lock.Lock()
	for s.watching {
		s.cond.Wait()
	}
	s.lock.Unlock()
}

func (s *sysSignal) watch() {
	s.lock.Lock()
	s.watching = true
	s.cond.Signal()
	s.lock.Unlock()

	sig := <-s.ch
	close(s.ch)
	s.done <- sig

	s.lock.Lock()
	s.watching = false
	s.cond.Signal()
	s.lock.Unlock()
}

func (s *sysSignal) WaitSignal() {
	if s.watching {
		<-s.done
		//Stop()或CTRL+C前执行清理
		if s.handler != nil {
			s.handler()
		}
		close(s.done)
	}
}

func (s *sysSignal) stop() {
	//通知监视器退出
	s.ch <- os.Interrupt
}

func (s *sysSignal) Stop() {
	if s.watching && s.flag[1].TestSet() {
		s.stop()
		s.Wait()
		s.flag[1].Reset()
	}
}
