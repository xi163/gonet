package color_linux

import "fmt"

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
	Color = []struct {
		color string
		desc  string
	}{
		{COLOR_Red, "Red"},
		{COLOR_Green, "Green"},
		{COLOR_Blue, "Blue"},
		{COLOR_Yellow, "Yellow"},
		{COLOR_Cyan, "Cyan"},
		{COLOR_Purple, "Cyan"},
		{COLOR_White, "White"},
		{COLOR_Gray, "Gray"},
		// {COLOR_Black, "Black"},
		// {COLOR_LightRed, "LightRed"},
		// {COLOR_HighGreen, "HighGreen"},
		// {COLOR_LightBlue, "LightBlue"},
		// {COLOR_LightYellow, "LightYellow"},
		// {COLOR_HighCyan, "HighCyan"},
		// {COLOR_Pink, "Pink"},
		// {COLOR_HighWhite, "HighWhite"},
	}
)

// \033     起始标记
// [d;b;fm  d-(0-终端默认设置 1-高亮显示 4-使用下划线 5-闪烁 7-反白显示 8-不可见) b-背景色 f-前景色/字体色
// \033[0m  结束标记 恢复终端默认样式
func EnumColorStyle() {
	for _, d := range []int{0, 1, 4, 5, 7, 8} { // 0,1,4,5,7,8
		for b := 40; b <= 47; b++ { // 背景色 = 40-47
			for f := 30; f <= 37; f++ { // 前景色 = 30-37
				// fmt.Printf(" %c[%d;%d;%dm%s(f=%d,b=%d,d=%d)%c[0m ", 0x1B, d, b, f, "", f, b, d, 0x1B)
				fmt.Printf("  \033[%d;%d;%dm%s\033[0m", d, b, f, fmt.Sprintf("\\033[%d;%d;%dmFont\\033[0m", d, b, f))
			}
			fmt.Println("")
		}
	}
	fmt.Println("")
}
