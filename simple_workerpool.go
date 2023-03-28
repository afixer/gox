package gox

import "sync"

// The workerPool type is a concurrent worker pool,
// that processes a list of elements E and returns a result of type R.
type simpleWorkerPool[E any, R any] struct {
	inputChannel     chan struct{}  // Channel to limit the number of concurrent workers.
	workerFunc       func(E) error  // Function that processes an element and returns an error.
	workersWaitGroup sync.WaitGroup // WaitGroup to wait for all workers to complete.
}

// NewWorkerPool creates a new worker pool with nWorkers concurrent workers.
// The workerFunc function is called for each element in the pool.
// The resultConsumer function is called for each result returned by the workerFunc.
func NewSimpleWorkerPool[E, R any](nWorkers int, workerFunc func(E) error) *simpleWorkerPool[E, R] {
	// Create a new simpleWorkerPool struct.
	p := &simpleWorkerPool[E, R]{
		inputChannel:     make(chan struct{}, nWorkers), // Buffered channel to limit the number of concurrent workers.
		workerFunc:       workerFunc,                    // Function to process an element.
		workersWaitGroup: sync.WaitGroup{},              // WaitGroup to wait for all workers to complete.
	}

	return p
}

// Add adds a new element e to the pool.
func (p *simpleWorkerPool[E, R]) Add(e E) {
	// Signal that a worker has started and increment the workersWaitGroup counter.
	p.inputChannel <- struct{}{}
	p.workersWaitGroup.Add(1)

	// Launch a new goroutine to process the element.
	go func() {
		// Decrement the workersWaitGroup counter and signal that the worker has finished when it is done.
		defer func() {
			<-p.inputChannel
			p.workersWaitGroup.Done()
		}()

		// Call the workerFunc function to process the element.
		err := p.workerFunc(e)
		if err != nil {
			return
		}
	}()
}

// Wait waits for all workers to complete.
func (p *simpleWorkerPool[E, R]) Wait() {
	// Wait for all workers to complete.
	p.workersWaitGroup.Wait()
}
