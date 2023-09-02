package run

import (
	"time"

	"github.com/cwloo/gonet/core/base/timer"
)

// <summary>
// Args 协程启动参数
// <summary>
type Args interface {
	RunAfter(delay int32, args ...any) uint32
	RunAfterWith(delay int32, handler timer.TimerCallback, args ...any) uint32
	RunEvery(delay, interval int32, args ...any) uint32
	RunEveryWith(delay, interval int32, handler timer.TimerCallback, args ...any) uint32
	RemoveTimer(timerID uint32)
	RemoveTimers()
	Duration() time.Duration
	Reset(d time.Duration)
	Add(args ...any)
	Quit() bool
	Trigger() <-chan time.Time
	TimerCallback() (handler timer.TimerCallback)
}
