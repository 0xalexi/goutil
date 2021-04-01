package goutil

import (
	"reflect"
	"sort"
)

func CompareInterfaces(a, b interface{}) bool {
	a = sortVal(a)
	b = sortVal(b)
	return reflect.DeepEqual(a, b)
}

func sortVal(value interface{}) interface{} {
	if value == nil {
		return value
	}
	switch t := value.(type) {
	case string, float64, bool:
		return t
	case []string:
		sort.Strings(t)
		return t
	case []float64:
		sort.Float64s(t)
		return t
	case []bool:
		return t
	case []interface{}:
		// This may be an issue
		return t
	}
	return value
}
