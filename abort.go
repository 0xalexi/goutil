package goutil

import (
	"fmt"
	"github.com/pkg/errors"
)

type AbortHandler struct {
	steps []func() error
}

func NewAbortHandler() *AbortHandler {
	steps := make([]func() error, 0)
	return &AbortHandler{steps}
}

func (a *AbortHandler) Run() (err error) {
	if a == nil {
		return
	}
	deleted := 0
	for i := range a.steps {
		j := i - deleted
		if j < len(a.steps) {
			f := a.steps[j]
			_err := f()
			if _err != nil {
				errors.Wrap(err, fmt.Sprint(_err))
			}
			a.steps = a.steps[:j+copy(a.steps[j:], a.steps[j+1:])]
			deleted++
		}
	}
	return
}

func (a *AbortHandler) Append(f func() error) {
	if a == nil {
		return
	}
	a.steps = append(a.steps, f)
}

func (a *AbortHandler) Push(f func() error) {
	if a == nil {
		return
	}
	a.steps = append([]func() error{f}, a.steps...)
}
