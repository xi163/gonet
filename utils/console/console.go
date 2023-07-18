package console

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

var Clear = map[string]func(){
	"windows": func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	},
	"linux": func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	},
}

// 标准输入
func Read(callback func(string) int) {
	for {
		// 从标准输入读取字符串，以\n为分割
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}
		// 去掉读入内容的空白符
		text = strings.TrimSpace(text)
		rc := callback(text)
		if rc < 0 {
			break
		}
	}
}
