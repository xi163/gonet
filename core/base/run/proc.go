package run

import (
	"errors"
	"fmt"
	"time"

	"github.com/cwloo/gonet/core/base/run/event"
	"github.com/cwloo/gonet/core/base/timer"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/core/net/conn"
	"github.com/cwloo/gonet/utils/gid"
	"github.com/cwloo/gonet/utils/safe"
)

// <summary>
// Proc 处理单元
// <summary>
type Proc interface {
	cb.Proc
	timer.Proc
	event.Proc
	Tid() int
	Name() string
	Busing() bool
	AssertThis()
	Runner() Processor
	Args() Args
	Duration() time.Duration
	Reset(d time.Duration)
	Do(data any)
	Dispatch(c Proc)
	Dispatcher() Proc
	ResetDispatcher()
	Run()
	Quit()
}

type proc struct {
	tid        int
	name       string
	run        Processor
	args       Args
	dispatcher Proc
}

func NewProc(name string, r Processor) Proc {
	s := &proc{
		tid:  gid.Getgid(),
		name: name,
		run:  r,
	}
	s.assertRunner()
	s.args = s.run.NewArgs(s)
	s.toName()
	return s
}

func (s proc) Name() string {
	return s.name
}

func (s *proc) toName() {
	s.name = s.name + fmt.Sprintf(".proc.%v", s.tid)
}

func (s *proc) Duration() time.Duration {
	s.assertArgs()
	return s.args.Duration()
}

func (s *proc) Reset(d time.Duration) {
	s.assertArgs()
	s.args.Reset(d)
}

func (s *proc) Tid() int {
	if s.tid == 0 {
		panic(errors.New("proc.tid is nil"))
	}
	return s.tid
}

func (s *proc) this() bool {
	return gid.Getgid() == s.tid
}

func (s *proc) AssertThis() {
	if !s.this() {
		panic(fmt.Sprintf("非线程安全 %v", s.tid))
	}
}

func (s *proc) Runner() Processor {
	s.assertRunner()
	return s.run
}

func (s *proc) Args() Args {
	s.assertArgs()
	return s.args
}

func (s *proc) Run() {
	defer safe.Catch()
	s.assertRunner()
	s.assertArgs()
	s.run.Run(s)
	s.run = nil
}

func (s *proc) Busing() bool {
	s.assertArgs()
	return s.args.Busing()
}

func (s proc) assertArgs() {
	if s.args == nil {
		panic(errors.New("proc.args is nil"))
	}
}

func (s proc) assertRunner() {
	if s.run == nil {
		panic(errors.New("proc.run is nil"))
	}
}

func (s *proc) Do(data any) {
	if data != nil {
		s.assertRunner()
		s.run.Queue().Push(data)
	}
}

// s.Exec(func(v any) {
// }, []any{a, b, c})
func (s *proc) Exec(f cb.Functor) {
	if f != nil {
		if s.this() {
			f.Call()
		} else {
			s.Append(f)
		}
	}
}

// s.Append(func(v any) {
// }, []any{a, b, c})
func (s *proc) Append(f cb.Functor) {
	if f != nil {
		s.assertRunner()
		s.run.Queue().Push(f)
	}
}

func (s *proc) RunAfter(delay int32, args ...any) uint32 {
	s.AssertThis()
	s.assertArgs()
	return s.args.RunAfter(delay, args...)
}

func (s *proc) RunAfterWith(delay int32, handler timer.TimerCallback, args ...any) uint32 {
	s.AssertThis()
	s.assertArgs()
	return s.args.RunAfterWith(delay, handler, args...)
}

func (s *proc) RunEvery(delay, interval int32, args ...any) uint32 {
	s.AssertThis()
	s.assertArgs()
	return s.args.RunEvery(delay, interval, args...)
}

func (s *proc) RunEveryWith(delay, interval int32, handler timer.TimerCallback, args ...any) uint32 {
	s.AssertThis()
	s.assertArgs()
	return s.args.RunEveryWith(delay, interval, handler, args...)
}

func (s *proc) RemoveTimer(timerID uint32) {
	s.AssertThis()
	s.assertArgs()
	s.args.RemoveTimer(timerID)
}

func (s *proc) RemoveTimers() {
	s.AssertThis()
	s.assertArgs()
	s.args.RemoveTimers()
}

func (s *proc) Post(data *event.Data) {
	s.Do(data)
}

func (s *proc) PostRead(cmd uint32, msg any, peer conn.Session) {
	s.Post(event.Create(event.EVTRead, event.CreateRead(cmd, msg, peer), nil))
}

func (s *proc) PostReadWith(handler cb.ReadCallback, cmd uint32, msg any, peer conn.Session) {
	s.Post(event.Create(event.EVTRead, event.CreateReadWith(handler, cmd, msg, peer), nil))
}

func (s *proc) PostCustom(cmd uint32, msg any, peer conn.Session) {
	s.Post(event.Create(event.EVTCustom, event.CreateCustom(cmd, msg, peer), nil))
}

func (s *proc) PostCustomWith(handler cb.CustomCallback, cmd uint32, msg any, peer conn.Session) {
	s.Post(event.Create(event.EVTCustom, event.CreateCustomWith(handler, cmd, msg, peer), nil))
}

func (s *proc) PostClosing(d time.Duration, peer conn.Session) {
	s.Post(event.Create(event.EVTClosing, event.CreateClosing(d, peer), nil))
}

func (s *proc) Dispatch(c Proc) {
	s.AssertThis()
	s.dispatcher = c
}

func (s *proc) Dispatcher() Proc {
	s.AssertThis()
	return s.dispatcher
}

func (s *proc) ResetDispatcher() {
	s.AssertThis()
	s.dispatcher = nil
}

func (s *proc) Quit() {
	s.assertArgs()
	s.assertRunner()
	if !s.args.Quit() {
		s.run.Queue().Push(nil)
	}
}
