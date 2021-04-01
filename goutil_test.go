package goutil

import "testing"

func matchTypes(t *testing.T) {
	type a struct{}
	type b struct{}

	if TypesMatch(a{}, b{}) {
		t.Error("match true for a, b")
	}

	if !TypesMatch(a{}, a{}) {
		t.Error("match false for a, a")
	}
	if !TypesMatch(a{}, &a{}) {
		t.Error("match false for a, *a")
	}
	if !TypesMatch(&a{}, &a{}) {
		t.Error("match false for *a, *a")
	}
	if !TypesMatch(&a{}, a{}) {
		t.Error("match false for *a, a")
	}
}

func TestReflect(t *testing.T) {
	t.Run("MatchTypes", matchTypes)
}
