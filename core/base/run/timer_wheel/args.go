package timer_wheel

import (
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/core/base/timer"
	"github.com/cwloo/gonet/utils/timerwheel"
)

// RunAfter(delay int32, args ...any) uint32
// RunAfterWith(delay int32, handler timer.TimerCallback, args ...any) uint32
// RunEvery(delay, interval int32, args ...any) uint32
// RunEveryWith(delay, interval int32, handler timer.TimerCallback, args ...any) uint32
// RemoveTimer(timerID uint32)
// RemoveTimers()

// <summary>
// Args 协程启动参数
// <summary>
type Args struct {
	using    bool
	stopping cc.Singal
	ticker   *time.Ticker
	trigger  <-chan time.Time
	timer    timer.ScopedTimer
	// timerv2    *timerv2.SafeTimerScheduel
	timerCb    timer.TimerCallback
	timerWheel timerwheel.TimerWheel
}

func newArgs(proc run.Proc, size int32, d time.Duration, timerCb timer.TimerCallback) run.Args {
	ticker, trigger := run.NewTicker(d)
	s := &Args{
		stopping: cc.NewSingal(),
		ticker:   ticker,
		trigger:  trigger,
		timer:    timer.NewScopedTimer(proc.Tid()),
		// timerv2:    timerv2.NewSafeTimerScheduel(),
		timerCb:    timerCb,
		timerWheel: timerwheel.NewTimerWheel(proc.Tid(), size),
	}
	return s
}

func (s *Args) SetUsing(b bool) {
	s.using = b
}

func (s *Args) GetUsing() bool {
	return s.using
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

func (s *Args) PopBucket(interval int32) (v []any) {
	v = s.timerWheel.PopBucket(interval)
	return
}

func (s *Args) PushBucket(val any, timeout int32) int32 {
	return s.timerWheel.PushBucket(val, timeout)
}

func (s *Args) UpdateBucket(val any, cursor int32, timeout int32) int32 {
	return s.timerWheel.UpdateBucket(val, cursor, timeout)
}
