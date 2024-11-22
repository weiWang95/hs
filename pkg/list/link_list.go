package list

import "encoding/json"

type LinkList[T comparable] struct {
	head *linkNode[T]
	tail *linkNode[T]
	size int
}

type linkNode[T comparable] struct {
	next *linkNode[T]
	data *T
}

func NewLinkList[T comparable]() *LinkList[T] {
	return &LinkList[T]{
		head: nil,
		tail: nil,
		size: 0,
	}
}

func (l *LinkList[T]) Size() int {
	return l.size
}

func (l *LinkList[T]) Each(fn func(*T) bool) (b bool) {
	b = true
	for node := l.head; node != nil; node = node.next {
		if !fn(node.data) {
			b = false
			break
		}
	}
	return
}

func (l *LinkList[T]) Add(data *T) {
	node := &linkNode[T]{}
	node.data = data
	if l.head == nil {
		l.head = node
		l.tail = node
	} else {
		l.tail.next = node
		l.tail = node
	}
	l.size++
}

func (l *LinkList[T]) Get(idx int) (t *T) {
	if idx < 0 || idx >= l.size {
		return nil
	}
	node := l.head
	for i := 0; i < idx; i++ {
		node = node.next
	}
	return node.data
}

func (l *LinkList[T]) Remove(data *T) {
	if l.head == nil {
		return
	}
	if l.head.data == data {
		l.head = l.head.next
		l.size--
		return
	}
	node := l.head
	for node.next != nil {
		if node.next.data == data {
			node.next = node.next.next
			l.size--
			return
		}
		node = node.next
	}
}

func (l *LinkList[T]) RemoveAt(idx int) (t *T) {
	if idx < 0 || idx >= l.size {
		return nil
	}
	prev := l.head
	node := l.head
	for i := 0; i < idx; i++ {
		prev = node
		node = node.next
	}
	if node == l.head {
		l.head = node.next
	}
	if node == l.tail {
		prev.next = nil
		l.tail = prev
	}
	l.size--
	return node.data
}

func (l *LinkList[T]) Clear() {
	l.head = nil
	l.tail = nil
	l.size = 0
}

func (l *LinkList[T]) MarshalJSON() ([]byte, error) {
	var nodes []*T
	for node := l.head; node != nil; node = node.next {
		nodes = append(nodes, node.data)
	}
	return json.Marshal(nodes)
}

func (l *LinkList[T]) UnmarshalJSON(data []byte) error {
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
