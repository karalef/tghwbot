package queue

// New creates new queue.
func New[T any](worker func(T)) *Queue[T] {
	if worker == nil {
		return nil
	}
	q := Queue[T]{
		push: make(chan T),
		out:  make(chan T),
		buf:  new(fifo[T]),
		do:   worker,
	}
	go q.buffer()
	go q.listen()
	return &q
}

// Queue starts the worker for each value sent to Push.
type Queue[T any] struct {
	push chan T
	out  chan T
	buf  *fifo[T]
	do   func(T)
}

func (q *Queue[T]) Push(v T) {
	q.push <- v
}

func (q *Queue[T]) Len() int {
	return q.buf.len()
}

func (q *Queue[T]) Close() {
	close(q.push)
}

func (q *Queue[T]) listen() {
	for {
		r, ok := <-q.out
		if !ok {
			return
		}
		q.do(r)
	}
}

func (q *Queue[T]) buffer() {
	in := q.push
	var out chan T
	var val T

	for in != nil || out != nil {
		select {
		case elem, ok := <-in:
			if !ok {
				in = nil
				break
			}
			q.buf.push(elem)
		case out <- val:
			q.buf.remove()
		}

		if q.buf.len() == 0 {
			out = nil
			continue
		}
		out = q.out
		val = q.buf.peek()
	}

	close(q.out)
}
