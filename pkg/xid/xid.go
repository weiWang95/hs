package xid

import (
	"math/rand"
	"time"
)

var base, _ = time.ParseInLocation("20060102150405", "20240101000000", time.UTC)

func New() uint64 {
	t := time.Now().UnixNano() - base.UnixNano()
	return uint64((t << 5) + rand.Int63n(1000))
}
