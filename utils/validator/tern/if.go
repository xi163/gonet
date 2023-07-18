package tern

func IF(yes bool, a any, b any) any {
	if yes {
		return a
	} else {
		return b
	}
}
