package ch

import "github.com/cwloo/gonet/core/base/mq"

// chan消息队列
type Queue interface {
	mq.Queue
	Read() <-chan any
	Signal() <-chan bool
	Full() bool
	Length() int
	AssertEmpty()
	Reset()
}
