package env

import "runtime"

func G() (path, cmd, ext string) {
	switch runtime.GOOS {
	case "linux":
		path = "/"
		cmd = "./"
		ext = ""
	case "windows":
		path += "\\"
		cmd = ""
		ext = ".exe"
	}
	return
}
