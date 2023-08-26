package run

import (
	"fmt"
	"os"

	"github.com/cwloo/gonet/utils/macro"
)

// 捕获panic内容并恢复程序运行，在panic之后触发，所以必须defer方式调用
func Catch() {
	if err := recover(); err != nil {
		fmt.Fprint(os.Stderr, "run.panic: ", macro.SprintErrorf(6, "%v", err))
	}
}
