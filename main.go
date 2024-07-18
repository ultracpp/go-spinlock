/*
 * SpinLock - custom implementation of a spinlock in Go
 * Copyright (c) 2024 Eungsuk Jeon
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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
		mutex := &Locker{}
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
		fmt.Printf("Locker: %d %v\n", sum, elapsed)
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
