package goutil

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v3/mem"
)

type MemStats struct {
	m     runtime.MemStats
	psmem *mem.VirtualMemoryStat
}

func NewMemStats() *MemStats {
	return &MemStats{}
}

func (s *MemStats) Read() {
	runtime.ReadMemStats(&s.m)
	s.psmem, _ = mem.VirtualMemory()
}

func (s *MemStats) String() string {
	v := fmt.Sprint("mem-acquired:", s.m.Sys, "alloc:", s.m.Alloc, "total-alloc:", s.m.TotalAlloc)
	if s.psmem != nil {
		v += fmt.Sprint(" sig-total:", s.psmem.Total, "sig-free:", s.psmem.Free, "sig-used:", float64(s.psmem.Used), "sig-percent-used:", s.psmem.UsedPercent)
	}
	return v
}
