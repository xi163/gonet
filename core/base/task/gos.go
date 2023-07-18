package task

import (
	"sync"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/mq/lq"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/core/base/run/gos"
	"github.com/cwloo/gonet/core/base/watcher"
	"github.com/cwloo/gonet/core/cb"
)

func NewGos(name string, init, size int, fixed, nonblock bool, handler cb.Processor) Task {
	s := &task{
		name:     name,
		init:     init,
		size:     size,
		fixed:    fixed,
		nonblock: nonblock,
		i32:      cc.NewI32(),
		lock:     &sync.Mutex{},
		slots:    map[int32]run.Slot{},
		watcher:  watcher.NewWatcher(name, lq.NewQueue(0)),
		run:      gos.NewProcessor(handler),
		flag: [2]cc.AtomFlag{
			cc.NewAtomFlag(),
			cc.NewAtomFlag()},
	}
	// if New == nil {
	// 	s.New = s.New_chmq
	// } else {
	// 	s.New = New
	// }
	// if overload == nil {
	// 	s.overload = s.Overload
	// } else {
	// 	s.overload = overload
	// }
	// if gcCondition != nil {
	// 	s.run.SetGcCondition(gcCondition)
	// }
	// ptr := utils.IF(New == nil, s.New_chmq, New)
	// unsafe.Pointer(&ptr)
	// s.New = utils.IF(New == nil, s.New_chmq, New).(mq.New)
	// s.overload = utils.IF(overload == nil, s.Overload, overload).(run.Overload)
	return s
}
