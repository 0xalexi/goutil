package goutil

import (
	"sync"
)

type StringSet struct {
	m map[string]bool
	sync.RWMutex
}

func NewStringSet() *StringSet {
	return &StringSet{
		m: make(map[string]bool),
	}
}

// Add add
func (s *StringSet) Add(item string) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

// Remove deletes the specified item from the map
func (s *StringSet) Remove(item string) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, item)
}

// Has looks for the existence of an item
func (s *StringSet) Has(item string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

// Len returns the number of items in a set.
func (s *StringSet) Len() int {
	return len(s.List())
}

// Clear removes all items from the set
func (s *StringSet) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[string]bool)
}

// IsEmpty checks for emptiness
func (s *StringSet) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

// StringSet returns a slice of all items
func (s *StringSet) List() []string {
	s.RLock()
	defer s.RUnlock()
	list := make([]string, 0)
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

type Uint64Set struct {
	m map[uint64]bool
	sync.RWMutex
}

func NewUint64Set() *Uint64Set {
	return &Uint64Set{
		m: make(map[uint64]bool),
	}
}

// Add add
func (s *Uint64Set) Add(item uint64) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

// Remove deletes the specified item from the map
func (s *Uint64Set) Remove(item uint64) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, item)
}

// Has looks for the existence of an item
func (s *Uint64Set) Has(item uint64) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

// Len returns the number of items in a set.
func (s *Uint64Set) Len() int {
	return len(s.List())
}

// Clear removes all items from the set
func (s *Uint64Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[uint64]bool)
}

// IsEmpty checks for emptiness
func (s *Uint64Set) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

// Uint64Set returns a slice of all items
func (s *Uint64Set) List() Uint64s {
	s.RLock()
	defer s.RUnlock()
	list := make(Uint64s, 0)
	for item := range s.m {
		list = append(list, item)
	}
	return list
}
