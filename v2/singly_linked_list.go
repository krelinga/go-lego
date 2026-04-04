package v2

import "iter"

type SinglyLinkedListPosition struct {
	owner any
	node  any
}

func (p SinglyLinkedListPosition) IsValid() bool {
	return p.node != nil && p.owner != nil
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

func (l *SinglyLinkedList[V]) pos2node(p SinglyLinkedListPosition) *singlyLinkedListNode[V] {
	if p.node == nil || p.owner == nil {
		return nil
	}
	listPtr, ok := p.owner.(*SinglyLinkedList[V])
	if !ok || listPtr != l {
		return nil
	}
	nodePtr, ok := p.node.(*singlyLinkedListNode[V])
	if !ok {
		return nil
	}
	return nodePtr
}

func (l *SinglyLinkedList[V]) Get(p SinglyLinkedListPosition) (V, bool) {
	node := l.pos2node(p)
	if node == nil {
		var zero V
		return zero, false
	}
	return node.value, true
}

func (l *SinglyLinkedList[V]) All() iter.Seq2[SinglyLinkedListPosition, V] {
	return func(yield func(SinglyLinkedListPosition, V) bool) {
		for node := l.head; node != nil; node = node.next {
			pos := SinglyLinkedListPosition{l, node}
			if !yield(pos, node.value) {
				return
			}
		}
	}
}

func (l *SinglyLinkedList[V]) Positions() iter.Seq[SinglyLinkedListPosition] {
	return func(yield func(SinglyLinkedListPosition) bool) {
		for node := l.head; node != nil; node = node.next {
			pos := SinglyLinkedListPosition{l, node}
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

func (l *SinglyLinkedList[V]) First() (SinglyLinkedListPosition, V, bool) {
	if l.head == nil {
		var zero V
		return SinglyLinkedListPosition{}, zero, false
	}
	return SinglyLinkedListPosition{l, l.head}, l.head.value, true
}

func (l *SinglyLinkedList[V]) Last() (SinglyLinkedListPosition, V, bool) {
	if l.tail == nil {
		var zero V
		return SinglyLinkedListPosition{}, zero, false
	}
	return SinglyLinkedListPosition{l, l.tail}, l.tail.value, true
}

func (l *SinglyLinkedList[V]) String() string {
	return listStringHelper(l)
}

func (l *SinglyLinkedList[V]) Set(p SinglyLinkedListPosition, v V) {
	node := l.pos2node(p)
	if node == nil {
		panic("position does not belong to this list")
	}
	node.value = v
}

func (l *SinglyLinkedList[V]) InsertBefore(p SinglyLinkedListPosition, v V) SinglyLinkedListPosition {
	if l.Length() == 0 && !p.IsValid() {
		l.Add(v)
		return SinglyLinkedListPosition{l, l.head}
	}
	node := l.pos2node(p)
	if node == nil {
		panic("position does not belong to this list")
	}
	newNode := &singlyLinkedListNode[V]{value: v}
	if node == l.head {
		newNode.next = l.head
		l.head = newNode
	} else {
		prev := l.head
		for prev != nil && prev.next != node {
			prev = prev.next
		}
		if prev == nil {
			panic("position not found in list")
		}
		prev.next = newNode
		newNode.next = node
	}
	l.len++
	return SinglyLinkedListPosition{l, newNode}
}

func (l *SinglyLinkedList[V]) InsertAfter(p SinglyLinkedListPosition, v V) SinglyLinkedListPosition {
	if l.Length() == 0 && !p.IsValid() {
		l.Add(v)
		return SinglyLinkedListPosition{l, l.head}
	}
	node := l.pos2node(p)
	if node == nil {
		panic("position does not belong to this list")
	}
	newNode := &singlyLinkedListNode[V]{value: v}
	newNode.next = node.next
	node.next = newNode
	if node == l.tail {
		l.tail = newNode
	}
	l.len++
	return SinglyLinkedListPosition{l, newNode}
}

func (l *SinglyLinkedList[V]) Add(v V) {
	newNode := &singlyLinkedListNode[V]{value: v}
	if l.tail == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		l.tail = newNode
	}
	l.len++
}

func (l *SinglyLinkedList[V]) Remove(p SinglyLinkedListPosition) {
	node := l.pos2node(p)
	if node == nil {
		panic("position does not belong to this list")
	}
	if node == l.head {
		l.head = node.next
	} else {
		prev := l.head
		for prev != nil && prev.next != node {
			prev = prev.next
		}
		if prev == nil {
			panic("position not found in list")
		}
		prev.next = node.next
		if node == l.tail {
			l.tail = prev
		}
	}
	l.len--
}

func (l *SinglyLinkedList[V]) Clear() {
	l.head = nil
	l.tail = nil
	l.len = 0
}
