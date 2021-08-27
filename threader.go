package goutil

import (
	"sync"
	"sync/atomic"
)

type Threader struct {
	RunningCond sync.Cond
	Running     int64
	NumThreads  int64
}

func NewThreader(numThreads int) *Threader {
	return &Threader{
		RunningCond: sync.Cond{L: &sync.Mutex{}},
		Running:     0,
		NumThreads:  int64(numThreads),
	}
}

func (t *Threader) Acquire() {
	t.RunningCond.L.Lock()
	for t.Running == atomic.LoadInt64(&t.NumThreads) {
		t.RunningCond.Wait()
	}
	t.Running++
	t.RunningCond.L.Unlock()
}

func (t *Threader) SetThreadCount(numThreads int) {
	atomic.StoreInt64(&t.NumThreads, int64(numThreads))
}

func (t *Threader) GetThreadCount() int {
	return int(t.NumThreads)
}

func (t *Threader) Release() {
	t.RunningCond.L.Lock()
	t.Running--
	t.RunningCond.Signal()
	t.RunningCond.L.Unlock()
}

func (t *Threader) Wait() {
	t.RunningCond.L.Lock()
	for t.Running > 0 {
		t.RunningCond.Wait()
	}
	t.RunningCond.L.Unlock()
}
