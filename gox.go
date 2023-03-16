package gox

import "sync"

// goX struct to handle concurrent execution of functions
type goX struct {
	waitGroup sync.WaitGroup // WaitGroup to wait for all functions to complete
}

// goE struct to handle concurrent execution of functions with error handling
type goE struct {
	waitGroup sync.WaitGroup   // WaitGroup to wait for all functions to complete
	errMap    map[string]error // Map to hold errors returned by functions
	mutex     sync.Mutex       // Mutex to synchronize access to errMap
}

// Wait method for goX struct to wait for all functions to complete
func (gx *goX) Wait() {
	gx.waitGroup.Wait()
}

// WaitE method for goE struct to wait for all functions to complete with error handling
func (ge *goE) WaitE() map[string]error {
	ge.waitGroup.Wait()
	if len(ge.errMap) == 0 {
		return nil
	}
	return ge.errMap
}

// Run method for goX struct to start concurrent execution of functions
func (gx *goX) Run(funcs ...func()) *goX {
	for _, f := range funcs {
		gx.waitGroup.Add(1) // Add function to WaitGroup
		go func(f func()) {
			defer gx.waitGroup.Done() // Signal completion of function
			f()                       // Execute function
		}(f) // Pass function as argument to goroutine
	}
	return gx
}

// RunE method for goE struct to start concurrent execution of functions with error handling
func (ge *goE) RunE(name string, f func() error) *goE {
	ge.waitGroup.Add(1) // Add function to WaitGroup
	go func() {
		defer ge.waitGroup.Done()   // Signal completion of function
		if err := f(); err != nil { // Execute function and check for errors
			ge.mutex.Lock()       // Acquire mutex to access errMap
			ge.errMap[name] = err // Add error to errMap with function name as key
			ge.mutex.Unlock()     // Release mutex
		}
	}()

	return ge
}

// Start function to create new goX struct
func Start() *goX {
	return &goX{
		waitGroup: sync.WaitGroup{},
	}
}

// StartE function to create new goE struct
func StartE() *goE {
	return &goE{
		waitGroup: sync.WaitGroup{},
		errMap:    make(map[string]error),
		mutex:     sync.Mutex{},
	}
}

// Run function to start concurrent execution of functions using a new goX struct
func Run(funcs ...func()) *goX {
	return Start().Run(funcs...)
}

// RunE function to start concurrent execution of functions with error handling using a new goE struct
func RunE(name string, f func() error) *goE {
	return StartE().RunE(name, f)
}
