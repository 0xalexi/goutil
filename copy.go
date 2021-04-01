package goutil

import "github.com/mohae/deepcopy"

func CopyMapStringInterface(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return nil
	}
	t := make(map[string]interface{}, len(m))
	for k, v := range m {
		t[k] = v
	}
	return t
}

func CopyStringSlice(s []string) []string {
	if s == nil {
		return nil
	}
	t := make([]string, len(s))
	copy(t, s)
	return t
}

func CopyByteSlice(s []byte) []byte {
	if s == nil {
		return nil
	}
	t := make([]byte, len(s))
	copy(t, s)
	return t
}

func CopyUint32Slice(s []uint32) []uint32 {
	if s == nil {
		return nil
	}
	t := make([]uint32, len(s))
	copy(t, s)
	return t
}

func CopyBoolSlice(s []bool) []bool {
	if s == nil {
		return nil
	}
	t := make([]bool, len(s))
	copy(t, s)
	return t
}

func DeepCopy(v interface{}) interface{} {
	return deepcopy.Copy(v)
}
