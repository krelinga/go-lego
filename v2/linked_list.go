package v2

import "iter"

type LinkedListPosition struct {
	owner any
	node  any
}

func (p LinkedListPosition) IsValid() bool {
	return p.node != nil && p.owner != nil
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

func (l *LinkedList[V]) pos2node(p LinkedListPosition) *linkedListNode[V] {
	if p.node == nil || p.owner == nil {
		return nil
	}
	listPtr, ok := p.owner.(*LinkedList[V])
	if !ok || listPtr != l {
		return nil
	}
	nodePtr, ok := p.node.(*linkedListNode[V])
	if !ok {
		return nil
	}
	return nodePtr
}

func (l *LinkedList[V]) Get(p LinkedListPosition) (V, bool) {
	node := l.pos2node(p)
	if node == nil {
		var zero V
		return zero, false
	}
	return node.value, true
}

func (l *LinkedList[V]) All() iter.Seq2[LinkedListPosition, V] {
	return func(yield func(LinkedListPosition, V) bool) {
		for node := l.head; node != nil; node = node.next {
			pos := LinkedListPosition{l, node}
			if !yield(pos, node.value) {
				return
			}
		}
	}
}

func (l *LinkedList[V]) Positions() iter.Seq[LinkedListPosition] {
	return func(yield func(LinkedListPosition) bool) {
		for node := l.head; node != nil; node = node.next {
			pos := LinkedListPosition{l, node}
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

func (l *LinkedList[V]) ReverseAll() iter.Seq2[LinkedListPosition, V] {
	return func(yield func(LinkedListPosition, V) bool) {
		for node := l.tail; node != nil; node = node.prev {
			pos := LinkedListPosition{l, node}
			if !yield(pos, node.value) {
				return
			}
		}
	}
}

func (l *LinkedList[V]) ReversePositions() iter.Seq[LinkedListPosition] {
	return func(yield func(LinkedListPosition) bool) {
		for node := l.tail; node != nil; node = node.prev {
			pos := LinkedListPosition{l, node}
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

func (l *LinkedList[V]) First() (LinkedListPosition, V, bool) {
	if l.head == nil {
		var zero V
		return LinkedListPosition{}, zero, false
	}
	return LinkedListPosition{l, l.head}, l.head.value, true
}

func (l *LinkedList[V]) Last() (LinkedListPosition, V, bool) {
	if l.tail == nil {
		var zero V
		return LinkedListPosition{}, zero, false
	}
	return LinkedListPosition{l, l.tail}, l.tail.value, true
}

func (l *LinkedList[V]) String() string {
	return listStringHelper(l)
}

func (l *LinkedList[V]) Set(p LinkedListPosition, v V) {
	node := l.pos2node(p)
	if node == nil {
		panic("position does not belong to this list")
	}
	node.value = v
}

func (l *LinkedList[V]) InsertBefore(p LinkedListPosition, v V) LinkedListPosition {
	if l.Length() == 0 && !p.IsValid() {
		l.Add(v)
		return LinkedListPosition{l, l.head}
	}
	node := l.pos2node(p)
	if node == nil {
		panic("position does not belong to this list")
	}
	newNode := &linkedListNode[V]{value: v}
	if node == l.head {
		newNode.next = l.head
		l.head.prev = newNode
		l.head = newNode
	} else {
		prev := node.prev
		prev.next = newNode
		newNode.prev = prev
		newNode.next = node
		node.prev = newNode
	}
	l.len++
	return LinkedListPosition{l, newNode}
}

func (l *LinkedList[V]) InsertAfter(p LinkedListPosition, v V) LinkedListPosition {
	if l.Length() == 0 && !p.IsValid() {
		l.Add(v)
		return LinkedListPosition{l, l.head}
	}
	node := l.pos2node(p)
	if node == nil {
		panic("position does not belong to this list")
	}
	newNode := &linkedListNode[V]{value: v}
	if node == l.tail {
		newNode.prev = l.tail
		l.tail.next = newNode
		l.tail = newNode
	} else {
		next := node.next
		next.prev = newNode
		newNode.next = next
		newNode.prev = node
		node.next = newNode
	}
	l.len++
	return LinkedListPosition{l, newNode}
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

func (l *LinkedList[V]) Remove(p LinkedListPosition) {
	node := l.pos2node(p)
	if node == nil {
		panic("position does not belong to this list")
	}
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}
	l.len--
}

func (l *LinkedList[V]) Clear() {
	l.head = nil
	l.tail = nil
	l.len = 0
}