package protocol

import (
	"fmt"
	"testing"
)

func TestHead_Timestamp(t *testing.T) {
	h := NewHead()
	h.Server().SetTimestamp()

	fmt.Println("----", h.Timestamp().Format("2006-01-02 15:04:05"))
	t.Fail()
}
