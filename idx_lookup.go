package goutil

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type IdxLookup struct {
	IdIdx map[string]int `json:"id_to_idx"`
	IdxId []string
	mu    sync.RWMutex
}

func NewIdxLookup() *IdxLookup {
	return &IdxLookup{
		IdIdx: map[string]int{},
		IdxId: []string{},
	}
}

func (v *IdxLookup) MarshalBinary() ([]byte, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return GetBytes(v.IdIdx)
}

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (v *IdxLookup) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	v.IdIdx = map[string]int{}
	err := DecodeBytes(data, &v.IdIdx)
	if err != nil {
		return err
	}
	v.IdxId = make([]string, len(v.IdIdx))
	for id, idx := range v.IdIdx {
		v.IdxId[idx] = id
	}
	return nil
}

func LoadIdxLookup(filepath string) (lookup IdxLookup, err error) {
	var b []byte
	if b, err = ioutil.ReadFile(filepath); err != nil {
		return
	}
	lookup.IdIdx = make(map[string]int)
	if err = json.Unmarshal(b, &lookup.IdIdx); err != nil {
		return
	}
	lookup.IdxId = make([]string, len(lookup.IdIdx))
	for id, idx := range lookup.IdIdx {
		lookup.IdxId[idx] = id
	}
	return
}

func (l *IdxLookup) GetId(idx int) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if idx >= len(l.IdxId) || idx < 0 {
		return "", ErrOutOfRange
	}
	return l.IdxId[idx], nil
}

func (l *IdxLookup) CheckId(id string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	_, ok := l.IdIdx[id]
	return ok
}

func (l *IdxLookup) GetIdx(id string) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	idx, ok := l.IdIdx[id]
	if !ok {
		l.insert(id)
		return l.IdIdx[id], nil
	}
	return idx, nil
}

func (l *IdxLookup) Num() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.IdIdx)
}

func (l *IdxLookup) insert(id string) {
	l.IdIdx[id] = len(l.IdxId)
	l.IdxId = append(l.IdxId, id)
}

type StaticIdxLookup struct {
	IdIdx map[string]int `json:"name_idx"`
	IdxId []string
}

func LoadStaticIdxLookup(filepath string) (lookup StaticIdxLookup, err error) {
	var b []byte
	if b, err = ioutil.ReadFile(filepath); err != nil {
		return
	}
	lookup.IdIdx = make(map[string]int)
	if err = json.Unmarshal(b, &lookup.IdIdx); err != nil {
		return
	}
	lookup.IdxId = make([]string, len(lookup.IdIdx))
	for id, idx := range lookup.IdIdx {
		lookup.IdxId[idx] = id
	}
	return
}

func (l *StaticIdxLookup) GetId(idx int) (string, error) {
	if idx >= len(l.IdxId) || idx < 0 {
		return "", ErrOutOfRange
	}
	return l.IdxId[idx], nil
}

func (l *StaticIdxLookup) GetIdx(id string) (int, error) {
	idx, ok := l.IdIdx[id]
	if !ok {
		return idx, ErrNotFound
	}
	return idx, nil
}

func (l *StaticIdxLookup) Num() int {
	return len(l.IdIdx)
}
