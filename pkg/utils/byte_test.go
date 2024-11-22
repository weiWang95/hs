package utils

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestM(t *testing.T) {
	b1 := Uint64ToBytes(2<<8 + 1<<0)
	b2 := []byte{8}
	fmt.Println(b1, b1)
	b3 := append(b1, b2...)

	i1 := BytesToUint64(b3[:8])
	fmt.Println(i1)

	i := binary.BigEndian.Uint64(b3[:8])
	t1 := uint8(i >> 8)
	a1 := uint8(i >> 0)
	fmt.Println(i, t1, a1, b3[8:])
}
