package Fn

import (
	"path"
	"strings"
)

func Split(name string) (string, string) {
	//path/pkg.(type).func
	//path/pkg.(type[...]).func
	name = strings.ReplaceAll(name, "...", "T")
	_, f := path.Split(name)
	v := strings.Split(f, ".")
	if len(v) >= 3 {
		//pkg.(type).func
		v[1] = strings.Replace(v[1], "(", "", 1)
		if []byte(v[1])[0] == '*' {
			v[1] = strings.Replace(v[1], "*", "", 1)
		}
		v[1] = strings.Replace(v[1], ")", "", 1)
		//pkg type.func
		return v[0], v[1] + "." + v[2]
	}
	//pkg func
	return v[0], v[1]
}
