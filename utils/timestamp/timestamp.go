package timestamp

import "time"

type T interface {
	Valid() bool
	Add(sec int32) T
	Less(t T) bool
	Equal(t T) bool
	Greater(t T) bool
	SinceUnixEpoch() int64
}

type tmstamp struct {
	val int64 // second/millisecond/microsecond/nanosecond
}

func New(val int64) T {
	return &tmstamp{val: val}
}

func (s *tmstamp) Valid() bool {
	return s.val > int64(0)
}

func (s *tmstamp) Add(sec int32) T {
	s.val = s.val + int64(sec)
	return s
}

func (s *tmstamp) Less(t T) bool {
	return s.val < t.SinceUnixEpoch()
}

func (s *tmstamp) Equal(t T) bool {
	return s.val == t.SinceUnixEpoch()
}

func (s *tmstamp) Greater(t T) bool {
	return s.val > t.SinceUnixEpoch()
}

func (s *tmstamp) SinceUnixEpoch() int64 {
	return s.val
}

func Add(t T, val int32) T {
	return New(t.SinceUnixEpoch() + int64(val))
}

// 当前时间(秒)
func Now() T {
	return New(time.Now().Unix())
}

// 当前时间(毫秒)
func NowMilliSec() T {
	return New(time.Now().UnixNano() / 1e6)
}

// 当前时间(微秒)
func NowMicroSec() T {
	return New(time.Now().UnixNano() / 1e3)
}

// 当前时间(纳秒)
func NowNanoSec() T {
	return New(time.Now().UnixNano())
}

// 前后间隔时间差(second/millisecond/microsecond/nanosecond)
func Diff(high, low T) int32 {
	diff := int32(high.SinceUnixEpoch() - low.SinceUnixEpoch())
	return diff
}
