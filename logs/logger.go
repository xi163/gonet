package logs

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/mq/lq"
	"github.com/cwloo/gonet/core/base/pipe"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/utils/Fn"
	"github.com/cwloo/gonet/utils/conv"
	"github.com/cwloo/gonet/utils/gid"
)

var (
	i32  = cc.NewI32()
	TAG  = []byte{'T', 'P'}
	CHR  = []string{"F", "E", "W", "I", "T", "D"}
	LVL  = []string{"FATAL", "ERROR", "WARNING", "INFO", "TRACE", "DEBUG"}
	MODE = []string{"M_STDOUT_ONLY", "M_FILE_ONLY", "M_STDOUT_FILE"}
	bio  = 0
)

// 异步日志系统
type Logger interface {
	SetPrename(name string)
	GetPrename() string
	TimezoneString() string
	GetTimezone() Timezone
	SetTimezone(timezone Timezone)
	ModeString() string
	GetMode() Mode
	SetMode(mode Mode)
	StyleString() string
	GetStyle() Style
	SetStyle(style Style)
	LevelString() string
	GetLevel() Level
	SetLevel(level Level)
	SetColor(level Level, prefix, context int)
	Init(dir string, prename string, logsize int64)
	Sprint(level Level, style Style, skip int, format string, v ...any) (string, string)
	Write(stack string, level Level, style Style, skip int, format string, v ...any)
	Wait()
	Close()
}

// logger 异步日志系统

type logger struct {
	utcOk   bool
	mkdir   bool
	sync    bool
	pid     int
	day     int
	size    int64
	prename string
	prefix  string
	path    string
	fd      *os.File
	tm      time.Time
	arg     *unsafeArg
	pipe    pipe.Pipe
	l       *sync.RWMutex
	bio     *bufio.Writer
	l_sync  *sync.Mutex
	cond    *sync.Cond
	flag    cc.AtomFlag
}

func NewLogger() Logger {
	s := &logger{
		day:    -1,
		utcOk:  true,
		pid:    os.Getpid(),
		arg:    newUnsafeArg(),
		l:      &sync.RWMutex{},
		l_sync: &sync.Mutex{},
		flag:   cc.NewAtomFlag()}
	s.cond = sync.NewCond(s.l_sync)
	s.start()
	return s
}

// SetPrename
func (s *logger) SetPrename(name string) {
	s.prename = name
}

// GetPrename
func (s *logger) GetPrename() string {
	return s.prename
}

// TimezoneString
func (s *logger) TimezoneString() string {
	return s.arg.timezoneString()
}

// GetTimezone
func (s *logger) GetTimezone() Timezone {
	return s.arg.getTimezone()
}

// SetTimezone
func (s *logger) SetTimezone(timezone Timezone) {
	switch timezone == s.arg.getTimezone() {
	case false:
		switch s.arg.setTimezone(timezone) {
		case true:
			s.setting(true)
		default:
			s.setting(true)
		}
	}
}

// ModeString
func (s *logger) ModeString() string {
	return s.arg.modeString()
}

// GetMode
func (s *logger) GetMode() Mode {
	return s.arg.getMode()
}

// SetMode
func (s *logger) SetMode(mode Mode) {
	switch mode == s.arg.getMode() {
	case false:
		switch s.arg.setMode(mode) {
		case true:
			s.setting(false)
		default:
			s.setting(false)
		}
	}
}

// StyleString
func (s *logger) StyleString() string {
	return s.arg.styleString()
}

// GetStyle
func (s *logger) GetStyle() Style {
	return s.arg.getStyle()
}

// SetStyle
func (s *logger) SetStyle(style Style) {
	switch style == s.arg.getStyle() {
	case false:
		switch s.arg.setStyle(style) {
		case true:
			s.setting(false)
		default:
			s.setting(false)
		}
	}
}

// LevelString
func (s *logger) LevelString() string {
	return s.arg.levelString()
}

// GetLevel
func (s *logger) GetLevel() Level {
	return s.arg.getLevel()
}

// SetLevel
func (s *logger) SetLevel(level Level) {
	switch level == s.arg.getLevel() {
	case false:
		switch s.arg.setLevel(level) {
		case true:
			s.setting(false)
		default:
			s.setting(false)
		}
	}
}

// check
func (s *logger) check(level Level) bool {
	return level <= s.arg.getLevel()
}

// utc_Ok
func (s *logger) utc_Ok() (ok bool) {
	s.l.RLock()
	ok = s.utcOk
	s.l.RUnlock()
	return
}

// update
func (s *logger) update(tm *time.Time) (ok bool) {
	switch s.utc_Ok() {
	case true:
		t := time.Now()
		s.l.Lock()
		s.utcOk = convertUTC(&t, &s.tm, s.arg.getTimezone())
		switch s.utcOk {
		case true:
			*tm = s.tm
			s.l.Unlock()
			ok = true
		default:
			s.l.Unlock()
			goto ERR
		}
	default:
	}
	return
ERR:
	ErrorfTLF("error")
	return
}

// get
func (s *logger) get(tm *time.Time) {
	s.l.RLock()
	*tm = s.tm
	s.l.RUnlock()
}

// SetColor
func (s *logger) SetColor(level Level, prefix, context int) {

}

// setting
func (s *logger) setting(update bool) {
	switch update {
	case true:
		t := time.Now()
		var tm time.Time
		s.l.Lock()
		s.utcOk = convertUTC(&t, &tm, s.arg.getTimezone())
		switch s.utcOk {
		case true:
			s.l.Unlock()
			setting(&tm, s.arg.getTimezone())
		default:
			s.l.Unlock()
			setting(&tm, s.arg.getTimezone())
		}
	default:
		switch s.utc_Ok() {
		case true:
			t := time.Now()
			var tm time.Time
			s.l.Lock()
			s.utcOk = convertUTC(&t, &tm, s.arg.getTimezone())
			switch s.utcOk {
			case true:
				s.l.Unlock()
				setting(&tm, s.arg.getTimezone())
			default:
				s.l.Unlock()
				setting(&tm, s.arg.getTimezone())
			}
		default:
			setting(nil, s.arg.getTimezone())
		}
	}
}

// checkSync
// func (s *logger) checkSync(style Style) {
// 	if (style & F_SYNC) > 0 {
// 		s.notify()
// 	}
// }

// notify
func (s *logger) notify() {
	s.l_sync.Lock()
	s.sync = true
	s.cond.Signal()
	s.l_sync.Unlock()
}

// Wait
func (s *logger) Wait() {
	s.l_sync.Lock()
	for !s.sync {
		s.cond.Wait()
	}
	s.sync = false
	s.l_sync.Unlock()
}

// checkDir
func (s *logger) checkDir() bool {
	switch s.prefix {
	case "":
	default:
		switch s.mkdir {
		case false:
			dir := filepath.Dir(s.prefix)
			_, err := os.Stat(dir)
			if err != nil && os.IsNotExist(err) {
				err := os.MkdirAll(dir, 0777)
				if err != nil {
					panic(err.Error())
				}
				s.mkdir = true
			} else {
				s.mkdir = true
			}
		}
	}
	return s.mkdir
}

// Init
func (s *logger) Init(dir string, prename string, logsize int64) {
	if dir == "" {
		dir = "."
	}
	s.size = logsize
	if s.prename = prename; prename != "" {
		s.prefix = dir + "/" + prename + "."
	} else {
		s.prefix = dir + "/"
	}
}

// name
func (s *logger) name(space bool) string {
	switch s.prename {
	case "":
		return ""
	default:
		switch space {
		case true:
			return strings.Join([]string{" <", s.prename, "> "}, "")
		default:
			return strings.Join([]string{" <", s.prename, ">"}, "")
		}

	}
}

// format
func (s *logger) format(level Level, style Style, skip int) (prefix string) {
	var tm time.Time
	ok := s.update(&tm)
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
		b.WriteString(strconv.Itoa(s.pid))
		b.WriteString(s.name(false))
		b.WriteString(" ")
		switch ok {
		case true:
			b.WriteString(String(s.arg.getTimezone()))
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
		b.WriteString(strconv.Itoa(s.pid))
		b.WriteString(s.name(false))
		switch ok {
		case true:
			b.WriteString(" ")
			b.WriteString(String(s.arg.getTimezone()))
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
		b.WriteString(strconv.Itoa(s.pid))
		b.WriteString(s.name(false))
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
		b.WriteString(strconv.Itoa(s.pid))
		b.WriteString(s.name(false))
		switch ok {
		case true:
			b.WriteString(" ")
			b.WriteString(String(s.arg.getTimezone()))
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
		b.WriteString(strconv.Itoa(s.pid))
		b.WriteString(s.name(false))
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
		b.WriteString(strconv.Itoa(s.pid))
		b.WriteString(s.name(false))
		b.WriteString(" ")
		switch ok {
		case true:
			b.WriteString(String(s.arg.getTimezone()))
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
		b.WriteString(strconv.Itoa(s.pid))
		b.WriteString(s.name(false))
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
		b.WriteString(strconv.Itoa(s.pid))
		b.WriteString(s.name(false))
		b.WriteString(" ")
		switch ok {
		case true:
			b.WriteString(String(s.arg.getTimezone()))
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
		b.WriteString(strconv.Itoa(s.pid))
		b.WriteString(s.name(false))
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
		b.WriteString(s.name(true))
		prefix = b.String()
	}
	return
}

// open
func (s *logger) open(path string) {
	if path == "" {
		panic(errors.New("error"))
	}
	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(errors.New(err.Error()))
	} else {
		s.fd = fd
		s.bio = bufio.NewWriter(s.fd)
	}
}

// write
func (s *logger) writeMsg(msg *Msg, pos int, style Style) {
	str := strings.Join([]string{msg.first, msg.second}, "")
	switch bio {
	case 1:
		s.write_bio(str, len(msg.first), style)
	default:
		s.write(str, len(msg.first), style)
	}
}

// writeStack
func (s *logger) writeStack(stack string) {
	switch bio {
	case 1:
		s.write_bio_(stack)
	default:
		s.write_(stack)
	}
}

// write
func (s *logger) write(msg string, pos int, style Style) {
	switch style {
	case F_DETAIL, F_DETAIL_SYNC:
		s.write_(msg[1:])
	case F_TMSTMP, F_TMSTMP_SYNC:
		s.write_(msg[1:])
	case F_FN, F_FN_SYNC:
		s.write_(msg[1:])
	case F_TMSTMP_FN, F_TMSTMP_FN_SYNC:
		s.write_(msg[1:])
	case F_FL, F_FL_SYNC:
		s.write_(msg[1:])
	case F_TMSTMP_FL, F_TMSTMP_FL_SYNC:
		s.write_(msg[1:])
	case F_FL_FN, F_FL_FN_SYNC:
		s.write_(msg[1:])
	case F_TMSTMP_FL_FN, F_TMSTMP_FL_FN_SYNC:
		s.write_(msg[1:])
	case F_TEXT, F_TEXT_SYNC:
		s.write_(msg[1:])
	case F_PURE, F_PURE_SYNC:
		fallthrough
	default:
		s.write_(msg[2:])
	}
}

// write_bio
func (s *logger) write_bio(msg string, pos int, style Style) {
	switch style {
	case F_DETAIL, F_DETAIL_SYNC:
		s.write_bio_(msg[1:])
	case F_TMSTMP, F_TMSTMP_SYNC:
		s.write_bio_(msg[1:])
	case F_FN, F_FN_SYNC:
		s.write_bio_(msg[1:])
	case F_TMSTMP_FN, F_TMSTMP_FN_SYNC:
		s.write_bio_(msg[1:])
	case F_FL, F_FL_SYNC:
		s.write_bio_(msg[1:])
	case F_TMSTMP_FL, F_TMSTMP_FL_SYNC:
		s.write_bio_(msg[1:])
	case F_FL_FN, F_FL_FN_SYNC:
		s.write_bio_(msg[1:])
	case F_TMSTMP_FL_FN, F_TMSTMP_FL_FN_SYNC:
		s.write_bio_(msg[1:])
	case F_TEXT, F_TEXT_SYNC:
		s.write_bio_(msg[1:])
	case F_PURE, F_PURE_SYNC:
		fallthrough
	default:
		s.write_bio_(msg[2:])
	}
}

// write_bio_
func (s *logger) write_bio_(msg string) {
	if s.bio != nil {
		_, err := s.bio.WriteString(msg)
		if err != nil {
			panic(errors.New(err.Error()))
		}
		s.bio.Flush()
	}
}

// write_
func (s *logger) write_(msg string) {
	if s.fd != nil {
		_, err := s.fd.WriteString(msg)
		if err != nil {
			panic(errors.New(err.Error()))
		}
	}
}

// close
func (s *logger) close() {
	if s.fd != nil {
		s.fd.Close()
		s.fd = nil
		s.bio = nil
	}
}

// Sprint
func (s *logger) Sprint(level Level, style Style, skip int, format string, v ...any) (prefix, content string) {
	prefix = s.format(level, style, skip)
	content = fmt.Sprintf(format, v...)
	return
}

// Write
func (s *logger) Write(stack string, level Level, style Style, skip int, format string, v ...any) {
	if s.check(level) {
		prefix, content := s.Sprint(level, style, skip, format, v...)
		s.push(prefix, content+"\n", len(prefix), style, stack)
	}
}

// push
func (s *logger) push(prefix, content string, pos int, style Style, stack string) {
	s.start()
	s.pipe.Do(NewMessageT(NewMessage(NewMsg(prefix, content), stack), NewFlags(pos, style)))
}

// shift
func (s *logger) shift(tm *time.Time) {
	if tm.Day() != s.day {
		s.close()
		// 2006/01/02 15:04:05.000000
		YMD := tm.Format("2006-01-02")
		HMS := tm.Format("15.04.05")
		s.path = strings.Join([]string{
			s.prefix, strconv.Itoa(s.pid), "_", YMD, ".", HMS, ".log",
		}, "")
		_, err := os.Stat(s.path)
		if err != nil && os.IsNotExist(err) {
		} else {
			os.Remove(s.path)
		}
		s.open(s.path)
		s.day = tm.Day()
	} else {
		sta, err := os.Stat(s.path)
		if err != nil && os.IsNotExist(err) {
			s.close()
			s.open(s.path)
			return
		}
		if sta.Size() < s.size {
		} else {
			s.close()
			YMD := tm.Format("2006-01-02")
			HMS := tm.Format("15.04.05.000000")
			s.path = strings.Join([]string{
				s.prefix, strconv.Itoa(s.pid), "_", YMD, ".", HMS, ".log",
			}, "")
			_, err := os.Stat(s.path)
			if err != nil && os.IsNotExist(err) {
				s.open(s.path)
			} else {
				for i := 0; ; {
					s.path = strings.Join([]string{
						s.prefix, strconv.Itoa(s.pid), "_", YMD, ".", HMS, ".", strconv.Itoa(i), ".log",
					}, "")
					_, err := os.Stat(s.path)
					if err != nil && os.IsNotExist(err) {
						s.open(s.path)
						break
					} else {
						i++
					}
				}
			}
		}
	}
}

// func (s *logger) onTimer(timerID uint32, dt int32, args ...any) bool {
// 	if len(args) == 0 {
// 		panic(errors.New("logs.args 0"))
// 	}
// 	if args[0] == nil {
// 		panic(errors.New("logs.args[0] is nil"))
// 	}
// 	switch args[0].(type) {
// 	default:
// 		break
// 	}
// 	return true
// }

func getlevel(c byte) Level {
	switch c {
	case 'F':
		return LVL_FATAL
	case 'E':
		return LVL_ERROR
	case 'W':
		return LVL_WARN
	case 'I':
		return LVL_INFO
	case 'T':
		return LVL_TRACE
	case 'D':
		return LVL_DEBUG
	}
	panic(errors.New("error"))
}

func (s *logger) handler(msg any, args ...any) (exit bool) {
	switch msg := msg.(type) {
	case *MessageT:
		// messageT, _ := msg.(*MessageT)
		messageT := msg
		message := messageT.first
		flags := messageT.second
		pos := flags.first
		style := flags.second
		msgData := message.first
		stack := message.second
		prefix := msgData.first
		// content := msgData.second
		mode := s.arg.getMode()
		switch mode {
		case M_FILE_ONLY, M_STDOUT_FILE:
			switch prefix[0] {
			case TAG[0]:
				switch s.checkDir() {
				default:
					mode = M_STDOUT_ONLY
				case true:
					var tm time.Time
					s.get(&tm)
					s.shift(&tm)
				}
			case TAG[1]:
				mode = M_STDOUT_ONLY
			}
		}
		level := getlevel(conv.StrToByte(prefix)[1])
		switch level {
		case LVL_FATAL:
			switch mode {
			case M_FILE_ONLY, M_STDOUT_FILE:
				s.writeMsg(msgData, pos, style)
				s.writeStack(stack)
			}
			switch mode {
			case M_STDOUT_ONLY, M_STDOUT_FILE:
				s.stdoutbuf(msgData, pos, level, style, stack)
			}
		case LVL_ERROR, LVL_WARN, LVL_INFO, LVL_TRACE, LVL_DEBUG:
			switch mode {
			case M_FILE_ONLY, M_STDOUT_FILE:
				s.writeMsg(msgData, pos, style)
			}
			switch mode {
			case M_STDOUT_ONLY, M_STDOUT_FILE:
				s.stdoutbuf(msgData, pos, level, style, "")
			}
		}
		msgData.Put()
		flags.Put()
		message.Put()
		messageT.Put()
		exit = (msg.second.second & F_SYNC) > 0
	}
	return
}

func (s *logger) stdoutbuf(msg *Msg, pos int, level Level, style Style, stack string) {
	switch level {
	case LVL_FATAL:
		switch style {
		case F_DETAIL, F_DETAIL_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
			Print(color[level][0], stack)
		case F_TMSTMP, F_TMSTMP_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
			Print(color[level][0], stack)
		case F_FN, F_FN_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
			Print(color[level][0], stack)
		case F_TMSTMP_FN, F_TMSTMP_FN_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
			Print(color[level][0], stack)
		case F_FL, F_FL_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
			Print(color[level][0], stack)
		case F_TMSTMP_FL, F_TMSTMP_FL_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
			Print(color[level][0], stack)
		case F_FL_FN, F_FL_FN_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
			Print(color[level][0], stack)
		case F_TMSTMP_FL_FN, F_TMSTMP_FL_FN_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
			Print(color[level][0], stack)
		case F_TEXT, F_TEXT_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
			Print(color[level][0], stack)
		case F_PURE, F_PURE_SYNC:
			fallthrough
		default:
			Print(color[level][0], msg.second)
			Print(color[level][0], stack)
		}
	case LVL_ERROR, LVL_WARN, LVL_INFO, LVL_TRACE, LVL_DEBUG:
		switch style {
		case F_DETAIL, F_DETAIL_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
		case F_TMSTMP, F_TMSTMP_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
		case F_FN, F_FN_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
		case F_TMSTMP_FN, F_TMSTMP_FN_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
		case F_FL, F_FL_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
		case F_TMSTMP_FL, F_TMSTMP_FL_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
		case F_FL_FN, F_FL_FN_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
		case F_TMSTMP_FL_FN, F_TMSTMP_FL_FN_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
		case F_TEXT, F_TEXT_SYNC:
			Print(color[level][0], msg.first[1:])
			Print(color[level][1], msg.second)
		case F_PURE, F_PURE_SYNC:
			fallthrough
		default:
			Print(color[level][0], msg.second)
		}
	}
}

func (s *logger) onQuit(slot run.Slot) {
	s.close()
	s.notify()
	s.reset()
}

// Close
func (s *logger) Close() {
	s.pipe.Close()
}

// start
func (s *logger) start() {
	if s.pipe == nil && s.flag.TestSet() {
		mq := lq.NewQueue(1000)
		runner := NewProcessor(s.handler)
		s.pipe = pipe.NewPipeWithQuit(i32.New(), "logger.pipe", mq, runner, s.onQuit)
		s.flag.Reset()
	}
	s.wait_started()
}

// wait_started
func (s *logger) wait_started() {
	for {
		if s.pipe != nil {
			break
		}
	}
}

// reset
func (s *logger) reset() {
	s.pipe = nil
}
