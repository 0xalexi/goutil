package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

const (
	defaultMaxBackups int = 500
)

type Logger struct {
	dir           string
	basename      string
	level         int
	byteLimit     int64
	noPrefix      bool
	logFilePrefix string
	logTimePrefix string
	logToStdout   bool
	mu            sync.Mutex
	bytecount     int64
	gznum         int
	maxbackups    int
	outlog        *os.File
	rotateQueue   chan int
}

func NewLoggerWithDir(dir, name string, level int, byteLimit int64) *Logger {
	if !strings.HasSuffix(name, ".log") {
		name += ".log"
	}

	l := &Logger{
		dir:           dir,
		basename:      name,
		level:         level,
		byteLimit:     byteLimit,
		logFilePrefix: logFilePrefix,
		logTimePrefix: logTimePrefix,
		maxbackups:    defaultMaxBackups,
		rotateQueue:   make(chan int),
	}
	if err := l.createLog(); err != nil {
		panic(err)
	}
	l.initStats()
	l.runRotator()
	return l
}

func New(name string, level int) *Logger {
	return NewLoggerWithDir(".", name, level, 1<<25)
}

func (l *Logger) SetLogLevel(level int) {
	l.level = level
}

func (l *Logger) SetNoPrefix(disabled bool) {
	l.noPrefix = disabled
}

func (l *Logger) SetMaxBackups(maxbackups int) {
	l.maxbackups = maxbackups
}

func (l *Logger) SetLogToStdout(logToStdout bool) {
	l.logToStdout = logToStdout
}

func (l *Logger) createLog() error {
	outlog, err := os.OpenFile(filepath.Join(l.dir, l.basename), os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0777)
	l.outlog = outlog
	return err
}

func (l *Logger) initStats() {
	outlogStat, err := l.outlog.Stat()
	if err != nil {
		panic(err)
	}
	l.bytecount = outlogStat.Size()
	l.gznum = 0
	for {
		_, err = os.Stat(fmt.Sprintf(filepath.Join(l.dir, l.basename+".%d.gz"), l.gznum+1))
		if os.IsNotExist(err) {
			break
		} else if err != nil {
			panic(err)
		} else {
			l.gznum++
		}
	}
}

func (l *Logger) runRotator() {
	go func() {
		for ngzip := range l.rotateQueue {
			n, err := rotateLogFiles(filepath.Join(l.dir, l.basename), l.maxbackups, ngzip)
			if err == nil {
				// remove num rotated from logger gznum
				l.mu.Lock()
				l.gznum -= (ngzip - n)
				l.mu.Unlock()
			}

		}
	}()
}

func (l *Logger) _doLog(stackSize int, level int, args ...interface{}) {
	if l == nil {
		DoLog(level, args...)
		return
	}
	if level > l.level {
		return
	}
	// Apply prefix args
	if !l.noPrefix {
		pc, file, line, ok := runtime.Caller(stackSize)
		file = filepath.Base(file)
		if ok {
			f := runtime.FuncForPC(pc)
			args = append(args, nil)
			copy(args[1:], args)
			args[0] = fmt.Sprintf(l.logFilePrefix, file, line, f.Name())
		}
	}
	args = append(args, nil)
	copy(args[1:], args)
	args[0] = time.Now().Format(l.logTimePrefix)
	levelStr := getLevelStr(level)
	if levelStr != "" {
		args = append(args, nil)
		copy(args[1:], args)
		args[0] = levelStr
	}
	// Do the needful
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.outlog == nil {
		fmt.Println(args...)
		return
	}
	n, err := fmt.Fprintln(l.outlog, args...)
	if err != nil {
		n, err = fmt.Fprintln(l.outlog, args...)
		if err != nil {
			panic(err)
		}
	}
	if l.logToStdout {
		fmt.Println(args...)
	}
	l.bytecount += int64(n)
	if l.bytecount >= l.byteLimit {
		newf := fmt.Sprintf(l.basename+".%d", l.gznum+1)
		l.outlog.Close()
		err := os.Rename(filepath.Join(l.dir, l.basename), filepath.Join(l.dir, newf))
		l.createLog()
		l.bytecount = 0
		if err == nil {
			l.gznum++
			go func(gznum int) {
				gzipOldLog(filepath.Join(l.dir, newf))
				if gznum > l.maxbackups {
					l.rotateQueue <- gznum
				}
			}(l.gznum)
		} else {
			fmt.Println("Failed to copy old logs. Erasing data instead:", err)
			l.outlog.Seek(0, os.SEEK_SET)
			err := l.outlog.Truncate(0)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (l *Logger) doLog(level int, args ...interface{}) {
	// stackSize: 0 is _doLog, 1 is doLog, 2 is parent logger.go func, 3 is actual caller
	l._doLog(3, level, args...)
}

type LoggerInterface interface {
	Fatal(v ...interface{})
	Error(v ...interface{})
	Warn(v ...interface{})
	Info(v ...interface{})
	Debug(v ...interface{})
	Trace(v ...interface{})

	DoLog(level int, v ...interface{})
	Log(v ...interface{})
	Println(v ...interface{})
	Print(v ...interface{})
	Printf(fmts string, v ...interface{})
}

func (l *Logger) DoLog(level int, v ...interface{}) {
	l.doLog(level, v...)
}

func (l *Logger) Log(v ...interface{}) {
	l.doLog(l.level, v...)
}

func (l *Logger) Println(v ...interface{}) {
	l.doLog(l.level, v...)
}

func (l *Logger) Print(v ...interface{}) {
	l.doLog(l.level, v...)
}

func (l *Logger) Printf(fmts string, v ...interface{}) {
	l.doLog(l.level, fmt.Sprintf(fmts, v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.doLog(LOG_ERROR, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.doLog(LOG_WARN, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.doLog(LOG_INFO, v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.doLog(LOG_DEBUG, v...)
}

func (l *Logger) Trace(v ...interface{}) {
	l.doLog(LOG_TRACE, v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.doLog(LOG_FATAL, v...)
	os.Exit(1)
}

func (l *Logger) Panic(v ...interface{}) {
	panic(fmt.Sprintln(v...))
}

func (l *Logger) StackTrace() {
	out := string(debug.Stack())
	for _, line := range strings.Split(out, "\n") {
		l.doLog(LOG_ERROR, "Stack trace: ", line)
	}
}
