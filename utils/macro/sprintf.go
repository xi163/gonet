package macro

import (
	"strings"
	"time"
)

func name(space bool) string {
	return ""
}

func SprintErrorf(skip int, format string, v ...any) string {
	timezone := MY_CST
	var tm time.Time
	t := time.Now()
	ok := convertUTC(&t, &tm, timezone)
	prefix, content := Sprint(ok, &tm, timezone, PID, name, LVL_ERROR, F_FL_FN, skip, format, v...)
	idx := strings.Index(prefix, " ")
	return strings.Join([]string{prefix[idx+1:], content}, "")
}

func SprintWarnf(skip int, format string, v ...any) string {
	timezone := MY_CST
	var tm time.Time
	t := time.Now()
	ok := convertUTC(&t, &tm, timezone)
	prefix, content := Sprint(ok, &tm, timezone, PID, name, LVL_WARN, F_FL_FN, skip, format, v...)
	idx := strings.Index(prefix, " ")
	return strings.Join([]string{prefix[idx+1:], content}, "")
}

func SprintCriticalf(skip int, format string, v ...any) string {
	timezone := MY_CST
	var tm time.Time
	t := time.Now()
	ok := convertUTC(&t, &tm, timezone)
	prefix, content := Sprint(ok, &tm, timezone, PID, name, LVL_CRITICAL, F_FL_FN, skip, format, v...)
	idx := strings.Index(prefix, " ")
	return strings.Join([]string{prefix[idx+1:], content}, "")
}

func SprintInfof(skip int, format string, v ...any) string {
	timezone := MY_CST
	var tm time.Time
	t := time.Now()
	ok := convertUTC(&t, &tm, timezone)
	prefix, content := Sprint(ok, &tm, timezone, PID, name, LVL_INFO, F_FL_FN, skip, format, v...)
	idx := strings.Index(prefix, " ")
	return strings.Join([]string{prefix[idx+1:], content}, "")
}

func SprintDebugf(skip int, format string, v ...any) string {
	timezone := MY_CST
	var tm time.Time
	t := time.Now()
	ok := convertUTC(&t, &tm, timezone)
	prefix, content := Sprint(ok, &tm, timezone, PID, name, LVL_DEBUG, F_FL_FN, skip, format, v...)
	idx := strings.Index(prefix, " ")
	return strings.Join([]string{prefix[idx+1:], content}, "")
}

func SprintTracef(skip int, format string, v ...any) string {
	timezone := MY_CST
	var tm time.Time
	t := time.Now()
	ok := convertUTC(&t, &tm, timezone)
	prefix, content := Sprint(ok, &tm, timezone, PID, name, LVL_TRACE, F_FL_FN, skip, format, v...)
	idx := strings.Index(prefix, " ")
	return strings.Join([]string{prefix[idx+1:], content}, "")
}
