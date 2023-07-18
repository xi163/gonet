package logs

import "strings"

func SprintErrorf(format string, v ...any) string {
	prefix, content := inst.Sprint(LVL_ERROR, F_FL_FN, 3, format, v...)
	return strings.Join([]string{prefix, content}, "")
}

func SprintWarnf(format string, v ...any) string {
	prefix, content := inst.Sprint(LVL_WARN, F_FL_FN, 3, format, v...)
	return strings.Join([]string{prefix, content}, "")
}

func SprintInfof(format string, v ...any) string {
	prefix, content := inst.Sprint(LVL_INFO, F_FL_FN, 3, format, v...)
	return strings.Join([]string{prefix, content}, "")
}

func SprintTracef(format string, v ...any) string {
	prefix, content := inst.Sprint(LVL_TRACE, F_FL_FN, 3, format, v...)
	return strings.Join([]string{prefix, content}, "")
}

func SprintDebugf(format string, v ...any) string {
	prefix, content := inst.Sprint(LVL_DEBUG, F_FL_FN, 3, format, v...)
	return strings.Join([]string{prefix, content}, "")
}
