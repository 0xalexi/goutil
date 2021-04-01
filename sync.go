package goutil

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/sync/semaphore"
)

func AtomicAddFloat64(val *float64, delta float64) (new float64) {
	for {
		old := *val
		new = old + delta
		if atomic.CompareAndSwapUint64((*uint64)(unsafe.Pointer(val)), math.Float64bits(old), math.Float64bits(new)) {
			break
		}
	}
	return
}

type TimeoutMutex struct {
	l       *semaphore.Weighted
	timeout time.Duration
}

func NewTimeoutMutex(timeout time.Duration) TimeoutMutex {
	return TimeoutMutex{
		l:       semaphore.NewWeighted(1),
		timeout: timeout,
	}
}

func (m TimeoutMutex) Lock() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	return m.l.Acquire(ctx, 1)
}

func (m TimeoutMutex) Unlock() {
	m.l.Release(1)
}

type ConcurrentIdIdxMap struct {
	Data map[uint64]int
	Mu   sync.RWMutex
}

func NewConcurrentIdIdxMap() *ConcurrentIdIdxMap {
	return &ConcurrentIdIdxMap{Data: make(map[uint64]int)}
}

func NewConcurrentIdIdxMapWithData(m map[uint64]int) *ConcurrentIdIdxMap {
	return &ConcurrentIdIdxMap{Data: m}
}

func (m *ConcurrentIdIdxMap) Len() int {
	if m == nil {
		return 0
	}
	return len(m.Data)
}

func (m ConcurrentIdIdxMap) ReadOk(id uint64) (int, bool) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	v, ok := m.Data[id]
	return v, ok
}

func (m ConcurrentIdIdxMap) Read(id uint64) int {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	return m.Data[id]
}

func (m *ConcurrentIdIdxMap) Write(k uint64, v int) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.Data[k] = v
}

func (m *ConcurrentIdIdxMap) SetData(data map[uint64]int) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.Data = data
}

// type Locker interface {
// 	Lock()
// }

// type RLocker interface {
// 	Lock()
// }

// func LockWithTimeout(m Locker, timeout time.Duration) error {

// }

// func RLockWithTimeout(m RLocker, timeout time.Duration) {

// }
