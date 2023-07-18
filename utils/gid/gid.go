package gid

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	// "github.com/cwloo/gonet/logs"
)

func Getgid() int {
	// Recover()
	// b := make([]byte, 64)
	// b = b[:runtime.Stack(b, false)]
	// b = bytes.TrimPrefix(b, []byte("goroutine "))
	// b = b[:bytes.IndexByte(b, ' ')]
	// n, _ := strconv.ParseUint(string(b), 10, 64)
	// return uint32(n)
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	str := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	ID, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println(err)
		// logs.Errorf(err.Error())
	}
	return ID
}
