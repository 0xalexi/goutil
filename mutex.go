package goutil

import (
	"sync"
)

type MultiMutex struct {
	Locks      map[string]*sync.RWMutex
	LocksMutex sync.RWMutex
}

func NewMultiMutex() *MultiMutex {
	return &MultiMutex{
		Locks: make(map[string]*sync.RWMutex),
	}
}

func (m *MultiMutex) GetLock(k string) *sync.RWMutex {
	// attempt rlocking read
	m.LocksMutex.RLock()
	lock, ok := m.Locks[k]
	if ok {
		m.LocksMutex.RUnlock()
		return lock
	}
	m.LocksMutex.RUnlock()

	// write new lock if necessary
	m.LocksMutex.Lock()
	defer m.LocksMutex.Unlock()
	lock, ok = m.Locks[k]
	if !ok || lock == nil {
		lock = &sync.RWMutex{}
		m.Locks[k] = lock
	}
	return lock
}

func (m *MultiMutex) Lock(k string) {
	lock := m.GetLock(k)
	lock.Lock()
}

func (m *MultiMutex) Unlock(k string) {
	lock := m.GetLock(k)
	lock.Unlock()
}

func (m *MultiMutex) RLock(k string) {
	lock := m.GetLock(k)
	lock.RLock()
}

func (m *MultiMutex) RUnlock(k string) {
	lock := m.GetLock(k)
	lock.RUnlock()
}

type MultiMutexUint32 struct {
	Locks      map[uint32]*sync.Mutex
	LocksMutex sync.RWMutex
}

func NewMultiMutexUint32() *MultiMutexUint32 {
	return &MultiMutexUint32{
		Locks: make(map[uint32]*sync.Mutex),
	}
}

func (m *MultiMutexUint32) GetLock(k uint32) *sync.Mutex {
	m.LocksMutex.Lock()
	defer m.LocksMutex.Unlock()
	lock := m.Locks[k]
	if lock == nil {
		lock = &sync.Mutex{}
		m.Locks[k] = lock
	}
	return lock
}

func (m *MultiMutexUint32) Lock(k uint32) {
	lock := m.GetLock(k)
	lock.Lock()
}

func (m *MultiMutexUint32) Unlock(k uint32) {
	lock := m.GetLock(k)
	lock.Unlock()
}

type MultiMutexUint64 struct {
	Locks      map[uint64]*sync.Mutex
	LocksMutex sync.RWMutex
}

func NewMultiMutexUint64() *MultiMutexUint64 {
	return &MultiMutexUint64{
		Locks: make(map[uint64]*sync.Mutex),
	}
}

func (m *MultiMutexUint64) GetLock(k uint64) *sync.Mutex {
	m.LocksMutex.Lock()
	defer m.LocksMutex.Unlock()
	lock := m.Locks[k]
	if lock == nil {
		lock = &sync.Mutex{}
		m.Locks[k] = lock
	}
	return lock
}

func (m *MultiMutexUint64) Lock(k uint64) {
	lock := m.GetLock(k)
	lock.Lock()
}

func (m *MultiMutexUint64) Unlock(k uint64) {
	lock := m.GetLock(k)
	lock.Unlock()
}

type OptionalMode bool

const (
	OPT_ENABLED  OptionalMode = true
	OPT_DISABLED OptionalMode = false
)

type OptionalWaitGroup struct {
	sync.WaitGroup
}

func NewOptionalWaitGroup(enabled OptionalMode) *OptionalWaitGroup {
	if !enabled {
		return nil
	}
	return new(OptionalWaitGroup)
}

func (wg *OptionalWaitGroup) Add(n int) {
	if wg == nil {
		return
	}
	wg.WaitGroup.Add(n)
}

func (wg *OptionalWaitGroup) Done() {
	if wg == nil {
		return
	}
	wg.WaitGroup.Done()
}

func (wg *OptionalWaitGroup) Wait() {
	if wg == nil {
		return
	}
	wg.WaitGroup.Wait()
}

type OptionalRWMutex struct {
	sync.RWMutex
}

func NewOptionalRWMutex(enabled OptionalMode) *OptionalRWMutex {
	if !enabled {
		return nil
	}
	return new(OptionalRWMutex)
}

func NewDisabledOptionalRWMutex() *OptionalRWMutex {
	return nil
}

func (lock *OptionalRWMutex) Lock() {
	if lock == nil {
		return
	}
	lock.RWMutex.Lock()
}

func (lock *OptionalRWMutex) Unlock() {
	if lock == nil {
		return
	}
	lock.RWMutex.Unlock()
}

func (lock *OptionalRWMutex) RLock() {
	if lock == nil {
		return
	}
	lock.RWMutex.RLock()
}

func (lock *OptionalRWMutex) RUnlock() {
	if lock == nil {
		return
	}
	lock.RWMutex.RUnlock()
}
