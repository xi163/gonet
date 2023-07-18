package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	src = rand.NewSource(time.Now().UnixNano())
	arr = []string{
		"0123456789",
		"abcdefghijklmnopqrstuvwxyz",
		"abcdefghijklmnopqrstuvwxyz0123456789",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		"AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz",
		"AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789",
	}
)

func NumberStr(n int) string {
	result := make([]byte, n)
	r := rand.New(src)
	for i := 0; i < n; i++ {
		x := r.Intn(len(arr[0]))
		if i == 0 && arr[0][x] == '0' {
			i--
		} else {
			result = append(result, arr[0][x])
		}
	}
	return string(result)
}

func LowerCharStr(n int) string {
	result := make([]byte, n)
	r := rand.New(src)
	for i := 0; i < n; i++ {
		result = append(result, arr[1][r.Intn(len(arr[1]))])
	}
	return string(result)
}

func LowerCharNumStr(n int) string {
	result := make([]byte, n)
	r := rand.New(src)
	for i := 0; i < n; i++ {
		result = append(result, arr[2][r.Intn(len(arr[2]))])
	}
	return string(result)
}

func UpperCharStr(n int) string {
	result := make([]byte, n)
	r := rand.New(src)
	for i := 0; i < n; i++ {
		result = append(result, arr[3][r.Intn(len(arr[3]))])
	}
	return string(result)
}

func UpperCharNumStr(n int) string {
	result := make([]byte, n)
	r := rand.New(src)
	for i := 0; i < n; i++ {
		result = append(result, arr[4][r.Intn(len(arr[4]))])
	}
	return string(result)
}

func CharStr(n int) string {
	result := make([]byte, n)
	r := rand.New(src)
	for i := 0; i < n; i++ {
		result = append(result, arr[5][r.Intn(len(arr[5]))])
	}
	return string(result)
}

func CharNumStr(n int) string {
	result := make([]byte, n)
	r := rand.New(src)
	for i := 0; i < n; i++ {
		result = append(result, arr[6][r.Intn(len(arr[6]))])
	}
	return string(result)
}
