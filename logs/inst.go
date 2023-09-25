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
func FatalfD(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_DETAIL_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func ErrorfD(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_DETAIL, 4, format, v...)
}

func WarnfD(format string, v ...any) {
	inst.Write("", LVL_WARN, F_DETAIL, 4, format, v...)
}

func InfofD(format string, v ...any) {
	inst.Write("", LVL_INFO, F_DETAIL, 4, format, v...)
}

func TracefD(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_DETAIL, 4, format, v...)
}

func DebugfD(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_DETAIL, 4, format, v...)
}

// F_TMSTMP
func FatalfT(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_TMSTMP_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func ErrorfT(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_TMSTMP, 4, format, v...)
}

func WarnfT(format string, v ...any) {
	inst.Write("", LVL_WARN, F_TMSTMP, 4, format, v...)
}

func InfofT(format string, v ...any) {
	inst.Write("", LVL_INFO, F_TMSTMP, 4, format, v...)
}

func TracefT(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_TMSTMP, 4, format, v...)
}

func DebugfT(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_TMSTMP, 4, format, v...)
}

// F_FN
func FatalfF(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_FN_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func ErrorfF(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_FN, 4, format, v...)
}

func WarnfF(format string, v ...any) {
	inst.Write("", LVL_WARN, F_FN, 4, format, v...)
}

func InfofF(format string, v ...any) {
	inst.Write("", LVL_INFO, F_FN, 4, format, v...)
}

func TracefF(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_FN, 4, format, v...)
}

func DebugfF(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_FN, 4, format, v...)
}

// F_TMSTMP_FN
func FatalfTF(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_TMSTMP_FN_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func ErrorfTF(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_TMSTMP_FN, 4, format, v...)
}

func WarnfTF(format string, v ...any) {
	inst.Write("", LVL_WARN, F_TMSTMP_FN, 4, format, v...)
}

func InfofTF(format string, v ...any) {
	inst.Write("", LVL_INFO, F_TMSTMP_FN, 4, format, v...)
}

func TracefTF(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_TMSTMP_FN, 4, format, v...)
}

func DebugfTF(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_TMSTMP_FN, 4, format, v...)
}

// F_FL
func FatalfL(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_FL_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func ErrorfL(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_FL, 4, format, v...)
}

func WarnfL(format string, v ...any) {
	inst.Write("", LVL_WARN, F_FL, 4, format, v...)
}

func InfofL(format string, v ...any) {
	inst.Write("", LVL_INFO, F_FL, 4, format, v...)
}

func TracefL(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_FL, 4, format, v...)
}

func DebugfL(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_FL, 4, format, v...)
}

// F_TMSTMP_FL
func FatalfTL(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_TMSTMP_FL_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func ErrorfTL(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_TMSTMP_FL, 4, format, v...)
}

func WarnfTL(format string, v ...any) {
	inst.Write("", LVL_WARN, F_TMSTMP_FL, 4, format, v...)
}

func InfofTL(format string, v ...any) {
	inst.Write("", LVL_INFO, F_TMSTMP_FL, 4, format, v...)
}

func TracefTL(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_TMSTMP_FL, 4, format, v...)
}

func DebugfTL(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_TMSTMP_FL, 4, format, v...)
}

// F_FL_FN
func FatalfLF(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_FL_FN_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func ErrorfLF(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_FL_FN, 4, format, v...)
}

func WarnfLF(format string, v ...any) {
	inst.Write("", LVL_WARN, F_FL_FN, 4, format, v...)
}

func InfofLF(format string, v ...any) {
	inst.Write("", LVL_INFO, F_FL_FN, 4, format, v...)
}

func TracefLF(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_FL_FN, 4, format, v...)
}

func DebugfLF(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_FL_FN, 4, format, v...)
}

// F_TMSTMP_FL_FN
func FatalfTLF(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_TMSTMP_FL_FN_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func ErrorfTLF(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_TMSTMP_FL_FN, 4, format, v...)
}

func WarnfTLF(format string, v ...any) {
	inst.Write("", LVL_WARN, F_TMSTMP_FL_FN, 4, format, v...)
}

func InfofTLF(format string, v ...any) {
	inst.Write("", LVL_INFO, F_TMSTMP_FL_FN, 4, format, v...)
}

func TracefTLF(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_TMSTMP_FL_FN, 4, format, v...)
}

func DebugfTLF(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_TMSTMP_FL_FN, 4, format, v...)
}

// F_TEXT
func FatalfTT(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_TEXT_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func ErrorfTT(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_TEXT, 4, format, v...)
}

func WarnfTT(format string, v ...any) {
	inst.Write("", LVL_WARN, F_TEXT, 4, format, v...)
}

func InfofTT(format string, v ...any) {
	inst.Write("", LVL_INFO, F_TEXT, 4, format, v...)
}

func TracefTT(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_TEXT, 4, format, v...)
}

func DebugfTT(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_TEXT, 4, format, v...)
}

// F_PURE
func FatalfP(format string, v ...any) {
	stack := string(debug.Stack())
	inst.Write(stack, LVL_FATAL, F_PURE_SYNC, 4, format, v...)
	inst.Wait()
	panic(stack)
}

func ErrorfP(format string, v ...any) {
	inst.Write("", LVL_ERROR, F_PURE, 4, format, v...)
}

func WarnfP(format string, v ...any) {
	inst.Write("", LVL_WARN, F_PURE, 4, format, v...)
}

func InfofP(format string, v ...any) {
	inst.Write("", LVL_INFO, F_PURE, 4, format, v...)
}

func TracefP(format string, v ...any) {
	inst.Write("", LVL_TRACE, F_PURE, 4, format, v...)
}

func DebugfP(format string, v ...any) {
	inst.Write("", LVL_DEBUG, F_PURE, 4, format, v...)
}
