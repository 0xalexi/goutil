package goutil

import (
	"math/rand"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano())

const (
	charBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charIdxBits = 6                  // 6 bits to represent a char index
	charIdxMask = 1<<charIdxBits - 1 // All 1-bits, as many as charIdxBits
	charIdxMax  = 63 / charIdxBits   // # of char indices fitting in 63 bits
)

func FastRandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for charIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), charIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), charIdxMax
		}
		if idx := int(cache & charIdxMask); idx < len(charBytes) {
			b[i] = charBytes[idx]
			i--
		}
		cache >>= charIdxBits
		remain--
	}

	return string(b)
}
