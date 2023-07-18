package logs_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/logs/color_linux"
	//"github.com/cwloo/gonet/logs/color_win"
)

func TestMain(m *testing.M) {
	m.Run()
}

func Test(t *testing.T) {
	logs.SetPrename("gonet")
	logs.SetTimezone(logs.MY_CST)
	logs.SetMode(logs.M_STDOUT_FILE)
	logs.SetStyle(logs.F_DETAIL)
	logs.SetLevel(logs.LVL_DEBUG)
	logs.Init("/home", "gonet", 100000000)
	t.Run("logs_test:", out_test)
	t.Run("logs_test:", path_test)
}

func out_test(t *testing.T) {
	color_linux.EnumColorStyle()

	logs.InfofP("errMsg=%v", logs.SprintErrorf("error"))

	logs.InfofP("------------------------ F_DETAIL")
	logs.SetStyle(logs.F_DETAIL)
	// logs.Fatalf("hello,word")
	logs.Errorf("hello,word")
	logs.Warnf("hello,word")
	logs.Infof("hello,word")
	logs.Tracef("hello,word")
	logs.Debugf("hello,word")

	logs.InfofP("------------------------ F_TMSTMP")
	logs.SetStyle(logs.F_TMSTMP)
	// logs.Fatalf("hello,word")
	logs.Errorf("hello,word")
	logs.Warnf("hello,word")
	logs.Infof("hello,word")
	logs.Tracef("hello,word")
	logs.Debugf("hello,word")

	logs.InfofP("------------------------ F_FN")
	logs.SetStyle(logs.F_FN)
	// logs.Fatalf("hello,word")
	logs.Errorf("hello,word")
	logs.Warnf("hello,word")
	logs.Infof("hello,word")
	logs.Tracef("hello,word")
	logs.Debugf("hello,word")

	logs.InfofP("------------------------ F_TMSTMP_FN")
	logs.SetStyle(logs.F_TMSTMP_FN)
	// logs.Fatalf("hello,word")
	logs.Errorf("hello,word")
	logs.Warnf("hello,word")
	logs.Infof("hello,word")
	logs.Tracef("hello,word")
	logs.Debugf("hello,word")

	logs.InfofP("------------------------ F_FL")
	logs.SetStyle(logs.F_FL)
	// logs.Fatalf("hello,word")
	logs.Errorf("hello,word")
	logs.Warnf("hello,word")
	logs.Infof("hello,word")
	logs.Tracef("hello,word")
	logs.Debugf("hello,word")

	logs.InfofP("------------------------ F_TMSTMP_FL")
	logs.SetStyle(logs.F_TMSTMP_FL)
	// logs.Fatalf("hello,word")
	logs.Errorf("hello,word")
	logs.Warnf("hello,word")
	logs.Infof("hello,word")
	logs.Tracef("hello,word")
	logs.Debugf("hello,word")

	logs.InfofP("------------------------ F_FL_FN")
	logs.SetStyle(logs.F_FL_FN)
	// logs.Fatalf("hello,word")
	logs.Errorf("hello,word")
	logs.Warnf("hello,word")
	logs.Infof("hello,word")
	logs.Tracef("hello,word")
	logs.Debugf("hello,word")

	logs.InfofP("------------------------ F_TMSTMP_FL_FN")
	logs.SetStyle(logs.F_TMSTMP_FL_FN)
	// logs.Fatalf("hello,word")
	logs.Errorf("hello,word")
	logs.Warnf("hello,word")
	logs.Infof("hello,word")
	logs.Tracef("hello,word")
	logs.Debugf("hello,word")

	logs.InfofP("------------------------ F_TEXT")
	logs.SetStyle(logs.F_TEXT)
	// logs.Fatalf("hello,word")
	logs.Errorf("hello,word")
	logs.Warnf("hello,word")
	logs.Infof("hello,word")
	logs.Tracef("hello,word")
	logs.Debugf("hello,word")

	logs.InfofP("------------------------ F_PURE")
	logs.SetStyle(logs.F_PURE)
	// logs.Fatalf("hello,word")
	logs.Errorf("hello,word")
	logs.Warnf("hello,word")
	logs.Infof("hello,word")
	logs.Tracef("hello,word")
	logs.Debugf("hello,word")

	logs.Close()
}

func path_test(t *testing.T) {
	_, dir, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(dir), "../../..")
	logs.Errorf("\n%v\n%v\n", dir, path)

	dir1, _ := os.Getwd()
	path1 := filepath.Join(dir1, "../config")
	logs.Errorf("\n%v\n%v\n", dir1, path1)

	dir2, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path2 := filepath.Join(dir2, "../config")
	logs.Errorf("\n%v\n%v\n", dir2, path2)

	path3, _ := os.Executable()
	_, exec := filepath.Split(path)
	logs.Errorf("\n%v\n%v\n", path3, exec)
}
