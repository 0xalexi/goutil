package goutil

import (
	"errors"
	"math"
	"time"
)

var NO_TIMEOUT = time.Duration(math.MaxInt64)

func WaitFor(f func() bool, period, timeout time.Duration) error {
	if period < 1 {
		period = time.Duration(2 * time.Second)
	}
	if timeout < 1 {
		timeout = time.Duration(60 * time.Second)
	}
	t := time.Duration(0)
	for t < timeout || timeout == NO_TIMEOUT {
		s := time.Now()
		if f() {
			return nil
		}
		d := time.Now().Sub(s)
		if d < period {
			time.Sleep(period - d)
			t += period
		} else {
			t += d
		}
	}
	return errors.New("Timed out after " + timeout.String())
}
