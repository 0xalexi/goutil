package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	defaultMaxBackups int = 500
)

type TmbLogger struct {
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

func NewTmbLoggerWithDir(dir, name string, level int, byteLimit int64) *TmbLogger {
	if !strings.HasSuffix(name, ".log") {
		name += ".log"
	}

	l := &TmbLogger{
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

func NewTmbLogger(name string, level int, byteLimit int64) *TmbLogger {
	return NewTmbLoggerWithDir(".", name, level, byteLimit)
}

func (l *TmbLogger) SetLogLevel(level int) {
	l.level = level
}

func (l *TmbLogger) SetNoPrefix(disabled bool) {
	l.noPrefix = disabled
}

func (l *TmbLogger) SetMaxBackups(maxbackups int) {
	l.maxbackups = maxbackups
}

func (l *TmbLogger) SetLogToStdout(logToStdout bool) {
	l.logToStdout = logToStdout
}

func (l *TmbLogger) createLog() error {
	outlog, err := os.OpenFile(filepath.Join(l.dir, l.basename), os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0777)
	l.outlog = outlog
	return err
}

func (l *TmbLogger) initStats() {
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

func (l *TmbLogger) runRotator() {
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

func (l *TmbLogger) doLog(level int, args ...interface{}) {
	if l == nil {
		DoLog(level, args...)
		return
	}
	if level > l.level {
		return
	}
	// Apply prefix args
	if !l.noPrefix {
		pc, file, line, ok := runtime.Caller(2) // 0 is doLog, 1 is internal logger.go func, 2 is caller
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

type TmbLoggerInterface interface {
	LogFatal(v ...interface{})
	LogError(v ...interface{})
	LogWarn(v ...interface{})
	LogInfo(v ...interface{})
	LogDebug(v ...interface{})
	LogTrace(v ...interface{})

	DoLog(level int, v ...interface{})
	Log(v ...interface{})
	Println(v ...interface{})
	Print(v ...interface{})
	Printf(fmts string, v ...interface{})
}

func (l *TmbLogger) DoLog(level int, v ...interface{}) {
	l.doLog(level, v...)
}

func (l *TmbLogger) Log(v ...interface{}) {
	l.doLog(l.level, v...)
}

func (l *TmbLogger) Println(v ...interface{}) {
	l.doLog(l.level, v...)
}

func (l *TmbLogger) Print(v ...interface{}) {
	l.doLog(l.level, v...)
}

func (l *TmbLogger) Printf(fmts string, v ...interface{}) {
	l.doLog(l.level, fmt.Sprintf(fmts, v...))
}

func (l *TmbLogger) LogError(v ...interface{}) {
	l.doLog(LOG_ERROR, v...)
}

func (l *TmbLogger) LogWarn(v ...interface{}) {
	l.doLog(LOG_WARN, v...)
}

func (l *TmbLogger) LogInfo(v ...interface{}) {
	l.doLog(LOG_INFO, v...)
}

func (l *TmbLogger) LogDebug(v ...interface{}) {
	l.doLog(LOG_DEBUG, v...)
}

func (l *TmbLogger) LogTrace(v ...interface{}) {
	l.doLog(LOG_TRACE, v...)
}

func (l *TmbLogger) LogFatal(v ...interface{}) {
	l.doLog(LOG_FATAL, v...)
	os.Exit(1)
}
