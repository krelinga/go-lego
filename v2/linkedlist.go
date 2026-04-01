package v2

import "iter"

type LinkedListPosition[V any] struct {
	owner *LinkedList[V]
	node  *linkedListNode[V]
}

type linkedListNode[V any] struct {
	value V
	next  *linkedListNode[V]
}

type LinkedList[V any] struct {
	head *linkedListNode[V]
	tail *linkedListNode[V]
	len  int
}

func (l *LinkedList[V]) Length() int {
	return l.len
}

func (l *LinkedList[V]) Get(p LinkedListPosition[V]) (V, bool) {
	if p.node == nil || p.owner != l {
		var zero V
		return zero, false
	}
	return p.node.value, true
}

func (l *LinkedList[V]) List() iter.Seq2[LinkedListPosition[V], V] {
	return func(yield func(LinkedListPosition[V], V) bool) {
		for node := l.head; node != nil; node = node.next {
			pos := LinkedListPosition[V]{l, node}
			if !yield(pos, node.value) {
				return
			}
		}
	}
}

func (l *LinkedList[V]) First() (LinkedListPosition[V], V, bool) {
	if l.head == nil {
		var zero V
		return LinkedListPosition[V]{l, nil}, zero, false
	}
	return LinkedListPosition[V]{l, l.head}, l.head.value, true
}

func (l *LinkedList[V]) Last() (LinkedListPosition[V], V, bool) {
	if l.tail == nil {
		var zero V
		return LinkedListPosition[V]{l, nil}, zero, false
	}
	return LinkedListPosition[V]{l, l.tail}, l.tail.value, true
}

func (l *LinkedList[V]) Set(p LinkedListPosition[V], v V) {
	if p.node == nil || p.owner != l {
		panic("invalid position")
	}
	p.node.value = v
}

func (l *LinkedList[V]) InsertBefore(p LinkedListPosition[V], v V) LinkedListPosition[V] {
	if p.node == nil || p.owner != l {
		panic("invalid position")
	}
	newNode := &linkedListNode[V]{value: v}
	if p.node == l.head {
		newNode.next = l.head
		l.head = newNode
	} else {
		prev := l.head
		for prev != nil && prev.next != p.node {
			prev = prev.next
		}
		if prev == nil {
			panic("position not found in list")
		}
		prev.next = newNode
		newNode.next = p.node
	}
	l.len++
	return LinkedListPosition[V]{l, newNode}
}

func (l *LinkedList[V]) InsertAfter(p LinkedListPosition[V], v V) LinkedListPosition[V] {
	if p.node == nil || p.owner != l {
		panic("invalid position")
	}
	newNode := &linkedListNode[V]{value: v}
	newNode.next = p.node.next
	p.node.next = newNode
	if p.node == l.tail {
		l.tail = newNode
	}
	l.len++
	return LinkedListPosition[V]{l, newNode}
}

func (l *LinkedList[V]) Append(v V) LinkedListPosition[V] {
	newNode := &linkedListNode[V]{value: v}
	if l.tail == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		l.tail = newNode
	}
	l.len++
	return LinkedListPosition[V]{l, newNode}
}

func (l *LinkedList[V]) Remove(p LinkedListPosition[V]) {
	if p.node == nil || p.owner != l {
		panic("invalid position")
	}
	if p.node == l.head {
		l.head = p.node.next
	} else {
		prev := l.head
		for prev != nil && prev.next != p.node {
			prev = prev.next
		}
		if prev == nil {
			panic("position not found in list")
		}
		prev.next = p.node.next
		if p.node == l.tail {
			l.tail = prev
		}
	}
	l.len--
}

func NewLinkedList[V any](entries ...V) *LinkedList[V] {
	l := &LinkedList[V]{}
	for _, v := range entries {
		l.Append(v)
	}
	return l
}