package logs

import "strings"

func SprintErrorf(skip int, format string, v ...any) string {
	prefix, content := inst.Sprint(LVL_ERROR, F_FL_FN, skip, format, v...)
	idx := strings.Index(prefix, " ")
	return strings.Join([]string{prefix[idx+1:], content}, "")
}

func SprintWarnf(skip int, format string, v ...any) string {
	prefix, content := inst.Sprint(LVL_WARN, F_FL_FN, skip, format, v...)
	idx := strings.Index(prefix, " ")
	return strings.Join([]string{prefix[idx+1:], content}, "")
}

func SprintInfof(skip int, format string, v ...any) string {
	prefix, content := inst.Sprint(LVL_INFO, F_FL_FN, skip, format, v...)
	idx := strings.Index(prefix, " ")
	return strings.Join([]string{prefix[idx+1:], content}, "")
}

func SprintTracef(skip int, format string, v ...any) string {
	prefix, content := inst.Sprint(LVL_TRACE, F_FL_FN, skip, format, v...)
	idx := strings.Index(prefix, " ")
	return strings.Join([]string{prefix[idx+1:], content}, "")
}

func SprintDebugf(skip int, format string, v ...any) string {
	prefix, content := inst.Sprint(LVL_DEBUG, F_FL_FN, skip, format, v...)
	idx := strings.Index(prefix, " ")
	return strings.Join([]string{prefix[idx+1:], content}, "")
}
