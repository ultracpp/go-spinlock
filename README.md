# SpinLock Benchmark

## Overview

This project demonstrates the implementation and benchmarking of a simple spin lock in Go, compared against Go's built-in `sync.Mutex`. Spin locks are useful in scenarios where you want to avoid the overhead of operating system context switches that traditional mutexes may incur.

## SpinLock Implementation

The `SpinLock` implemented here is a basic spin lock using atomic operations. It spins (busy-waits) until the lock becomes available, making it suitable for scenarios where the critical section is expected to be held for a short duration.

## Implementation Details

- **Locking**: Uses `atomic.CompareAndSwapInt32` to attempt to acquire the lock.
- **Unlocking**: Uses `atomic.StoreInt32` to release the lock.

## Features

- **Basic SpinLock**: A simple spinlock using an atomic boolean to manage the lock state.
- **Backoff Strategy**: Incorporates a backoff strategy that includes yielding and optional sleeping to reduce CPU usage during contention.
- **Max attempts Lock**: Adds a timeout feature to the lock acquisition, returning an error if the lock cannot be obtained after a specified number of attempts.

## Benchmark Environment

- **CPU**: AMD Ryzen 7 6800H
- **Memory**: 32GB

## Benchmark Results

### SpinLock vs. Sync.Mutex

#### Test Scenario
- **Threads**: 32 concurrent goroutines
- **Jobs per Thread**: 1,000,000

#### Results
- **SpinLock**: 32,000,000 operations in 175.0841ms
- **Sync.Mutex**: 32,000,000 operations in 1.629234s

### Performance Comparison
- The spin lock outperformed `sync.Mutex` significantly, approximately 10 times faster in this specific test scenario.
- Spin locks are particularly advantageous in scenarios with short critical sections, where the overhead of operating system-level mutexes can be relatively high.

## Test Code

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	ThreadCount = 32
	JobCount    = 1000000
)

var sum int64 = 0

func testSpinLock() {
	fmt.Println("========test_spin_lock========")

	var wg sync.WaitGroup
	sum = 0

	{
		spinLock := NewSpinLock()
		start := time.Now()

		for i := 0; i < ThreadCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for j := 0; j < JobCount; j++ {
					spinLock.Lock()
					sum += 1
					spinLock.Unlock()
				}
			}()
		}

		wg.Wait()

		elapsed := time.Since(start)
		fmt.Printf("SpinLock: %d %v\n", sum, elapsed)
	}

	sum = 0

	{
		mutex := &sync.Mutex{}
		start := time.Now()

		for i := 0; i < ThreadCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for j := 0; j < JobCount; j++ {
					mutex.Lock()
					sum += 1
					mutex.Unlock()
				}
			}()
		}

		wg.Wait()

		elapsed := time.Since(start)
		fmt.Printf("Sync.Mutex: %d %v\n", sum, elapsed)
	}
}

func main() {
	testSpinLock()
}
