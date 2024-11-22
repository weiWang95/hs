package list

type List[T comparable] struct {
	head *listNode[T]
	tail *listNode[T]
	size int
}

type listNode[T comparable] struct {
	next *listNode[T]
	data T
}

func NewList[T comparable]() *List[T] {
	return &List[T]{
		head: nil,
		tail: nil,
		size: 0,
	}
}

func (l *List[T]) Size() int {
	return l.size
}

func (l *List[T]) Add(data T) {
	node := &listNode[T]{}
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

func (l *List[T]) Remove(data T) {
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

func (l *List[T]) RemoveAt(idx int) (t T) {
	if idx < 0 || idx >= l.size {
		return
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
