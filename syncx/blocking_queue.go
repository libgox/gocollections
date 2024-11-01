package syncx

import "sync"

type BlockingQueue[T any] interface {
	// Put enqueue one item, block if the queue is full
	Put(item T)

	// Take dequeue one item, block until it's available
	Take() T

	// Poll dequeue one item, return zero value if queue is empty
	Poll() T

	// CompareAndPoll compare the first item and poll it if meet the conditions
	CompareAndPoll(compare func(item T) bool) T

	// Peek return the first item without dequeuing, return zero value if queue is empty
	Peek() T

	// PeekLast return last item in queue without dequeuing, return zero value if queue is empty
	PeekLast() T

	// Size return the current size of the queue
	Size() int

	// ReadableSlice returns a new view of the readable items in the queue
	ReadableSlice() []T
}

type blockingQueue[T any] struct {
	items   []T
	headIdx int
	tailIdx int
	size    int
	maxSize int

	mutex      sync.Mutex
	isNotEmpty *sync.Cond
	isNotFull  *sync.Cond
}

func NewBlockingQueue[T any](maxSize int) BlockingQueue[T] {
	bq := &blockingQueue[T]{
		items:   make([]T, maxSize),
		headIdx: 0,
		tailIdx: 0,
		size:    0,
		maxSize: maxSize,
	}

	bq.isNotEmpty = sync.NewCond(&bq.mutex)
	bq.isNotFull = sync.NewCond(&bq.mutex)
	return bq
}

func (bq *blockingQueue[T]) Put(item T) {
	bq.mutex.Lock()
	defer bq.mutex.Unlock()

	for bq.size == bq.maxSize {
		bq.isNotFull.Wait()
	}

	wasEmpty := bq.size == 0

	bq.items[bq.tailIdx] = item
	bq.size++
	bq.tailIdx++
	if bq.tailIdx >= bq.maxSize {
		bq.tailIdx = 0
	}

	if wasEmpty {
		// Wake up eventual reader waiting for next item
		bq.isNotEmpty.Signal()
	}
}

func (bq *blockingQueue[T]) Take() T {
	bq.mutex.Lock()
	defer bq.mutex.Unlock()

	for bq.size == 0 {
		bq.isNotEmpty.Wait()
	}

	return bq.dequeue()
}

func (bq *blockingQueue[T]) Poll() T {
	bq.mutex.Lock()
	defer bq.mutex.Unlock()

	if bq.size == 0 {
		var zero T
		return zero
	}

	return bq.dequeue()
}

func (bq *blockingQueue[T]) CompareAndPoll(compare func(T) bool) T {
	bq.mutex.Lock()
	defer bq.mutex.Unlock()

	if bq.size == 0 {
		var zero T
		return zero
	}

	if compare(bq.items[bq.headIdx]) {
		return bq.dequeue()
	}
	var zero T
	return zero
}

func (bq *blockingQueue[T]) Peek() T {
	bq.mutex.Lock()
	defer bq.mutex.Unlock()

	if bq.size == 0 {
		var zero T
		return zero
	}
	return bq.items[bq.headIdx]
}

func (bq *blockingQueue[T]) PeekLast() T {
	bq.mutex.Lock()
	defer bq.mutex.Unlock()

	if bq.size == 0 {
		var zero T
		return zero
	}
	idx := (bq.headIdx + bq.size - 1) % bq.maxSize
	return bq.items[idx]
}

func (bq *blockingQueue[T]) dequeue() T {
	item := bq.items[bq.headIdx]
	var zero T
	bq.items[bq.headIdx] = zero

	bq.headIdx++
	if bq.headIdx == len(bq.items) {
		bq.headIdx = 0
	}

	bq.size--
	bq.isNotFull.Signal()
	return item
}

func (bq *blockingQueue[T]) Size() int {
	bq.mutex.Lock()
	defer bq.mutex.Unlock()

	return bq.size
}

func (bq *blockingQueue[T]) ReadableSlice() []T {
	bq.mutex.Lock()
	defer bq.mutex.Unlock()

	res := make([]T, bq.size)
	readIdx := bq.headIdx
	for i := 0; i < bq.size; i++ {
		res[i] = bq.items[readIdx]
		readIdx++
		if readIdx == bq.maxSize {
			readIdx = 0
		}
	}

	return res
}
