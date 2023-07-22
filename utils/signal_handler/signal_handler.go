package signal_handler

import (
	"os"
	"os/signal"
	"syscall"
)

var sig Handler

func Wait(f func()) {
	sig = newHandler(f)
	go func() {
		<-sig.Signal()
		sig.Done()
	}()
	sig.Wait()
}

func Stop() {
	sig.Stop()
}

type Handler interface {
	Signal() <-chan os.Signal
	Stop()
	Done()
	Wait()
}

type handler struct {
	sig  chan os.Signal
	done chan bool
	f    func()
}

func newHandler(f func()) Handler {
	s := &handler{
		sig:  make(chan os.Signal, 1),
		done: make(chan bool, 1),
		f:    f}
	// SIGINT SIGTERM SIGKILL SIGHUP
	signal.Notify(s.sig, syscall.SIGINT, syscall.SIGTERM)
	return s
}

func (s *handler) Signal() <-chan os.Signal {
	return s.sig
}

func (s *handler) Stop() {
	s.sig <- syscall.SIGINT
	// syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}

func (s *handler) Done() {
	close(s.sig)
	s.done <- true
}

func (s *handler) Wait() {
	<-s.done
	s.f()
	close(s.done)
}
