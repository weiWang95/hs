package list

import (
	"fmt"
	"testing"
)

type B struct {
	Name string
}

func TestDoubleLinkList_Clone(t *testing.T) {

	b1 := &B{Name: "a"}
	b2 := &B{Name: "b"}
	fmt.Printf("b1: %+v %T, b2: %+v %T\n", b1, b1, b2, b2)

	l := NewDoubleLinkList[B]()
	l.Add(b1)
	l.Add(b2)

	l2 := l.Clone()
	b3 := l2.Get(0)
	b4 := l2.Get(1)

	b3.Name = "c"

	fmt.Printf("b3: %+v %T, b4: %+v %T\n", b3, b3, b4, b4)
	t.Fail()
}
