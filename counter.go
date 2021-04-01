package goutil

import "sync"

type Uint64Counter struct {
	m map[uint64]int
	sync.RWMutex
}

func NewUint64Counter() *Uint64Counter {
	return &Uint64Counter{
		m: make(map[uint64]int),
	}
}

func (s *Uint64Counter) Increment(k uint64) {
	s.Lock()
	defer s.Unlock()
	s.m[k] += 1
}

func (s *Uint64Counter) Decrement(k uint64) {
	s.Lock()
	defer s.Unlock()
	if s.m[k] > 0 {
		s.m[k] -= 1
	}
}

func (s *Uint64Counter) Remove(k uint64) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, k)
}

func (s *Uint64Counter) Get(k uint64) int {
	s.RLock()
	defer s.RUnlock()
	return s.m[k]
}
