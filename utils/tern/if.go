package tern

func IF[T any](cmp bool, a T, b T) T {
	if cmp {
		return a
	} else {
		return b
	}
}
