package logs

import (
	"runtime/debug"
)

var (
	inst = NewLogger()
)

func SetPrename(name string) {
	inst.SetPrename(name)
}

func GetPrename() string {
	return inst.GetPrename()
}

func TimezoneString() string {
	return inst.TimezoneString()
}

func SetTimezone(timezone Timezone) {
	inst.SetTimezone(timezone)
}

func GetTimezone() Timezone {
	return inst.GetTimezone()
}

func LevelString() string {
	return inst.LevelString()
}

func SetLevel(level Level) {
	inst.SetLevel(level)
}

func GetLevel() Level {
	return inst.GetLevel()
}

func ModeString() string {
	return inst.ModeString()
}

func SetMode(mode Mode) {
	inst.SetMode(mode)
}

func GetMode() Mode {
	return inst.GetMode()
}

func StyleString() string {
	return inst.StyleString()
}

func SetStyle(style Style) {
	inst.SetStyle(style)
}

func GetStyle() Style {
	return inst.GetStyle()
}

func Init(dir string, prename string, logsize int64) {
	inst.Init(dir, prename, logsize)
}

func Close() {
	inst.Close()
}

// F_DETAIL/F_TMSTMP/F_FN/F_TMSTMP_FN/F_FL/F_TMSTMP_FL/F_FL_FN/F_TMSTMP_FL_FN/F_TEXT/F_PURE
func Fatalf(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, inst.GetStyle()|F_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func Errorf(format string, v ...any) {
	inst.Write("", LVL_ERROR, inst.GetStyle(), 4, format, v...)
}

func Warnf(format string, v ...any) {
	inst.Write("", LVL_WARN, inst.GetStyle(), 4, format, v...)
}

func Infof(format string, v ...any) {
	inst.Write("", LVL_INFO, inst.GetStyle(), 4, format, v...)
}

func Tracef(format string, v ...any) {
	inst.Write("", LVL_TRACE, inst.GetStyle(), 4, format, v...)
}

func Debugf(format string, v ...any) {
	inst.Write("", LVL_DEBUG, inst.GetStyle(), 4, format, v...)
}

// F_DETAIL
func Fatalf_detail(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_DETAIL_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func Errorf_detail(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_DETAIL, 4, format, v...)
}

func Warnf_detail(format string, v ...any) {
	inst.Write("", LVL_WARN, F_DETAIL, 4, format, v...)
}

func Infof_detail(format string, v ...any) {
	inst.Write("", LVL_INFO, F_DETAIL, 4, format, v...)
}

func Tracef_detail(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_DETAIL, 4, format, v...)
}

func Debugf_detail(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_DETAIL, 4, format, v...)
}

// F_TMSTMP
func Fatalf_tmsp(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_TMSTMP_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func Errorf_tmsp(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_TMSTMP, 4, format, v...)
}

func Warnf_tmsp(format string, v ...any) {
	inst.Write("", LVL_WARN, F_TMSTMP, 4, format, v...)
}

func Infof_tmsp(format string, v ...any) {
	inst.Write("", LVL_INFO, F_TMSTMP, 4, format, v...)
}

func Tracef_tmsp(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_TMSTMP, 4, format, v...)
}

func Debugf_tmsp(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_TMSTMP, 4, format, v...)
}

// F_FN
func Fatalf_fn(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_FN_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func Errorf_fn(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_FN, 4, format, v...)
}

func Warnf_fn(format string, v ...any) {
	inst.Write("", LVL_WARN, F_FN, 4, format, v...)
}

func Infof_fn(format string, v ...any) {
	inst.Write("", LVL_INFO, F_FN, 4, format, v...)
}

func Tracef_fn(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_FN, 4, format, v...)
}

func Debugf_fn(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_FN, 4, format, v...)
}

// F_TMSTMP_FN
func Fatalf_tmsp_fn(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_TMSTMP_FN_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func Errorf_tmsp_fn(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_TMSTMP_FN, 4, format, v...)
}

func Warnf_tmsp_fn(format string, v ...any) {
	inst.Write("", LVL_WARN, F_TMSTMP_FN, 4, format, v...)
}

func Infof_tmsp_fn(format string, v ...any) {
	inst.Write("", LVL_INFO, F_TMSTMP_FN, 4, format, v...)
}

func Tracef_tmsp_fn(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_TMSTMP_FN, 4, format, v...)
}

func Debugf_tmsp_fn(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_TMSTMP_FN, 4, format, v...)
}

// F_FL
func Fatalf_fl(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_FL_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func Errorf_fl(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_FL, 4, format, v...)
}

func Warnf_fl(format string, v ...any) {
	inst.Write("", LVL_WARN, F_FL, 4, format, v...)
}

func Infof_fl(format string, v ...any) {
	inst.Write("", LVL_INFO, F_FL, 4, format, v...)
}

func Tracef_fl(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_FL, 4, format, v...)
}

func Debugf_fl(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_FL, 4, format, v...)
}

// F_TMSTMP_FL
func Fatalf_tmsp_fl(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_TMSTMP_FL_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func Errorf_tmsp_fl(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_TMSTMP_FL, 4, format, v...)
}

func Warnf_tmsp_fl(format string, v ...any) {
	inst.Write("", LVL_WARN, F_TMSTMP_FL, 4, format, v...)
}

func Infof_tmsp_fl(format string, v ...any) {
	inst.Write("", LVL_INFO, F_TMSTMP_FL, 4, format, v...)
}

func Tracef_tmsp_fl(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_TMSTMP_FL, 4, format, v...)
}

func Debugf_tmsp_fl(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_TMSTMP_FL, 4, format, v...)
}

// F_FL_FN
func Fatalf_fl_fn(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_FL_FN_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func Errorf_fl_fn(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_FL_FN, 4, format, v...)
}

func Warnf_fl_fn(format string, v ...any) {
	inst.Write("", LVL_WARN, F_FL_FN, 4, format, v...)
}

func Infof_fl_fn(format string, v ...any) {
	inst.Write("", LVL_INFO, F_FL_FN, 4, format, v...)
}

func Tracef_fl_fn(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_FL_FN, 4, format, v...)
}

func Debugf_fl_fn(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_FL_FN, 4, format, v...)
}

// F_TMSTMP_FL_FN
func Fatalf_tmsp_fl_fn(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_TMSTMP_FL_FN_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func Errorf_tmsp_fl_fn(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_TMSTMP_FL_FN, 4, format, v...)
}

func Warnf_tmsp_fl_fn(format string, v ...any) {
	inst.Write("", LVL_WARN, F_TMSTMP_FL_FN, 4, format, v...)
}

func Infof_tmsp_fl_fn(format string, v ...any) {
	inst.Write("", LVL_INFO, F_TMSTMP_FL_FN, 4, format, v...)
}

func Tracef_tmsp_fl_fn(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_TMSTMP_FL_FN, 4, format, v...)
}

func Debugf_tmsp_fl_fn(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_TMSTMP_FL_FN, 4, format, v...)
}

// F_TEXT
func Fatalf_text(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_TEXT_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func Errorf_text(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_TEXT, 4, format, v...)
}

func Warnf_text(format string, v ...any) {
	inst.Write("", LVL_WARN, F_TEXT, 4, format, v...)
}

func Infof_text(format string, v ...any) {
	inst.Write("", LVL_INFO, F_TEXT, 4, format, v...)
}

func Tracef_text(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_TEXT, 4, format, v...)
}

func Debugf_text(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_TEXT, 4, format, v...)
}

// F_PURE
func Fatalf_pure(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_PURE_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func Errorf_pure(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_PURE, 4, format, v...)
}

func Warnf_pure(format string, v ...any) {
	inst.Write("", LVL_WARN, F_PURE, 4, format, v...)
}

func Infof_pure(format string, v ...any) {
	inst.Write("", LVL_INFO, F_PURE, 4, format, v...)
}

func Tracef_pure(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_PURE, 4, format, v...)
}

func Debugf_pure(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_PURE, 4, format, v...)
}
