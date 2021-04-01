package goutil

import "time"

func MillisecondsToTime(n int64) time.Time {
	return time.Unix(0, n*int64(time.Millisecond))
}
