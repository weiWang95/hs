package list

import (
	"encoding/json"
	"reflect"
)

type DoubleLinkList[T any] struct {
	head *node[T]
	tail *node[T]
	size int
}

type node[T any] struct {
	prev *node[T]
	next *node[T]
	data *T
}

func (n *node[T]) Clear() {
	if n == nil {
		return
	}
	n.prev = nil
	n.next = nil
}

func NewDoubleLinkList[T any]() *DoubleLinkList[T] {
	return &DoubleLinkList[T]{}
}

func (l *DoubleLinkList[T]) Size() int {
	return l.size
}

func (l *DoubleLinkList[T]) Each(fn func(*T) bool) {
	for node := l.head; node != nil; node = node.next {
		if !fn(node.data) {
			break
		}
	}
}

func (l *DoubleLinkList[T]) EachWithIdx(fn func(int, *T) bool) {
	i := 0
	for node := l.head; node != nil; node = node.next {
		if !fn(i, node.data) {
			break
		}
		i += 1
	}
}

func (l *DoubleLinkList[T]) Add(data *T) {
	node := &node[T]{}

	node.data = data
	if l.head == nil {
		l.head = node
		l.tail = node
	} else {
		tail := l.tail
		tail.next = node
		node.prev = tail
		l.tail = node
	}
	l.size++
}

// func (l *DoubleLinkList[T]) Remove(data *T) {
// 	for node := l.head; node != nil; node = node.next {
// 		if node.data == data {
// 			if node == l.head {
// 				l.head = l.head.next
// 			}
// 			if node == l.tail {
// 				l.tail = l.tail.prev
// 			}
// 			if node.prev != nil {
// 				node.prev.next = node.next
// 			}
// 			if node.next != nil {
// 				node.next.prev = node.prev
// 			}
// 			l.size--

// 			node.Clear()

// 			break
// 		}
// 	}
// }

func (l *DoubleLinkList[T]) Get(index int) (t *T) {
	if index < 0 || index >= l.size {
		return nil
	}
	node := l.head
	for i := 0; i < index; i++ {
		node = node.next
	}
	return node.data
}

func (l *DoubleLinkList[T]) AddAt(index int, data *T) {
	if index < 0 || index > l.size {
		return
	}
	if l.size == 0 || index == l.size {
		l.Add(data)
		return
	}
	newNode := &node[T]{}
	newNode.data = data

	if index == 0 {
		l.head.prev = newNode
		newNode.next = l.head
		l.head = newNode
		l.size++
		return
	}

	n := l.head
	for i := 1; i < index; i++ {
		n = n.next
	}

	if n.next != nil {
		n.next.prev = newNode
		newNode.next = n.next
	}
	n.next = newNode
	newNode.prev = n

	l.size++
}

func (l *DoubleLinkList[T]) Del(index int) (t *T) {
	if index < 0 || index >= l.size {
		return nil
	}
	node := l.head
	for i := 0; i < index; i++ {
		node = node.next
	}
	if node == l.head {
		l.head = l.head.next
	}
	if node == l.tail {
		l.tail = l.tail.prev
	}
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	l.size--

	v := node.data

	node.Clear()

	return v
}

func (l *DoubleLinkList[T]) Move(from, to int) {
	if from < 0 || from >= l.size || to < 0 || to >= l.size || from == to {
		return
	}

	v := l.Del(from)
	l.AddAt(to, v)
}

func (l *DoubleLinkList[T]) Clear() {
	for node := l.head; node != nil; {
		if node.next != nil {
			node.prev.Clear()
			node = node.next
		} else {
			node.Clear()
			break
		}
	}

	l.head = nil
	l.tail = nil
	l.size = 0
}

func (l *DoubleLinkList[T]) MarshalJSON() ([]byte, error) {
	var nodes []*T
	for node := l.head; node != nil; node = node.next {
		nodes = append(nodes, node.data)
	}
	return json.Marshal(nodes)
}

func (l *DoubleLinkList[T]) UnmarshalJSON(data []byte) error {
	var v []*T

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	l.Clear()
	for _, d := range v {
		l.Add(d)
	}

	return nil
}

func (l *DoubleLinkList[T]) Clone() *DoubleLinkList[T] {
	list := NewDoubleLinkList[T]()
	for node := l.head; node != nil; node = node.next {
		t := reflect.TypeOf(node.data)
		if t.Kind() == reflect.Ptr {
			t2 := reflect.New(t.Elem())
			v := reflect.ValueOf(node.data)
			t2.Elem().Set(v.Elem())
			list.Add(t2.Interface().(*T))
		} else {
			list.Add(node.data)
		}
	}
	return list
}
