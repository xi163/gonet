package gos

import (
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/core/base/timer"
)

// 协程启动参数
type Args struct {
	stopping cc.Singal
}

func newArgs(proc run.Proc) run.Args {
	s := &Args{
		stopping: cc.NewSingal(),
	}
	return s
}

func (s *Args) Quit() bool {
	s.stopping.Signal()
	return true
}

func (s *Args) Trigger() <-chan time.Time {
	return nil
}

func (s *Args) TimerCallback() (handler timer.TimerCallback) {
	return
}

func (s *Args) RunAfter(delay int32, args ...any) uint32 {
	return 0
}

func (s *Args) RunAfterWith(delay int32, handler timer.TimerCallback, args ...any) uint32 {
	return 0
}

func (s *Args) RunEvery(delay, interval int32, args ...any) uint32 {
	return 0
}

func (s *Args) RunEveryWith(delay, interval int32, handler timer.TimerCallback, args ...any) uint32 {
	return 0
}

func (s *Args) RemoveTimer(timerID uint32) {
}

func (s *Args) RemoveTimers() {
}

func (s *Args) Duration() time.Duration {
	return 0
}

func (s *Args) Reset(d time.Duration) {

}

func (s *Args) Add(args ...any) {

}
