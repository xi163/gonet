package Signal

import (
	"os"
	"os/signal"
)

var sig = newHandler()

func RegisterStopFunc(f func()) {
	go func() {
		for {
			select {
			case <-sig.stopSignal():
				f()
			}
		}
	}()
}

func Stop() {
	// syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	sig.stop()
}

// const (
// 	SIGUSR1 = os.Signal(0xa)
// 	SIGUSR2 = os.Signal(0xc)
// )

type handler struct {
	chStop chan os.Signal
}

func newHandler() *handler {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	//signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP)
	return &handler{chStop: stop}
}

func (s *handler) stopSignal() <-chan os.Signal {
	return s.chStop
}

func (s *handler) stop() {
	s.chStop <- os.Interrupt
	// s.chStop <- syscall.SIGINT
	//syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}
