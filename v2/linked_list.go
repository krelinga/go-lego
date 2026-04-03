package v2

import "iter"

type LinkedListPosition[V any] struct {
	owner *LinkedList[V]
	node  *linkedListNode[V]
}

type linkedListNode[V any] struct {
	value V
	prev  *linkedListNode[V]
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

func (l *LinkedList[V]) All() iter.Seq2[LinkedListPosition[V], V] {
	return func(yield func(LinkedListPosition[V], V) bool) {
		for node := l.head; node != nil; node = node.next {
			pos := LinkedListPosition[V]{l, node}
			if !yield(pos, node.value) {
				return
			}
		}
	}
}

func (l *LinkedList[V]) Positions() iter.Seq[LinkedListPosition[V]] {
	return func(yield func(LinkedListPosition[V]) bool) {
		for node := l.head; node != nil; node = node.next {
			pos := LinkedListPosition[V]{l, node}
			if !yield(pos) {
				return
			}
		}
	}
}

func (l *LinkedList[V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for node := l.head; node != nil; node = node.next {
			if !yield(node.value) {
				return
			}
		}
	}
}

func (l *LinkedList[V]) ReverseAll() iter.Seq2[LinkedListPosition[V], V] {
	return func(yield func(LinkedListPosition[V], V) bool) {
		for node := l.tail; node != nil; node = node.prev {
			pos := LinkedListPosition[V]{l, node}
			if !yield(pos, node.value) {
				return
			}
		}
	}
}

func (l *LinkedList[V]) ReversePositions() iter.Seq[LinkedListPosition[V]] {
	return func(yield func(LinkedListPosition[V]) bool) {
		for node := l.tail; node != nil; node = node.prev {
			pos := LinkedListPosition[V]{l, node}
			if !yield(pos) {
				return
			}
		}
	}
}

func (l *LinkedList[V]) ReverseValues() iter.Seq[V] {
	return func(yield func(V) bool) {
		for node := l.tail; node != nil; node = node.prev {
			if !yield(node.value) {
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

func (l *LinkedList[V]) String() string {
	return listStringHelper(l)
}

func (l *LinkedList[V]) Set(p LinkedListPosition[V], v V) {
	if p.node == nil || p.owner != l {
		panic("position does not belong to this list")
	}
	p.node.value = v
}

func (l *LinkedList[V]) InsertBefore(p LinkedListPosition[V], v V) LinkedListPosition[V] {
	if p.node == nil || p.owner != l {
		panic("position does not belong to this list")
	}
	newNode := &linkedListNode[V]{value: v}
	if p.node == l.head {
		newNode.next = l.head
		l.head.prev = newNode
		l.head = newNode
	} else {
		prev := p.node.prev
		prev.next = newNode
		newNode.prev = prev
		newNode.next = p.node
		p.node.prev = newNode
	}
	l.len++
	return LinkedListPosition[V]{l, newNode}
}

func (l *LinkedList[V]) InsertAfter(p LinkedListPosition[V], v V) LinkedListPosition[V] {
	if p.node == nil || p.owner != l {
		panic("position does not belong to this list")
	}
	newNode := &linkedListNode[V]{value: v}
	if p.node == l.tail {
		newNode.prev = l.tail
		l.tail.next = newNode
		l.tail = newNode
	} else {
		next := p.node.next
		next.prev = newNode
		newNode.next = next
		newNode.prev = p.node
		p.node.next = newNode
	}
	l.len++
	return LinkedListPosition[V]{l, newNode}
}

func (l *LinkedList[V]) Add(v V) {
	newNode := &linkedListNode[V]{value: v}
	if l.tail == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		newNode.prev = l.tail
		l.tail = newNode
	}
	l.len++
}

func (l *LinkedList[V]) Remove(p LinkedListPosition[V]) {
	if p.node == nil || p.owner != l {
		panic("position does not belong to this list")
	}
	if p.node.prev != nil {
		p.node.prev.next = p.node.next
	} else {
		l.head = p.node.next
	}
	if p.node.next != nil {
		p.node.next.prev = p.node.prev
	} else {
		l.tail = p.node.prev
	}
	l.len--
}
