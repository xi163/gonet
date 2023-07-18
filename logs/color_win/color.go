package color_win

import (
	"fmt"
	"syscall"
)

var (
	FOREGROUND_BLUE      uint8 = 0x01
	FOREGROUND_GREEN     uint8 = 0x02
	FOREGROUND_RED       uint8 = 0x04
	FOREGROUND_INTENSITY uint8 = 0x08
	FOREGROUND_HIGHLIGHT uint8 = 0x0F
	FOREGROUND_NORMAL    uint8 = 0x07
	FOREGROUND_BLACK     uint8 = 0x00
	BACKGROUND_BLUE      uint8 = 0x10
	BACKGROUND_GREEN     uint8 = 0x20
	BACKGROUND_RED       uint8 = 0x40
	BACKGROUND_INTENSITY uint8 = 0x80
	BACKGROUND_HIGHLIGHT uint8 = 0xF0
	BACKGROUND_NORMAL    uint8 = 0x70
	BACKGROUND_BLACK     uint8 = 0x00

	FOREGROUND_Red         int = int(FOREGROUND_RED)                                                             //红
	FOREGROUND_Green       int = int(FOREGROUND_GREEN)                                                           //绿
	FOREGROUND_Blue        int = int(FOREGROUND_BLUE)                                                            //蓝
	FOREGROUND_Yellow      int = int(FOREGROUND_RED | FOREGROUND_GREEN)                                          //黄
	FOREGROUND_Cyan        int = int(FOREGROUND_GREEN | FOREGROUND_BLUE)                                         //青
	FOREGROUND_Purple      int = int(FOREGROUND_RED | FOREGROUND_BLUE)                                           //紫
	FOREGROUND_White       int = int(FOREGROUND_RED | FOREGROUND_GREEN | FOREGROUND_BLUE)                        //白
	FOREGROUND_Gray        int = int(FOREGROUND_INTENSITY)                                                       //灰
	FOREGROUND_Black       int = int(FOREGROUND_BLACK)                                                           //黑
	FOREGROUND_LightRed    int = int(FOREGROUND_INTENSITY | FOREGROUND_RED)                                      //淡红
	FOREGROUND_HighGreen   int = int(FOREGROUND_INTENSITY | FOREGROUND_GREEN)                                    //亮绿
	FOREGROUND_LightBlue   int = int(FOREGROUND_INTENSITY | FOREGROUND_BLUE)                                     //淡蓝
	FOREGROUND_LightYellow int = int(FOREGROUND_INTENSITY | FOREGROUND_RED | FOREGROUND_GREEN)                   //淡黄
	FOREGROUND_HighCyan    int = int(FOREGROUND_INTENSITY | FOREGROUND_GREEN | FOREGROUND_BLUE)                  //亮青
	FOREGROUND_Pink        int = int(FOREGROUND_INTENSITY | FOREGROUND_RED | FOREGROUND_BLUE)                    //粉红
	FOREGROUND_HighWhite   int = int(FOREGROUND_INTENSITY | FOREGROUND_RED | FOREGROUND_GREEN | FOREGROUND_BLUE) //亮白

	BACKGROUND_Red         int = int(BACKGROUND_RED)
	BACKGROUND_Green       int = int(BACKGROUND_GREEN)
	BACKGROUND_Blue        int = int(BACKGROUND_BLUE)
	BACKGROUND_Yellow      int = int(BACKGROUND_RED | BACKGROUND_GREEN)
	BACKGROUND_Cyan        int = int(BACKGROUND_GREEN | BACKGROUND_BLUE)
	BACKGROUND_Purple      int = int(BACKGROUND_RED | BACKGROUND_BLUE)
	BACKGROUND_White       int = int(BACKGROUND_RED | BACKGROUND_GREEN | BACKGROUND_BLUE)
	BACKGROUND_Gray        int = int(BACKGROUND_INTENSITY)
	BACKGROUND_Black       int = int(BACKGROUND_BLACK)
	BACKGROUND_LightRed    int = int(BACKGROUND_INTENSITY | BACKGROUND_RED)
	BACKGROUND_HighGreen   int = int(BACKGROUND_INTENSITY | BACKGROUND_GREEN)
	BACKGROUND_LightBlue   int = int(BACKGROUND_INTENSITY | BACKGROUND_BLUE)
	BACKGROUND_LightYellow int = int(BACKGROUND_INTENSITY | BACKGROUND_RED | BACKGROUND_GREEN)
	BACKGROUND_HighCyan    int = int(BACKGROUND_INTENSITY | BACKGROUND_GREEN | BACKGROUND_BLUE)
	BACKGROUND_Pink        int = int(BACKGROUND_INTENSITY | BACKGROUND_RED | BACKGROUND_BLUE)
	BACKGROUND_HighWhite   int = int(BACKGROUND_INTENSITY | BACKGROUND_RED | BACKGROUND_GREEN | BACKGROUND_BLUE)
)

// SET GOOS=windows
var (
	kernel32                *syscall.LazyDLL  = syscall.NewLazyDLL(`kernel32.dll`)
	GetStdHandle            *syscall.LazyProc = kernel32.NewProc(`GetStdHandle`)
	SetConsoleTextAttribute *syscall.LazyProc = kernel32.NewProc(`SetConsoleTextAttribute`)
	CloseHandle             *syscall.LazyProc = kernel32.NewProc(`CloseHandle`)
)

func PrintFormat(color int, format string, v ...any) {
	handle, _, _ := SetConsoleTextAttribute.Call(uintptr(syscall.Stdout), uintptr(color))
	fmt.Printf(format, v...)
	CloseHandle.Call(handle)
}

func Print(color int, msg string) {
	handle, _, _ := SetConsoleTextAttribute.Call(uintptr(syscall.Stdout), uintptr(color))
	fmt.Printf("%v", msg)
	CloseHandle.Call(handle)
}
