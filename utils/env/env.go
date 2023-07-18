package env

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	Path, _     = os.Executable()
	Dir, Exe    = filepath.Split(Path)
	P, Cmd, Ext = G()
)

func CorrectArg(old string) (new string) {
	new = old
	exist := true
LOOP:
	for {
		switch len(new) >= 2 && new[0:2] == "--" {
		case true:
			new = strings.Replace(new, "--", "", 1)
		}
		switch len(new) > 0 && new[0:1] == "-" {
		case true:
			new = strings.Replace(new, "-", "", 1)
		default:
			exist = false
		}
		switch exist {
		case false:
			break LOOP
		}
	}
	return
}

func CorrectPath(old string) (new string) {
	new = old
	dp := strings.Join([]string{P, P}, "")
LOOP:
	for {
		switch strings.Contains(new, dp) {
		case true:
			new = strings.ReplaceAll(new, dp, P)
		default:
			break LOOP
		}
	}
	return
}
