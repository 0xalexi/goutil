package goutil

import (
	"fmt"
	"runtime"

	sigar "github.com/cloudfoundry/gosigar"
)

type MemStats struct {
	m      runtime.MemStats
	sigMem *sigar.Mem
}

func NewMemStats() *MemStats {
	return &MemStats{sigMem: new(sigar.Mem)}
}

func (s *MemStats) Read() {
	runtime.ReadMemStats(&s.m)
	s.sigMem.Get()
}

func (s *MemStats) String() string {
	v := fmt.Sprint("mem-acquired:", s.m.Sys, "alloc:", s.m.Alloc, "total-alloc:", s.m.TotalAlloc)
	if s.sigMem != nil {
		v += fmt.Sprint(" sig-total:", s.sigMem.Total, "sig-free:", s.sigMem.Free, "sig-used:", float64(s.sigMem.Used), "sig-percent-used:", fmt.Sprintf("%f%%", float64(s.sigMem.Used)/float64(s.sigMem.Total)*100))
	}
	return v
}
