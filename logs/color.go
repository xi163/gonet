package logs

import (
	"github.com/cwloo/gonet/logs/color_linux"
)

var (
	color = [][2]int{
		{color_linux.FOREGROUND_Red, color_linux.FOREGROUND_Cyan},     //LVL_FATAL
		{color_linux.FOREGROUND_Red, color_linux.FOREGROUND_Cyan},     //LVL_ERROR
		{color_linux.FOREGROUND_Cyan, color_linux.FOREGROUND_Purple},  //LVL_WARN
		{color_linux.FOREGROUND_Purple, color_linux.FOREGROUND_White}, //LVL_INFO
		{color_linux.FOREGROUND_Yellow, color_linux.FOREGROUND_Green}, //LVL_TRACE
		{color_linux.FOREGROUND_Green, color_linux.FOREGROUND_Yellow}, //LVL_DEBUG
	}
)

func PrintFormat(color int, format string, v ...any) {
	color_linux.PrintFormat(color, format, v...)
}

func Print(color int, msg string) {
	color_linux.Print(color, msg)
}
