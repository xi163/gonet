package color_win

import (
	"fmt"
	"syscall"
)

const (
	Red int = iota + 0
	Green
	Blue
	Yellow
	Cyan
	Purple
	White
	Gray
	Black
	LightRed
	HighGreen
	LightBlue
	LightYellow
	HighCyan
	Pink
	HighWhite
)

var (
	Fore = []struct {
		color int
		desc  string
	}{
		{FOREGROUND_Red, "Red"},
		{FOREGROUND_Green, "Green"},
		{FOREGROUND_Blue, "Blue"},
		{FOREGROUND_Yellow, "Yellow"},
		{FOREGROUND_Cyan, "Cyan"},
		{FOREGROUND_Purple, "Cyan"},
		{FOREGROUND_White, "White"},
		{FOREGROUND_Gray, "Gray"},
		{FOREGROUND_Black, "Black"},
		{FOREGROUND_LightRed, "LightRed"},
		{FOREGROUND_HighGreen, "HighGreen"},
		{FOREGROUND_LightBlue, "LightBlue"},
		{FOREGROUND_LightYellow, "LightYellow"},
		{FOREGROUND_HighCyan, "HighCyan"},
		{FOREGROUND_Pink, "Pink"},
		{FOREGROUND_HighWhite, "HighWhite"},
	}
	Back = []struct {
		color int
		desc  string
	}{
		{BACKGROUND_Red, "Red"},
		{BACKGROUND_Green, "Green"},
		{BACKGROUND_Blue, "Blue"},
		{BACKGROUND_Yellow, "Yellow"},
		{BACKGROUND_Cyan, "Cyan"},
		{BACKGROUND_Purple, "Purple"},
		{BACKGROUND_White, "White"},
		{BACKGROUND_Gray, "Gray"},
		{BACKGROUND_Black, "Black"},
		{BACKGROUND_LightRed, "LightRed"},
		{BACKGROUND_HighGreen, "HighGreen"},
		{BACKGROUND_LightBlue, "LightBlue"},
		{BACKGROUND_LightYellow, "LightYellow"},
		{BACKGROUND_HighCyan, "HighCyan"},
		{BACKGROUND_Pink, "Pink"},
		{BACKGROUND_HighWhite, "HighWhite"},
	}
)

func EnumColorStyle() {
	fmt.Printf("BACKGROUND\n")
	for b := Red; b <= HighWhite; b++ {
		handle, _, _ := SetConsoleTextAttribute.Call(uintptr(syscall.Stdout), uintptr(Back[b].color))
		fmt.Printf("  %s", Back[b].desc)
		CloseHandle.Call(handle)
	}
	fmt.Printf("FOREGROUND\n")
	for f := Red; f <= HighWhite; f++ {
		handle, _, _ := SetConsoleTextAttribute.Call(uintptr(syscall.Stdout), uintptr(Fore[f].color))
		fmt.Printf("  %s", Fore[f].desc)
		CloseHandle.Call(handle)
	}
	fmt.Printf("BACKGROUND.FOREGROUND\n")
	for b := Red; b <= HighWhite; b++ {
		for f := Red; f <= HighWhite; f++ {
			handle, _, _ := SetConsoleTextAttribute.Call(uintptr(syscall.Stdout), uintptr(Back[b].color|Fore[f].color))
			fmt.Printf("  %s.%s", Back[b].desc, Fore[f].desc)
			CloseHandle.Call(handle)
		}
		fmt.Println("")
	}
	fmt.Println("")
}
