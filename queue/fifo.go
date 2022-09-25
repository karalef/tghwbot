package queue

import "sync"

type lnode[T any] struct {
	v    T
	next *lnode[T]
}

type fifo[T any] struct {
	head *lnode[T]
	tail *lnode[T]
	size int
	mut  sync.Mutex
}

func (f *fifo[T]) len() int {
	f.mut.Lock()
	defer f.mut.Unlock()
	return f.size
}

func (f *fifo[T]) push(v T) {
	f.mut.Lock()
	defer f.mut.Unlock()
	if f.size == 0 {
		f.tail = &lnode[T]{v: v}
		f.head = f.tail
		f.size++
		return
	}
	f.tail.next = &lnode[T]{v: v}
	f.tail = f.tail.next
	f.size++
}

func (f *fifo[T]) peek() T {
	f.mut.Lock()
	defer f.mut.Unlock()
	return f.head.v
}

func (f *fifo[T]) remove() {
	f.mut.Lock()
	if f.size > 0 {
		f.head = f.head.next
		f.size--
	}
	f.mut.Unlock()
}
