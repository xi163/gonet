package safe

func Call(f func()) (err error) {
	defer Catch()
	f()
	return
}

func Call2(f func() (any, error)) (err error) {
	defer Catch()
	f()
	return
}

func Call3(f func(args ...any), args ...any) (err error) {
	defer Catch()
	f(args)
	return
}

func Call4(f func(bool), b bool) (err error) {
	defer Catch()
	f(b)
	return
}
