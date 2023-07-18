package color_linux

import (
	"fmt"
	"strings"
)

const (
	Reset = "\033[0m"
)

var (
	COLOR_Red    = "\033[1;31m"
	COLOR_Green  = "\033[1;32m"
	COLOR_Yellow = "\033[1;33m"
	COLOR_Blue   = "\033[1;34m"
	COLOR_Purple = "\033[1;35m"
	COLOR_Cyan   = "\033[1;36m"
	COLOR_Gray   = "\033[1;37m"
	COLOR_White  = "\033[1;97m"
)

const (
	FOREGROUND_Red    = 0
	FOREGROUND_Green  = 1
	FOREGROUND_Yellow = 2
	FOREGROUND_Blue   = 3
	FOREGROUND_Purple = 4
	FOREGROUND_Cyan   = 5
	FOREGROUND_Gray   = 6
	FOREGROUND_White  = 7
)

var (
	COLOR = []string{
		COLOR_Red,
		COLOR_Green,
		COLOR_Yellow,
		COLOR_Blue,
		COLOR_Purple,
		COLOR_Cyan,
		COLOR_Gray,
		COLOR_White,
	}
)

func PrintFormat(color int, format string, v ...any) {
	fmt.Printf(strings.Join([]string{COLOR[color], format, Reset}, ""), v...)
}

func Print(color int, msg string) {
	fmt.Printf(strings.Join([]string{COLOR[color], "%v", Reset}, ""), msg)
}
