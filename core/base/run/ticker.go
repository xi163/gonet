package run

import (
	"errors"
	"time"
)

func NewTicker(d time.Duration) (ticker *time.Ticker, trigger <-chan time.Time) {
	ticker = func(d time.Duration) *time.Ticker {
		if d <= 0 {
			return nil
		}
		return time.NewTicker(d)
	}(d)
	trigger = func(ticker *time.Ticker, d time.Duration) <-chan time.Time {
		if d <= 0 {
			return nil
		}
		return ticker.C
	}(ticker, d)
	return
}

// <summary>
// Trigger 定时器触发器
// <summary>
type Trigger interface {
	Trigger() <-chan time.Time
	Stop()
}

type trigger struct {
	ticker  *time.Ticker
	trigger <-chan time.Time
}

func NewTrigger(d time.Duration) Trigger {
	s := &trigger{}
	s.ticker, s.trigger = NewTicker(d)
	return s
}

func (s *trigger) Stop() {
	if s.ticker == nil {
		panic(errors.New("ticker is nil"))
	}
	s.ticker.Stop()
}

func (s *trigger) Trigger() <-chan time.Time {
	if s.trigger == nil {
		panic(errors.New("trigger is nil"))
	}
	return s.trigger
}
