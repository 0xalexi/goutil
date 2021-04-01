package goutil

import "math"

func StringPointer(v string) *string {
	return &v
}

type StringGetter interface {
	String() string
}

type Uint64s []uint64

func (l Uint64s) Max() uint64 {
	var max uint64
	for _, id := range l {
		if id > max {
			max = id
		}
	}
	return max
}

func (l Uint64s) GetBatches(batchsize int) []Uint64s {
	batches := int(math.Ceil(float64(len(l)) / float64(batchsize)))
	out := make([]Uint64s, batches)
	idx := 0
	bs := batchsize
	for i := 0; ; i += batchsize {
		if idx == batches {
			break
		}
		bs = batchsize
		if len(l)-i < batchsize {
			bs = len(l) - i
		}
		out[idx] = l[i : i+bs]
		idx++
	}
	return out
}

func (l Uint64s) Map() Uint64BoolMap {
	m := make(Uint64BoolMap)
	for _, v := range l {
		m[v] = true
	}
	return m
}

func (l Uint64s) Contains(v ...uint64) bool {
	m := l.Map()
	for _, _v := range v {
		if !m[_v] {
			return false
		}
	}
	return true
}

func (l *Uint64s) RemoveDuplicates() Uint64s {
	if l == nil {
		return nil
	}
	if *l == nil {
		return *l
	}
	found := make(map[uint64]bool)
	j := 0
	for i, x := range *l {
		if !found[x] {
			found[x] = true
			(*l)[j] = (*l)[i]
			j++
		}
	}
	*l = (*l)[:j]
	return *l
}

type Uint64BoolMap map[uint64]bool

func (m Uint64BoolMap) Keys() Uint64s {
	out := make(Uint64s, len(m))
	i := 0
	for k := range m {
		out[i] = k
		i++
	}
	return out
}

type StringBoolMap map[string]bool

func (m StringBoolMap) Keys() []string {
	out := make([]string, len(m))
	i := 0
	for k := range m {
		out[i] = k
		i++
	}
	return out
}
