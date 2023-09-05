package cmd

// Cmd
type Cmd interface {
	InitArgs(cb func(*ARG))
	SetConf(name string)
	AppendPattern(key string, v ...string)
	ParseArgs()

	Id() int
	FormatId(id int) string

	Root() string
	FormatLvl(lvl int) string

	Conf() string
	FormatConf(conf string) string

	Log() string
	FormatLog(log string) string

	Arg(key string) string
	FormatArg(key, val string) string

	PatternArg(key string) string
	FormatPatternArg(key, val string) string
}

type CMD struct {
	arg *ARG
}

func newCmd() Cmd {
	s := &CMD{arg: newArg()}
	return s
}

func (s *CMD) Inst() *ARG {
	return s.arg
}

func (s *CMD) InitArgs(cb func(*ARG)) {
	cb(s.arg)
}

func (s *CMD) SetConf(name string) {
	s.arg.SetConf(name)
}

func (s *CMD) AppendPattern(key string, v ...string) {
	s.arg.AppendPattern(key, v...)
}

func (s *CMD) ParseArgs() {
	s.arg.parse()
}

func (s *CMD) Id() int {
	return s.arg.id()
}

func (s *CMD) FormatId(id int) string {
	return s.arg.formatId(id)
}

func (s *CMD) Root() string {
	return s.arg.root()
}

func (s *CMD) FormatLvl(lvl int) string {
	return s.arg.formatLvl(lvl)
}

func (s *CMD) Conf() string {
	return s.arg.Conf(s.arg.root())
}

func (s *CMD) FormatConf(conf string) string {
	return s.arg.formatConf(conf)
}

func (s *CMD) Log() string {
	return s.arg.log()
}

func (s *CMD) FormatLog(log string) string {
	return s.arg.formatLog(log)
}

func (s *CMD) Arg(key string) string {
	return s.arg.arg(key)
}

func (s *CMD) FormatArg(key, val string) string {
	return s.arg.formatArg(key, val)
}

func (s *CMD) PatternArg(key string) string {
	return s.arg.patternArg(key)
}

func (s *CMD) FormatPatternArg(key, val string) string {
	return s.arg.formatPatternArg(key, val)
}
