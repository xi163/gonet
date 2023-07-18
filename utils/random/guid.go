package random

import "strings"

func CreateGUID() string {
	//7787711c-f168-4867-b4c0-82f1f566432e
	s0 := LowerCharNumStr(8)
	s1 := LowerCharNumStr(4)
	s2 := LowerCharNumStr(4)
	s3 := LowerCharNumStr(4)
	s4 := LowerCharNumStr(12)
	return strings.Join([]string{s0, s1, s2, s3, s4}, "-")
}
