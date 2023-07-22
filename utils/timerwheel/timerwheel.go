package timerwheel

import (
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils/bucket"
	"github.com/cwloo/gonet/utils/circular"
	"github.com/cwloo/gonet/utils/gid"
)

// <summary>
// TimerWheel 时间轮盘，处理空闲超时连接
// <summary>
type TimerWheel interface {
	PopBucket(interval int32) (v []any)
	PushBucket(val any, timeout int32) int32
	UpdateBucket(val any, cursor int32, timeout int32) int32
}

type timerWheel struct {
	tid  int
	ring circular.Buffer[*bucket.Bucket]
}

// 轮盘大小(size) >=
// 空闲超时时间(timeout) >
// 心跳间隔时间(interval)
// ----------------------------------------------------------
func NewTimerWheel(tid int, size int32) TimerWheel {
	s := &timerWheel{tid: tid, ring: circular.NewWitch[*bucket.Bucket](int(size), bucket.NewBucket)}
	s.ring.Resize(int(size))
	return s
}

func (s *timerWheel) PopBucket(interval int32) (v []any) {
	s.assertThisThread()
	v = s.ring.Front().Pop()
	s.ring.PushBack(bucket.NewBucket())
	// str := ""
	// for i := 0; i < s.ring.Size(); i++ {
	// 	str += "\n" + fmt.Sprintf("cap:%d size:%d begin:%d end:%d At(%d):%d",
	// 		s.ring.Capacity(), s.ring.Size(),
	// 		s.ring.Begin(), s.ring.End(),
	// 		i, s.ring.At(i).Len())
	// }
	// logs.Infof(str)
	return
}

func (s *timerWheel) PushBucket(val any, timeout int32) int32 {
	s.assertThisThread()
	s.ring.Back().Add(val)
	// str := ""
	// for i := 0; i < s.ring.Size(); i++ {
	// 	str += "\n" + fmt.Sprintf("cap:%d size:%d begin:%d end:%d At(%d):%d",
	// 		s.ring.Capacity(), s.ring.Size(),
	// 		s.ring.Begin(), s.ring.End(),
	// 		i, s.ring.At(i).Len())
	// }
	// logs.Debugf(str)
	return int32(s.ring.End())
}

func (s *timerWheel) UpdateBucket(val any, cursor int32, timeout int32) int32 {
	s.assertThisThread()
	s.ring.Range(func(bucket *bucket.Bucket) bool {
		return bucket.Remove(val)
	})
	s.ring.Back().Add(val)
	// str := ""
	// for i := 0; i < s.ring.Size(); i++ {
	// 	str += "\n" + fmt.Sprintf("cap:%d size:%d begin:%d end:%d At(%d):%d",
	// 		s.ring.Capacity(), s.ring.Size(),
	// 		s.ring.Begin(), s.ring.End(),
	// 		i, s.ring.At(i).Len())
	// }
	// logs.Warnf(str)
	return int32(s.ring.End())
}

func (s *timerWheel) this() bool {
	return gid.Getgid() == s.tid
}

func (s *timerWheel) assertThisThread() {
	if !s.this() {
		panic(logs.SprintErrorf("非线程安全 %v", s.tid))
	}
}
