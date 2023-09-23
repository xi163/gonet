package logs

import (
	"sync/atomic"
)

type unsafeArg struct {
	// byte alignment unsafe.Offsetof
	timezone int32
	level    int32
	mode     int32
	style    int32
}

func newUnsafeArg() *unsafeArg {
	s := &unsafeArg{}
	s.setTimezone(MY_CST)
	s.setLevel(LVL_DEBUG)
	s.setMode(M_STDOUT_FILE)
	s.setStyle(F_DETAIL)
	return s
}

func (s *unsafeArg) timezoneString() string {
	switch Timezone(s.timezone) {
	case MY_PST:
		return "MY_PST"
	case MY_MST:
		return "MY_MST"
	case MY_EST:
		return "MY_EST"
	case MY_BST:
		return "MY_BST"
	case MY_UTC:
		return "MY_UTC"
	case MY_GST:
		return "MY_GST"
	case MY_CST:
		return "MY_CST"
	case MY_JST:
		return "MY_JST"
	default:
		return "Timezone[unknow]"
	}
}

func (s *unsafeArg) getTimezone() Timezone {
	return Timezone(atomic.LoadInt32(&s.timezone))
}

func (s *unsafeArg) setTimezone(timezone Timezone) bool {
	switch timezone {
	case MY_PST, MY_MST, MY_EST, MY_BST, MY_UTC, MY_GST, MY_CST, MY_JST:
		atomic.StoreInt32(&s.timezone, int32(timezone))
	default:
		return false
	}
	return true
}

func (s *unsafeArg) modeString() string {
	return MODE[atomic.LoadInt32(&s.mode)]
}

func (s *unsafeArg) getMode() Mode {
	return Mode(atomic.LoadInt32(&s.mode))
}

func (s *unsafeArg) setMode(mode Mode) bool {
	switch mode {
	case M_STDOUT_ONLY, M_FILE_ONLY, M_STDOUT_FILE:
		atomic.StoreInt32(&s.mode, int32(mode))
	default:
		return false
	}
	return true
}

func (s *unsafeArg) styleString() string {
	switch Style(s.style) {
	case F_DETAIL:
		return "F_DETAIL"
	case F_TMSTMP:
		return "F_TMSTMP"
	case F_FN:
		return "F_FN"
	case F_TMSTMP_FN:
		return "F_TMSTMP_FN"
	case F_FL:
		return "F_FL"
	case F_TMSTMP_FL:
		return "F_TMSTMP_FL"
	case F_FL_FN:
		return "F_FL_FN"
	case F_TMSTMP_FL_FN:
		return "F_TMSTMP_FL_FN"
	case F_TEXT:
		return "F_TEXT"
	case F_PURE:
		return "F_PURE"
	default:
		return "F_UNKNOWN"
	}
}

func (s *unsafeArg) getStyle() (style Style) {
	style = Style(atomic.LoadInt32(&s.style))
	return
}

func (s *unsafeArg) setStyle(style Style) bool {
	switch style {
	case F_DETAIL: //F_DETAIL
		atomic.StoreInt32(&s.style, int32(F_DETAIL))
	case F_TMSTMP: //F_TMSTMP
		atomic.StoreInt32(&s.style, int32(F_TMSTMP))
	case F_FN: //F_FN
		atomic.StoreInt32(&s.style, int32(F_FN))
	case F_TMSTMP_FN: //F_TMSTMP_FN
		atomic.StoreInt32(&s.style, int32(F_TMSTMP_FN))
	case F_FL: //F_FL
		atomic.StoreInt32(&s.style, int32(F_FL))
	case F_TMSTMP_FL: //F_TMSTMP_FL
		atomic.StoreInt32(&s.style, int32(F_TMSTMP_FL))
	case F_FL_FN: //F_FL_FN
		atomic.StoreInt32(&s.style, int32(F_FL_FN))
	case F_TMSTMP_FL_FN: //F_TMSTMP_FL_FN
		atomic.StoreInt32(&s.style, int32(F_TMSTMP_FL_FN))
	case F_TEXT: //F_TEXT
		atomic.StoreInt32(&s.style, int32(F_TEXT))
	case F_PURE: //F_PURE
		atomic.StoreInt32(&s.style, int32(F_PURE))
	default:
		return false
	}
	return true
}

func (s *unsafeArg) levelString() string {
	return LVL[atomic.LoadInt32(&s.level)]
}

func (s *unsafeArg) getLevel() Level {
	return Level(atomic.LoadInt32(&s.level))
}

func (s *unsafeArg) setLevel(level Level) bool {
	switch level {
	case LVL_DEBUG, LVL_TRACE, LVL_INFO, LVL_WARN, LVL_ERROR, LVL_FATAL:
		atomic.StoreInt32(&s.level, int32(level))
	default:
		return false
	}
	return true
}
