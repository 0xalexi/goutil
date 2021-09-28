package goutil

import "math"

func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
		if s[0] == '\'' && s[len(s)-1] == '\'' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

type Strings []string

func (s Strings) Contains(v string) bool {
	if len(s) == 0 {
		return false
	}
	for _, _v := range s {
		if _v == v {
			return true
		}
	}
	return false
}

func (l Strings) Map() map[string]bool {
	m := make(map[string]bool)
	for _, v := range l {
		m[v] = true
	}
	return m
}

func (l *Strings) RemoveDuplicates() Strings {
	if l == nil {
		return nil
	}
	if *l == nil {
		return *l
	}
	found := make(map[string]bool)
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

func (l Strings) GetBatches(batchsize int) []Strings {
	batches := int(math.Ceil(float64(len(l)) / float64(batchsize)))
	out := make([]Strings, batches)
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
