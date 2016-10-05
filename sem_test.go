package sem

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func fetchSomething(requests *int32) error {
	atomic.AddInt32(requests, 1)
	time.Sleep(100 * time.Microsecond)
	return nil
}

func Example_exec() {
	s := New(2)
	var wg sync.WaitGroup
	var requests int32

	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < 10; k++ {
				s.Exec(func() {
					fetchSomething(&requests)
				})
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Finished %v requests.", atomic.LoadInt32(&requests))

	// Output:
	// Finished 300 requests.
}

func Example_acquireRelease() {
	s := New(2)
	var wg sync.WaitGroup
	var requests int32

	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < 10; k++ {
				s.Acquire()

				err := fetchSomething(&requests)
				if err != nil {
					fmt.Printf("ERROR: request failed: %v\n", err)
				}

				s.Release()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Finished %v requests.", atomic.LoadInt32(&requests))

	// Output:
	// Finished 300 requests.
}

func Example_exece() {
	s := New(2)
	var wg sync.WaitGroup
	var requests int32

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < 3; k++ {
				err := s.Exece(func() error {
					return fetchSomething(&requests)
				})
				if err != nil {
					fmt.Printf("ERROR: request failed: %v\n", err)
				}
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Finished %v requests.", atomic.LoadInt32(&requests))

	// Output:
	// Finished 9 requests.
}

func TestSem(t *testing.T) {
	s := New(7)
	var cur int32
	var max int32
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < 1000; k++ {
				s.Acquire()
				v := atomic.AddInt32(&cur, 1)

				old := atomic.LoadInt32(&max)
				for v > old {
					atomic.CompareAndSwapInt32(&max, old, v)
					old = atomic.LoadInt32(&max)
				}

				time.Sleep(1 * time.Microsecond)

				atomic.AddInt32(&cur, -1)
				s.Release()
			}
		}()
	}

	wg.Wait()

	maxv := int(atomic.LoadInt32(&max))
	if maxv != s.Cap() {
		t.Fatalf("Incorrect number of concurrent executors, want %v, actual %v", s.Cap(), maxv)
	}
}
