package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils/conv"
	"github.com/cwloo/gonet/utils/env"
)

// <summary>
// ARG
// <summary>
type ARG struct {
	conf struct {
		dir  string
		name string
	}
	pattern struct {
		ID   []string
		LVL  []string
		CONF []string
		LOG  []string
		NAME []string
		KEYS map[string][]string
	}
	dict map[string]string
}

func newArg() *ARG {
	return &ARG{
		pattern: struct {
			ID   []string
			LVL  []string
			CONF []string
			LOG  []string
			NAME []string
			KEYS map[string][]string
		}{
			ID:   []string{"id", "i"},
			CONF: []string{"config", "conf", "cfg", "c"},
			LOG:  []string{"logdir", "log_dir", "log-dir", "log", "l"},
			LVL:  []string{"dirlevel", "dir_level", "dir-level", "dirlvl", "dir_lvl", "dir-lvl"},
			NAME: []string{"configname", "config_name", "config-name", "confname", "conf_name", "conf-name", "cfgname", "cfg_name", "cfg-name"},
			KEYS: map[string][]string{},
		},
		dict: map[string]string{},
	}
}

func (s *ARG) SetConf(name string) {
	s.conf.name = env.CorrectPath(name)
}

func (s *ARG) AppendPattern(key string, v ...string) {
	s.pattern.KEYS[key] = append(s.pattern.KEYS[key], v...)
}

// id
func (s *ARG) id() (id int) {
	{
	ID:
		for _, c := range s.pattern.ID {
			v, ok := s.dict[c]
			switch ok {
			case true:
				id = conv.StrToInt(v)
				break ID
			}
		}
	}
	return
}

func (s *ARG) formatId(id int) string {
	s.assertId()
	return strings.Join([]string{"--", s.pattern.ID[0], "=", strconv.Itoa(id)}, "")
}

// dir
func (s *ARG) root() (dir string) {
	lvl := 0
LVL:
	for _, c := range s.pattern.LVL {
		v, ok := s.dict[c]
		switch ok {
		case true:
			lvl = conv.StrToInt(v)
			break LVL
		}
	}
	dir = env.Dir
	for l := 0; l < lvl; l++ {
		dir = filepath.Dir(dir)
	}
	return
}

func (s *ARG) formatLvl(lvl int) string {
	s.assertLvl()
	return strings.Join([]string{"--", s.pattern.LVL[0], "=", strconv.Itoa(lvl)}, "")
}

// conf
func (s *ARG) Conf(root string) (dir string) {
NAME:
	for _, c := range s.pattern.NAME {
		v, ok := s.dict[c]
		switch ok {
		case true:
			s.conf.name = env.CorrectPath(v)
			break NAME
		}
	}
CONF:
	for _, c := range s.pattern.CONF {
		v, ok := s.dict[c]
		switch ok {
		case true:
			dir = v
			break CONF
		}
	}
	switch dir {
	case "":
		switch len(s.conf.name) > 0 && s.conf.name[0:1] == env.P {
		case true:
			s.conf.name = strings.Replace(s.conf.name, env.P, "", 1)
		}
		switch s.conf.dir {
		case "":
			switch len(root) > 0 && root[len(root)-1:] == env.P {
			case true:
				dir = strings.Join([]string{root, s.conf.name}, "")
			default:
				dir = strings.Join([]string{root, env.P, s.conf.name}, "")
			}
		default:
			switch len(s.conf.dir) > 0 && s.conf.dir[0:1] == env.P {
			case true:
				s.conf.dir = strings.Replace(s.conf.dir, env.P, "", 1)
			}
			switch len(root) > 0 && root[len(root)-1:] == env.P {
			case true:
				switch len(s.conf.dir) > 0 && root[len(s.conf.dir)-1:] == env.P {
				case true:
					dir = strings.Join([]string{root, s.conf.dir, s.conf.name}, "")
				default:
					dir = strings.Join([]string{root, s.conf.dir, env.P, s.conf.name}, "")
				}
			default:
				switch len(s.conf.dir) > 0 && root[len(s.conf.dir)-1:] == env.P {
				case true:
					dir = strings.Join([]string{root, env.P, s.conf.dir, s.conf.name}, "")
				default:
					dir = strings.Join([]string{root, env.P, s.conf.dir, env.P, s.conf.name}, "")
				}
			}
		}
	default:
		dir = env.CorrectPath(dir)
	}
	return
}

func (s *ARG) formatConf(dir string) string {
	s.assertConf()
	return strings.Join([]string{"--", s.pattern.CONF[0], "=", dir}, "")
}

// log
func (s *ARG) log() (dir string) {
	{
	LOG:
		for _, c := range s.pattern.LOG {
			v, ok := s.dict[c]
			switch ok {
			case true:
				dir = v
				break LOG
			}
		}
		switch dir {
		case "":
		default:
			dir = env.CorrectPath(dir)
		}
	}
	return
}

func (s *ARG) formatLog(dir string) string {
	s.assertLog()
	return strings.Join([]string{"--", s.pattern.LOG[0], "=", dir}, "")
}

// arg
func (s *ARG) arg(key string) string {
	return s.dict[key]
}

func (s *ARG) formatArg(key, val string) (c string) {
	c = strings.Join([]string{"--", key, "=", val}, "")
	return
}

// patternArg
func (s *ARG) patternArg(key string) (val string) {
	keys, ok := s.pattern.KEYS[key]
	switch ok {
	case true:
	KEY:
		for _, c := range keys {
			v, ok := s.dict[c]
			switch ok {
			case true:
				val = v
				break KEY
			}
		}
	}
	return
}

func (s *ARG) formatPatternArg(key, val string) (c string) {
	keys, ok := s.pattern.KEYS[key]
	switch ok {
	case true:
		switch len(keys) > 0 {
		case true:
			c = strings.Join([]string{"--", keys[0], "=", val}, "")
		default:
			goto ERR
		}
	default:
		goto ERR
	}
	return
ERR:
	logs.Fatalf("error")
	return
}

// parse
func (s *ARG) parse() (id int, dir, conf, log string) {
	for _, v := range os.Args {
		m := strings.Split(v, "=")
		switch len(m) == 2 {
		case true:
			m[0] = env.CorrectArg(m[0])
			m[0] = strings.ToLower(m[0])
			s.dict[m[0]] = m[1]
		}
	}
	logs.Warnf("dir=%v conf=%v %v", s.root(), s.Conf(s.root()), os.Args)
	return
}

func (s *ARG) assertId() {
	switch len(s.pattern.ID) == 0 {
	case true:
		logs.Fatalf("error")
	}
}

func (s *ARG) assertLvl() {
	switch len(s.pattern.LVL) == 0 {
	case true:
		logs.Fatalf("error")
	}
}

func (s *ARG) assertConf() {
	switch len(s.pattern.CONF) == 0 {
	case true:
		logs.Fatalf("error")
	}
}

func (s *ARG) assertLog() {
	switch len(s.pattern.LOG) == 0 {
	case true:
		logs.Fatalf("error")
	}
}
