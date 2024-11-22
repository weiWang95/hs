package utils

import (
	"math/rand/v2"
)

func Shuffle[T any](a []T) []T {
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	return a
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.IntN(len(letters))]
	}
	return string(b)
}

func Random(probabilities []int) int {
	r := rand.IntN(100) + 1

	for i, v := range probabilities {
		if r < v {
			return i
		}
		r -= v
	}

	return len(probabilities) - 1
}

func RandIntWithOut(num int, out int) int {
	if num == 0 || (num == 1 && out == 0) {
		return -1
	}
	r := rand.IntN(num)
	if r != out {
		return r
	}

	if r-1 >= 0 {
		return r - 1
	}
	if r+1 < num {
		return r + 1
	}
	return -1
}
