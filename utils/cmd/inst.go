package cmd

var (
	cmd = newCmd()
)

func InitArgs(cb func(*ARG)) {
	cmd.InitArgs(cb)
}

func SetConf(name string) {
	cmd.SetConf(name)
}

func AppendPattern(key string, v ...string) {
	cmd.AppendPattern(key, v...)
}

func ParseArgs() {
	cmd.ParseArgs()
}

func Id() int {
	return cmd.Id()
}

func FormatId(id int) string {
	return cmd.FormatId(id)
}

func Root() string {
	return cmd.Root()
}

func FormatLvl(lvl int) string {
	return cmd.FormatLvl(lvl)
}

func Conf() string {
	return cmd.Conf()
}

func FormatConf(conf string) string {
	return cmd.FormatConf(conf)
}

func Log() string {
	return cmd.Log()
}

func FormatLog(log string) string {
	return cmd.FormatLog(log)
}

func Arg(key string) string {
	return cmd.Arg(key)
}

func FormatArg(key, val string) string {
	return cmd.FormatArg(key, val)
}

func PatternArg(key string) string {
	return cmd.PatternArg(key)
}

func FormatPatternArg(key, val string) string {
	return cmd.FormatPatternArg(key, val)
}
