package queue

// New creates new queue.
// count is the maximum number of concurrently running workers.
func New[T any](worker func(T), count ...int) *Queue[T] {
	c := 1
	if len(count) > 0 {
		c = count[0]
	}
	if worker == nil || c < 1 {
		return nil
	}
	q := Queue[T]{
		push: make(chan T),
		pool: make(chan func(T), c),
		buf:  new(fifo[T]),
	}
	for i := 0; i < c; i++ {
		q.pool <- worker
	}
	go q.buffer()
	return &q
}

// Queue starts the worker for each value sent to Push.
type Queue[T any] struct {
	push chan T
	pool chan func(T)
	buf  *fifo[T]
}

// Push pushes the value to the queue and returns true if the value passed to the handler without waiting.
func (q *Queue[T]) Push(v T) bool {
	a := cap(q.pool)-len(q.pool) > 0
	q.push <- v
	return a
}

func (q *Queue[T]) Len() int {
	return q.buf.len()
}

func (q *Queue[T]) Close() {
	close(q.push)
}

func (q *Queue[T]) buffer() {
	var pool chan func(T)
	var val T
	for {
		select {
		case elem, ok := <-q.push:
			if !ok {
				close(q.pool)
				return
			}
			q.buf.push(elem)
		case worker := <-pool:
			go func() {
				worker(val)
				q.pool <- worker
			}()
			q.buf.remove()
		}

		if q.buf.len() == 0 {
			pool = nil
			continue
		}
		pool = q.pool
		val = q.buf.peek()
	}
}
