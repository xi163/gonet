package tern

func IF[T any](cmp bool, a T, b T) T {
	switch cmp {
	case true:
		return a
	default:
		return b
	}
}
