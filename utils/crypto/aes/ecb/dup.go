package ecb

// Utility routines

func dup(p []byte) []byte {
	q := make([]byte, len(p))
	copy(q, p)
	return q
}
