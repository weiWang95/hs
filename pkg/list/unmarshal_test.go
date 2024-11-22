package list

import (
	"encoding/json"
	"fmt"
	"testing"
)

type A struct {
	Name string `json:"name"`
}

func Test_unmarshal(t *testing.T) {
	RegisterValue(A{})

	l := NewDoubleLinkList[A]()
	l.Add(&A{Name: "a"})

	bs, err := json.Marshal(l)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bs))

	l2 := NewDoubleLinkList[A]()
	err = json.Unmarshal(bs, &l2)
	if err != nil {
		t.Fatal(err)
	}

	t.Fail()
}
