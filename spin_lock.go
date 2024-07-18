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
	"runtime"
	"sync/atomic"
	"time"
)

type SpinLock struct {
	lock int32
}

const (
	spinLockSleepOneFrequency = 50
	spinLockMaxAttempts       = 200
)

func NewSpinLock() *SpinLock {
	return &SpinLock{}
}

func (s *SpinLock) Lock() {
	freq := 0

	for !s.TryLock() {
		runtime.Gosched()

		freq++

		if freq == spinLockSleepOneFrequency {
			time.Sleep(1 * time.Millisecond)
			freq = 0
		}
	}
}

func (s *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapInt32(&s.lock, 0, 1)
}

func (s *SpinLock) Unlock() {
	atomic.StoreInt32(&s.lock, 0)
}

func SpinWait(condition func() bool) {
	freq := 0

	for condition() {
		runtime.Gosched()

		freq++

		if freq == spinLockSleepOneFrequency {
			time.Sleep(1 * time.Millisecond)
			freq = 0
		}
	}
}

func (s *SpinLock) LockWithMaxAttempts() bool {
	freq := 0

	for !s.TryLock() {
		runtime.Gosched()

		freq++

		if freq == spinLockSleepOneFrequency {
			if freq == spinLockMaxAttempts {
				return false
			}

			time.Sleep(1 * time.Millisecond)
			freq = 0
		}
	}

	return true
}
