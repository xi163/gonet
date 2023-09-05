package workers

import (
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/core/base/run/cell"
	"github.com/cwloo/gonet/core/base/timer"
)

// 协程启动参数
type Args struct {
	stopping cc.Singal
	worker   cell.Worker
	ticker   *time.Ticker
	trigger  <-chan time.Time
	timer    timer.ScopedTimer
	// timerv2  *timerv2.SafeTimerScheduel
	timerCb timer.TimerCallback
}

func newArgs(proc run.Proc, d time.Duration, timerCb timer.TimerCallback, creator cell.WorkerCreator, args ...any) run.Args {
	ticker, trigger := run.NewTicker(d)
	s := &Args{
		stopping: cc.NewSingal(),
		ticker:   ticker,
		trigger:  trigger,
		timer:    timer.NewScopedTimer(proc.Tid()),
		// timerv2:  timerv2.NewSafeTimerScheduel(),
		timerCb: timerCb,
	}
	s.worker = creator.Create(proc, args...)
	return s
}

func (s *Args) Quit() bool {
	s.stopping.Signal()
	return true
}

func (s *Args) Duration() (d time.Duration) {
	return
}

func (s *Args) Reset(d time.Duration) {
	s.ticker.Reset(d)
}

func (s *Args) Add(args ...any) {

}

func (s *Args) SetTimerCallback(handler timer.TimerCallback) {
	s.timerCb = handler
}

func (s *Args) Trigger() <-chan time.Time {
	return s.trigger
}

func (s *Args) TimerCallback() (handler timer.TimerCallback) {
	handler = s.timerCb
	return
}

func (s *Args) RunAfter(delay int32, args ...any) uint32 {
	return s.timer.CreateTimer(delay, 0, args...)
}

func (s *Args) RunAfterWith(delay int32, handler timer.TimerCallback, args ...any) uint32 {
	return s.timer.CreateTimerWithCB(delay, 0, handler, args...)
}

func (s *Args) RunEvery(delay, interval int32, args ...any) uint32 {
	return s.timer.CreateTimer(delay, interval, args...)
}

func (s *Args) RunEveryWith(delay, interval int32, handler timer.TimerCallback, args ...any) uint32 {
	return s.timer.CreateTimerWithCB(delay, interval, handler, args...)
}

func (s *Args) RemoveTimer(timerID uint32) {
	s.timer.RemoveTimer(timerID)
}

func (s *Args) RemoveTimers() {
	s.timer.RemoveTimers()
}
