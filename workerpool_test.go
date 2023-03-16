package gox

import (
	"testing"
)

func TestWorkerPool(t *testing.T) {
	var (
		numWorkers = 12
		elements   = []int{1, 2, 3, 4, 5, 6}
		processor  = func(e int) (int, error) {
			return e * 2, nil
		}
		sum     = 0
		reducer = func(r int) {
			sum += r
		}
		expectedSum = 42
	)

	pool := NewWorkerPool(numWorkers, processor, reducer)
	for _, e := range elements {
		pool.Add(e)
	}
	pool.Wait()
	if sum != expectedSum {
		t.Errorf("expected sum to be %d, got %d", expectedSum, sum)
	}
}
