package timer

// <summary>
// Proc 定时器处理单元
// <summary>
type Proc interface {
	RunAfter(delay int32, args ...any) uint32
	RunAfterWith(delay int32, handler TimerCallback, args ...any) uint32
	RunEvery(delay, interval int32, args ...any) uint32
	RunEveryWith(delay, interval int32, handler TimerCallback, args ...any) uint32
	RemoveTimer(timerID uint32)
	RemoveTimers()
}
