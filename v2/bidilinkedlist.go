package v2

type BidiLinkedListPosition[V any] struct {
	owner *BidiLinkedList[V]
	node  *bidiLinkedListNode[V]
}

type bidiLinkedListNode[V any] struct {
	value V
	prev  *bidiLinkedListNode[V]
	next  *bidiLinkedListNode[V]
}

type FixedBidiLinkedList[V any] interface {
	FixedList[BidiLinkedListPosition[V], V]
	ReverseRange() ListSeq[BidiLinkedListPosition[V], V]
}

type BidiLinkedList[V any] struct {
	head *bidiLinkedListNode[V]
	tail *bidiLinkedListNode[V]
	len  int
}

func (l *BidiLinkedList[V]) Length() int {
	return l.len
}

func (l *BidiLinkedList[V]) Get(p BidiLinkedListPosition[V]) (V, bool) {
	if p.node == nil || p.owner != l {
		var zero V
		return zero, false
	}
	return p.node.value, true
}

func (l *BidiLinkedList[V]) Range() ListSeq[BidiLinkedListPosition[V], V] {
	return func(yield func(BidiLinkedListPosition[V], V) bool) {
		for node := l.head; node != nil; node = node.next {
			pos := BidiLinkedListPosition[V]{l, node}
			if !yield(pos, node.value) {
				return
			}
		}
	}
}

func (l *BidiLinkedList[V]) ReverseRange() ListSeq[BidiLinkedListPosition[V], V] {
	return func(yield func(BidiLinkedListPosition[V], V) bool) {
		for node := l.tail; node != nil; node = node.prev {
			pos := BidiLinkedListPosition[V]{l, node}
			if !yield(pos, node.value) {
				return
			}
		}
	}
}

func (l *BidiLinkedList[V]) First() (BidiLinkedListPosition[V], V, bool) {
	if l.head == nil {
		var zero V
		return BidiLinkedListPosition[V]{l, nil}, zero, false
	}
	return BidiLinkedListPosition[V]{l, l.head}, l.head.value, true
}

func (l *BidiLinkedList[V]) Last() (BidiLinkedListPosition[V], V, bool) {
	if l.tail == nil {
		var zero V
		return BidiLinkedListPosition[V]{l, nil}, zero, false
	}
	return BidiLinkedListPosition[V]{l, l.tail}, l.tail.value, true
}

func (l *BidiLinkedList[V]) Set(p BidiLinkedListPosition[V], v V) {
	if p.node == nil || p.owner != l {
		panic("position does not belong to this list")
	}
	p.node.value = v
}

func (l *BidiLinkedList[V]) InsertBefore(p BidiLinkedListPosition[V], v V) BidiLinkedListPosition[V] {
	if p.node == nil || p.owner != l {
		panic("position does not belong to this list")
	}
	newNode := &bidiLinkedListNode[V]{value: v}
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
	return BidiLinkedListPosition[V]{l, newNode}
}

func (l *BidiLinkedList[V]) InsertAfter(p BidiLinkedListPosition[V], v V) BidiLinkedListPosition[V] {
	if p.node == nil || p.owner != l {
		panic("position does not belong to this list")
	}
	newNode := &bidiLinkedListNode[V]{value: v}
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
	return BidiLinkedListPosition[V]{l, newNode}
}

func (l *BidiLinkedList[V]) Append(v V) BidiLinkedListPosition[V] {
	newNode := &bidiLinkedListNode[V]{value: v}
	if l.tail == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		newNode.prev = l.tail
		l.tail = newNode
	}
	l.len++
	return BidiLinkedListPosition[V]{l, newNode}
}

func (l *BidiLinkedList[V]) Remove(p BidiLinkedListPosition[V]) {
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

func NewBidiLinkedList[V any](entries ...V) *BidiLinkedList[V] {
	l := &BidiLinkedList[V]{}
	for _, entry := range entries {
		l.Append(entry)
	}
	return l
}
