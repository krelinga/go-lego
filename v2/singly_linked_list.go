package v2

import "iter"

type SinglyLinkedListPosition[V any] struct {
	owner *SinglyLinkedList[V]
	node  *singlyLinkedListNode[V]
}

type singlyLinkedListNode[V any] struct {
	value V
	next  *singlyLinkedListNode[V]
}

type SinglyLinkedList[V any] struct {
	head *singlyLinkedListNode[V]
	tail *singlyLinkedListNode[V]
	len  int
}

func (l *SinglyLinkedList[V]) Length() int {
	return l.len
}

func (l *SinglyLinkedList[V]) Get(p SinglyLinkedListPosition[V]) (V, bool) {
	if p.node == nil || p.owner != l {
		var zero V
		return zero, false
	}
	return p.node.value, true
}

func (l *SinglyLinkedList[V]) All() iter.Seq2[SinglyLinkedListPosition[V], V] {
	return func(yield func(SinglyLinkedListPosition[V], V) bool) {
		for node := l.head; node != nil; node = node.next {
			pos := SinglyLinkedListPosition[V]{l, node}
			if !yield(pos, node.value) {
				return
			}
		}
	}
}

func (l *SinglyLinkedList[V]) Positions() iter.Seq[SinglyLinkedListPosition[V]] {
	return func(yield func(SinglyLinkedListPosition[V]) bool) {
		for node := l.head; node != nil; node = node.next {
			pos := SinglyLinkedListPosition[V]{l, node}
			if !yield(pos) {
				return
			}
		}
	}
}

func (l *SinglyLinkedList[V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for node := l.head; node != nil; node = node.next {
			if !yield(node.value) {
				return
			}
		}
	}
}

func (l *SinglyLinkedList[V]) First() (SinglyLinkedListPosition[V], V, bool) {
	if l.head == nil {
		var zero V
		return SinglyLinkedListPosition[V]{l, nil}, zero, false
	}
	return SinglyLinkedListPosition[V]{l, l.head}, l.head.value, true
}

func (l *SinglyLinkedList[V]) Last() (SinglyLinkedListPosition[V], V, bool) {
	if l.tail == nil {
		var zero V
		return SinglyLinkedListPosition[V]{l, nil}, zero, false
	}
	return SinglyLinkedListPosition[V]{l, l.tail}, l.tail.value, true
}

func (l *SinglyLinkedList[V]) String() string {
	return listStringHelper(l)
}

func (l *SinglyLinkedList[V]) Set(p SinglyLinkedListPosition[V], v V) {
	if p.node == nil || p.owner != l {
		panic("invalid position")
	}
	p.node.value = v
}

func (l *SinglyLinkedList[V]) InsertBefore(p SinglyLinkedListPosition[V], v V) SinglyLinkedListPosition[V] {
	if p.node == nil || p.owner != l {
		panic("invalid position")
	}
	newNode := &singlyLinkedListNode[V]{value: v}
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
	return SinglyLinkedListPosition[V]{l, newNode}
}

func (l *SinglyLinkedList[V]) InsertAfter(p SinglyLinkedListPosition[V], v V) SinglyLinkedListPosition[V] {
	if p.node == nil || p.owner != l {
		panic("invalid position")
	}
	newNode := &singlyLinkedListNode[V]{value: v}
	newNode.next = p.node.next
	p.node.next = newNode
	if p.node == l.tail {
		l.tail = newNode
	}
	l.len++
	return SinglyLinkedListPosition[V]{l, newNode}
}

func (l *SinglyLinkedList[V]) Append(v V) SinglyLinkedListPosition[V] {
	newNode := &singlyLinkedListNode[V]{value: v}
	if l.tail == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		l.tail = newNode
	}
	l.len++
	return SinglyLinkedListPosition[V]{l, newNode}
}

func (l *SinglyLinkedList[V]) Remove(p SinglyLinkedListPosition[V]) {
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

func (l *SinglyLinkedList[V]) Add(v V) {
	l.Append(v)
}
