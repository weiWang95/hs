package utils

import (
	"fmt"
	"testing"
)

func TestRandom(t *testing.T) {
	rate := []int{10, 20, 40, 20, 10}

	ct := map[int]int{}
	for i := 0; i < 1000; i++ {
		l := Random(rate)
		ct[l] = ct[l] + 1
	}

	fmt.Println(ct)
}
