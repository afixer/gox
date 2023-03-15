package gox_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/afixer/gox"
)

// TestRun tests the Run function of the gox package.
func TestRun(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	gox.Run(
		func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
		},
		func() {
			defer wg.Done()
			time.Sleep(200 * time.Millisecond)
		},
	)

	// Wait for the completion of all functions
	wg.Wait()
}

// TestRunE tests the RunE function of the gox package.
func TestRunE(t *testing.T) {
	// Test with two functions that return errors
	errMap := gox.StartE().
		RunE("foo", func() error {
			return errors.New("foo error")
		}).
		RunE("bar", func() error {
			return errors.New("bar error")
		}).
		WaitE()

	// Check the expected errors in the error map
	expected := map[string]error{
		"foo": errors.New("foo error"),
		"bar": errors.New("bar error"),
	}
	if len(errMap) != len(expected) {
		t.Fatalf("Expected error map size %d, but got %d", len(expected), len(errMap))
	}
	for k, v := range expected {
		if errMap[k].Error() != v.Error() {
			t.Fatalf("Expected error %q for function %q, but got %q", v.Error(), k, errMap[k].Error())
		}
	}
}

// TestRunRace tests that Run is race-safe.
func TestRunRace(t *testing.T) {
	const n = 1000
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		gox.Run(func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
		})
	}

	// Wait for the completion of all functions
	wg.Wait()
}

// TestRunEWaitRace tests that RunE and WaitE are race-safe.
func TestRunEWaitRace(t *testing.T) {
	const n = 1000
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		gox.StartE().RunE("foo", func() error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}).WaitE()
		wg.Done()
	}

	// Wait for the completion of all functions
	wg.Wait()
}
