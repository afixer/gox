package gox

import "sync"

// The workerPool type is a concurrent worker pool,
// that processes a list of elements E and returns a result of type R.
type workerPool[E any, R any] struct {
	inputChannel     chan struct{}      // Channel to limit the number of concurrent workers.
	outputChannel    chan R             // Channel to receive results from workers.
	workerFunc       func(E) (R, error) // Function that processes an element and returns a result and an error.
	workersWaitGroup sync.WaitGroup     // WaitGroup to wait for all workers to complete.
	resultWaitGroup  sync.WaitGroup     // WaitGroup to wait for the result consumer to complete.
}

// NewWorkerPool creates a new worker pool with nWorkers concurrent workers.
// The workerFunc function is called for each element in the pool.
// The resultConsumer function is called for each result returned by the workerFunc.
func NewWorkerPool[E, R any](nWorkers int, workerFunc func(E) (R, error), resultConsumer func(R)) *workerPool[E, R] {
	// Create a new workerPool struct.
	p := &workerPool[E, R]{
		inputChannel:     make(chan struct{}, nWorkers), // Buffered channel to limit the number of concurrent workers.
		outputChannel:    make(chan R),                  // Channel to receive results from workers.
		workerFunc:       workerFunc,                    // Function to process an element.
		workersWaitGroup: sync.WaitGroup{},              // WaitGroup to wait for all workers to complete.
		resultWaitGroup:  sync.WaitGroup{},              // WaitGroup to wait for the result consumer to complete.
	}

	// Launch a goroutine to consume the results returned by the workers.
	p.resultWaitGroup.Add(1)
	go func() {
		defer p.resultWaitGroup.Done()
		for r := range p.outputChannel {
			resultConsumer(r)
		}
	}()

	return p
}

// Add adds a new element e to the pool.
func (p *workerPool[E, R]) Add(e E) {
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
		result, err := p.workerFunc(e)
		if err != nil {
			return
		}

		// Send the result to the outputChannel.
		p.outputChannel <- result
	}()
}

// Wait waits for all workers to complete.
func (p *workerPool[E, R]) Wait() {
	// Wait for all workers to complete.
	p.workersWaitGroup.Wait()

	// Close the outputChannel and wait for the result consumer to complete.
	close(p.outputChannel)
	p.resultWaitGroup.Wait()
}
