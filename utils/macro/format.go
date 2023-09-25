package macro

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/cwloo/gonet/utils/Fn"
	"github.com/cwloo/gonet/utils/gid"
)

var (
	PID  = os.Getpid()
	TAG  = []byte{'T', 'P'}
	CHR  = []string{"F", "E", "W", "I", "T", "D"}
	LVL  = []string{"FATAL", "ERROR", "WARNING", "INFO", "TRACE", "DEBUG"}
	MODE = []string{"M_STDOUT_ONLY", "M_FILE_ONLY", "M_STDOUT_FILE"}
)

func Sprint(ok bool, tm *time.Time, timezone Timezone, pid int, name func(bool) string, level Level, style Style, skip int, format string, v ...any) (prefix, content string) {
	prefix = Format(ok, tm, timezone, pid, name, level, style, skip)
	content = fmt.Sprintf(format, v...)
	return
}

func Format(ok bool, tm *time.Time, timezone Timezone, pid int, name func(bool) string, level Level, style Style, skip int) (prefix string) {
	// 2006/01/02 15:04:05.000000
	dt := tm.Format("15:04:05.000000")
	tid := gid.Getgid()
	switch style {
	case F_DETAIL, F_DETAIL_SYNC:
		//W101106 CST 21:17:00.024254 199 main.go:103][main] server.run xxx
		pc, f, line, _ := runtime.Caller(skip)
		_, file := path.Split(f)
		pg, fn := Fn.Split(runtime.FuncForPC(pc).Name())
		var b strings.Builder
		switch ok {
		case true:
			b.WriteByte(TAG[0])
		default:
			b.WriteByte(TAG[1])
		}
		b.WriteString(CHR[level])
		b.WriteString(strconv.Itoa(pid))
		b.WriteString(name(false))
		b.WriteString(" ")
		switch ok {
		case true:
			b.WriteString(String(timezone))
			b.WriteString(" ")
			b.WriteString(dt)
			b.WriteString(" ")
		}
		b.WriteString(strconv.Itoa(tid))
		b.WriteString(" ")
		b.WriteString(file)
		b.WriteString(":")
		b.WriteString(strconv.Itoa(line))
		b.WriteString("][")
		b.WriteString(pg)
		b.WriteString("] ")
		b.WriteString(fn)
		b.WriteString(" ")
		prefix = b.String()
	case F_TMSTMP, F_TMSTMP_SYNC:
		//W101106 CST 21:17:00.024254] xxx
		var b strings.Builder
		switch ok {
		case true:
			b.WriteByte(TAG[0])
		default:
			b.WriteByte(TAG[1])
		}
		b.WriteString(CHR[level])
		b.WriteString(strconv.Itoa(pid))
		b.WriteString(name(false))
		switch ok {
		case true:
			b.WriteString(" ")
			b.WriteString(String(timezone))
			b.WriteString(" ")
			b.WriteString(dt)
		}
		b.WriteString("] ")
		prefix = b.String()
	case F_FN, F_FN_SYNC:
		//W101106][main] server.run xxx
		pc, _, _, _ := runtime.Caller(skip)
		pg, fn := Fn.Split(runtime.FuncForPC(pc).Name())
		var b strings.Builder
		switch ok {
		case true:
			b.WriteByte(TAG[0])
		default:
			b.WriteByte(TAG[1])
		}
		b.WriteString(CHR[level])
		b.WriteString(strconv.Itoa(pid))
		b.WriteString(name(false))
		b.WriteString("][")
		b.WriteString(pg)
		b.WriteString("] ")
		b.WriteString(fn)
		b.WriteString(" ")
		prefix = b.String()
	case F_TMSTMP_FN, F_TMSTMP_FN_SYNC:
		//W101106 CST 21:17:00.024254][main] server.run xxx
		pc, _, _, _ := runtime.Caller(skip)
		pg, fn := Fn.Split(runtime.FuncForPC(pc).Name())
		var b strings.Builder
		switch ok {
		case true:
			b.WriteByte(TAG[0])
		default:
			b.WriteByte(TAG[1])
		}
		b.WriteString(CHR[level])
		b.WriteString(strconv.Itoa(pid))
		b.WriteString(name(false))
		switch ok {
		case true:
			b.WriteString(" ")
			b.WriteString(String(timezone))
			b.WriteString(" ")
			b.WriteString(dt)
		}
		b.WriteString("][")
		b.WriteString(pg)
		b.WriteString("] ")
		b.WriteString(fn)
		b.WriteString(" ")
		prefix = b.String()
	case F_FL, F_FL_SYNC:
		//W101106 main.go:103] xxx
		_, f, line, _ := runtime.Caller(skip)
		_, file := path.Split(f)
		var b strings.Builder
		switch ok {
		case true:
			b.WriteByte(TAG[0])
		default:
			b.WriteByte(TAG[1])
		}
		b.WriteString(CHR[level])
		b.WriteString(strconv.Itoa(pid))
		b.WriteString(name(false))
		b.WriteString(" ")
		b.WriteString(file)
		b.WriteString(":")
		b.WriteString(strconv.Itoa(line))
		b.WriteString("] ")
		prefix = b.String()
	case F_TMSTMP_FL, F_TMSTMP_FL_SYNC:
		//W101106 CST 21:17:00.024254 main.go:103] xxx
		_, f, line, _ := runtime.Caller(skip)
		_, file := path.Split(f)
		var b strings.Builder
		switch ok {
		case true:
			b.WriteByte(TAG[0])
		default:
			b.WriteByte(TAG[1])
		}
		b.WriteString(CHR[level])
		b.WriteString(strconv.Itoa(pid))
		b.WriteString(name(false))
		b.WriteString(" ")
		switch ok {
		case true:
			b.WriteString(String(timezone))
			b.WriteString(" ")
			b.WriteString(dt)
			b.WriteString(" ")
		}
		b.WriteString(file)
		b.WriteString(":")
		b.WriteString(strconv.Itoa(line))
		b.WriteString("] ")
		prefix = b.String()
	case F_FL_FN, F_FL_FN_SYNC:
		//W101106 main.go:103][main] server.run xxx
		pc, f, line, _ := runtime.Caller(skip)
		_, file := path.Split(f)
		pg, fn := Fn.Split(runtime.FuncForPC(pc).Name())
		var b strings.Builder
		switch ok {
		case true:
			b.WriteByte(TAG[0])
		default:
			b.WriteByte(TAG[1])
		}
		b.WriteString(CHR[level])
		b.WriteString(strconv.Itoa(pid))
		b.WriteString(name(false))
		b.WriteString(" ")
		b.WriteString(file)
		b.WriteString(":")
		b.WriteString(strconv.Itoa(line))
		b.WriteString("][")
		b.WriteString(pg)
		b.WriteString("] ")
		b.WriteString(fn)
		b.WriteString(" ")
		prefix = b.String()
	case F_TMSTMP_FL_FN, F_TMSTMP_FL_FN_SYNC:
		//W101106 CST 21:17:00.024254 main.go:103][main] server.run xxx
		pc, f, line, _ := runtime.Caller(skip)
		_, file := path.Split(f)
		pg, fn := Fn.Split(runtime.FuncForPC(pc).Name())
		var b strings.Builder
		switch ok {
		case true:
			b.WriteByte(TAG[0])
		default:
			b.WriteByte(TAG[1])
		}
		b.WriteString(CHR[level])
		b.WriteString(strconv.Itoa(pid))
		b.WriteString(name(false))
		b.WriteString(" ")
		switch ok {
		case true:
			b.WriteString(String(timezone))
			b.WriteString(" ")
			b.WriteString(dt)
			b.WriteString(" ")
		}
		b.WriteString(file)
		b.WriteString(":")
		b.WriteString(strconv.Itoa(line))
		b.WriteString("][")
		b.WriteString(pg)
		b.WriteString("] ")
		b.WriteString(fn)
		b.WriteString(" ")
		prefix = b.String()
	case F_TEXT, F_TEXT_SYNC:
		//W101106] xxx
		var b strings.Builder
		switch ok {
		case true:
			b.WriteByte(TAG[0])
		default:
			b.WriteByte(TAG[1])
		}
		b.WriteString(CHR[level])
		b.WriteString(strconv.Itoa(pid))
		b.WriteString(name(false))
		b.WriteString("] ")
		prefix = b.String()
	case F_PURE, F_PURE_SYNC:
		fallthrough
	default:
		//xxx
		var b bytes.Buffer
		switch ok {
		case true:
			b.WriteByte(TAG[0])
		default:
			b.WriteByte(TAG[1])
		}
		b.WriteString(CHR[level])
		b.WriteString(name(true))
		prefix = b.String()
	}
	return
}
