package gox

import "sync"

// goX struct to handle concurrent execution of functions
type goX struct {
	w sync.WaitGroup // WaitGroup to wait for all functions to complete
}

// goE struct to handle concurrent execution of functions with error handling
type goE struct {
	w      sync.WaitGroup   // WaitGroup to wait for all functions to complete
	errMap map[string]error // Map to hold errors returned by functions
	l      sync.Mutex       // Mutex to synchronize access to errMap
}

// Wait method for goX struct to wait for all functions to complete
func (w *goX) Wait() {
	w.w.Wait()
}

// WaitE method for goE struct to wait for all functions to complete with error handling
func (w *goE) WaitE() map[string]error {
	w.w.Wait()
	if len(w.errMap) == 0 {
		return nil
	}
	return w.errMap
}

// Run method for goX struct to start concurrent execution of functions
func (w *goX) Run(fs ...func()) *goX {
	for _, f := range fs {
		w.w.Add(1) // Add function to WaitGroup
		go func(f func()) {
			defer w.w.Done() // Signal completion of function
			f()              // Execute function
		}(f) // Pass function as argument to goroutine
	}
	return w
}

// RunE method for goE struct to start concurrent execution of functions with error handling
func (w *goE) RunE(name string, f func() error) *goE {
	w.w.Add(1) // Add function to WaitGroup
	go func() {
		defer w.w.Done()            // Signal completion of function
		if err := f(); err != nil { // Execute function and check for errors
			w.l.Lock()           // Acquire mutex to access errMap
			w.errMap[name] = err // Add error to errMap with function name as key
			w.l.Unlock()         // Release mutex
		}
	}()

	return w
}

// Start function to create new goX struct
func Start() *goX {
	return &goX{
		w: sync.WaitGroup{},
	}
}

// StartE function to create new goE struct
func StartE() *goE {
	return &goE{
		w:      sync.WaitGroup{},
		errMap: make(map[string]error),
		l:      sync.Mutex{},
	}
}

// Run function to start concurrent execution of functions using a new goX struct
func Run(fs ...func()) *goX {
	return Start().Run(fs...)
}

// RunE function to start concurrent execution of functions with error handling using a new goE struct
func RunE(name string, f func() error) *goE {
	return StartE().RunE(name, f)
}
