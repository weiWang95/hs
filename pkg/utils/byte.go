package utils

import "encoding/binary"

func Uint64ToBytes(a uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, a)
	return b
}

func BytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func Uint8ToBytes(a uint8) []byte {
	return []byte{a}
}

func BytesToUint8(b []byte) uint8 {
	return b[0]
}
