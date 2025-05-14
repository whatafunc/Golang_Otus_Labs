package hw04lrucache

const Debug = false // mode => go test -run '^TestList$'
type List interface {
	Len() int
	Front() *ListItem // head
	Back() *ListItem  // tail
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct { // It is simply one entry in the chain..
	Value interface{} // ..with Value pf any type like struct/int/..
	Next  *ListItem   // A pointer to the next item in the list.
	Prev  *ListItem
}

type list struct { // Doubly Linked List (Primary for list operations).
	head *ListItem
	tail *ListItem
	len  int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  l.head, // point forward to the old head
	}

	if l.head != nil {
		l.head.Prev = newItem // item gets as first, and ex-first goes to last
	} else {
		l.tail = newItem // first item added is both head and tail
	}

	l.head = newItem
	l.len++

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Prev:  l.tail, // point backward to the old head
	}
	if l.tail != nil {
		l.tail.Next = newItem // item gets as Last, and ex-last goes ...
	} else {
		l.head = newItem // first item added is both head and tail
	}

	l.tail = newItem
	l.len++

	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	// If the item has a previous node, update its Next pointer
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}

	// If the item has a next node, update its Prev pointer
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}

	l.len--

	// Optionally clear pointers to help GC
	i.Next = nil
	i.Prev = nil
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || l.head == i {
		return
	}

	// Unlink i from current position
	// if i.Prev != nil { // начальный существует для всех тк если начальный является стартовым то это отсечено выше уже
	i.Prev.Next = i.Next
	//}
	// if i.Next != nil { // последний существует всегда
	//	i.Next.Prev = i.Prev
	//}

	// If i was tail, update tail
	if l.tail == i {
		l.tail = i.Prev
	}

	// Insert i at front
	i.Prev = nil
	i.Next = l.head
	// if l.head != nil { // первый же есть всегда
	l.head.Prev = i
	//}
	l.head = i

	// If the list was empty (unlikely here), set tail too
	if l.tail == nil {
		l.tail = i
	}
}
