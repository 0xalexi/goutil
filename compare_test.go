package goutil

import (
	"fmt"
	"testing"
)

func shouldEq(t *testing.T, a, b interface{}) {
	if !CompareInterfaces(a, b) {
		t.Error(fmt.Sprintf("should be eq failed for %v : %v", a, b))
	}
}

func shouldNeq(t *testing.T, a, b interface{}) {
	if CompareInterfaces(a, b) {
		t.Error(fmt.Sprintf("should be neq failed for %v : %v", a, b))
	}
}

func TestCompare(t *testing.T) {
	shouldEq(t, []string{"S", "bottoms"}, []string{"bottoms", "S"})
	shouldEq(t, true, true)
	shouldEq(t, 0.43, 0.43)
	shouldEq(t, "9CVqIR74nk", "9CVqIR74nk")

	shouldNeq(t, "9CVqIR74nk", "9CVqIR74n")
	shouldNeq(t, []string{"S", "bottoms"}, []string{"bottoms"})
	shouldNeq(t, []string{"S", "bottoms"}, []int{1})
}
