package goutil

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type OpTrace struct {
	Function string        `json:"function,omitempty"`
	Duration time.Duration `json:"duration"`
	Children []*OpTrace    `json:"children,omitempty"`
	Count    int64         `json:"count,omitempty"`
	start    time.Time
	lock     *sync.RWMutex
}

func NewOpTrace(fname string) *OpTrace {
	return &OpTrace{Function: fname, start: time.Now(), lock: &sync.RWMutex{}}
}

func ident(lvl int) string {
	var tabs string
	for i := 0; i < lvl; i++ {
		tabs += "\t"
	}
	return tabs
}

func (l OpTrace) getString(lvl int) string {
	// if l.Duration == 0 {
	// 	l.CalcDuration()
	// }
	msg := fmt.Sprintf("%s: %v", l.Function, time.Duration(atomic.LoadInt64((*int64)(&l.Duration))))
	count := atomic.LoadInt64(&l.Count)
	if count > 0 {
		msg += fmt.Sprintf(", count: %d", count)
	}
	var tabs string = ident(lvl + 1)
	for _, child := range l.Children {
		msg += fmt.Sprintf("\n%s%s", tabs, child.getString(lvl+1))
	}
	return msg
}

func (l OpTrace) String() string {
	start := time.Now()
	str := l.getString(0)
	return fmt.Sprintf("process-trace: \n%s\ngen-string-time:%v", str, time.Now().Sub(start))
}

func (l *OpTrace) AddNewChild(fname string) *OpTrace {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.addNewChild(fname)
}

func (l *OpTrace) GetOrCreateChild(fname string) (*OpTrace, time.Time) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l == nil {
		return NewOpTrace(fname), time.Now()
	}
	for _, c := range l.Children {
		if c.Function == fname {
			return c, time.Now()
		}
	}
	return l.AddNewChild(fname), time.Now()
}

func (l *OpTrace) AddChild(child *OpTrace) {
	if l.Children == nil {
		l.Children = []*OpTrace{}
	}
	l.Children = append(l.Children, child)
	return
}

func (l *OpTrace) CalcDuration() *OpTrace {
	l.Duration = time.Now().Sub(l.start)
	return l
}

func (l *OpTrace) Restart() *OpTrace {
	return l.ResetTo(time.Now())
}

func (l *OpTrace) ResetTo(start time.Time) *OpTrace {
	l.start = start
	atomic.StoreInt64(&l.Count, 0)
	for _, c := range l.Children {
		c.ResetTo(start)
	}
	return l
}

// For repeating operations
func (l *OpTrace) Start() *OpTrace {
	l.start = time.Now()
	atomic.AddInt64(&l.Count, 1)
	return l
}

func (l *OpTrace) Stop() *OpTrace {
	atomic.AddInt64((*int64)(&l.Duration), int64(time.Now().Sub(l.start)))
	return l
}

// TODO: prevent double counting (see: fastfilter new-ff.handleIncReqCount trace)
func (l *OpTrace) Add(dur time.Duration) *OpTrace {
	atomic.AddInt64((*int64)(&l.Duration), int64(dur))
	atomic.AddInt64(&l.Count, 1)
	return l
}

func (l *OpTrace) AddFrom(start time.Time) *OpTrace {
	return l.Add(time.Now().Sub(start))
}

func (l *OpTrace) addNewChild(fname string) *OpTrace {
	child := NewOpTrace(fname)
	if l != nil {
		if l.Children == nil {
			l.Children = []*OpTrace{}
		}
		l.Children = append(l.Children, child)
	}
	return child
}
