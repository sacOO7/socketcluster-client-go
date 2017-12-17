package utils

import (
	"sync/atomic"
)

type AtomicCounter struct {
	Counter uint64
}

func (atomicCounter *AtomicCounter) IncrementAndGet() uint64 {
	return atomic.AddUint64(&atomicCounter.Counter, 1)
}

func (atomicCounter *AtomicCounter) GetAndIncrement() uint64 {
	defer atomic.AddUint64(&atomicCounter.Counter, 1)
	return atomic.LoadUint64(&atomicCounter.Counter)
}

func (atomicCounter *AtomicCounter) Reset() {
	atomic.StoreUint64(&atomicCounter.Counter, 0)
}

func (atomicCounter *AtomicCounter) Value() uint64 {
	return atomic.LoadUint64(&atomicCounter.Counter)
}
