package cb

// 回调处理单元
type Proc interface {

	// s.Exec(func(v any) {
	// }, []any{a, b, c})
	Exec(f Functor)

	// ExecTimeout(d time.Duration, f Functor, cb Functor)

	// s.Append(func(v any) {
	// }, []any{a, b, c})
	Append(f Functor)

	// AppendTimeout(d time.Duration, f Functor, cb Functor)
}
