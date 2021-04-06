package log

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alexi/goutil"
)

const (
	LOG_FATAL int = iota - 1 // -1
	LOG_ERROR
	LOG_WARN
	LOG_INFO
	LOG_DEBUG
	LOG_TRACE
)

var (
	LogLevel    int
	LogLock     *sync.Mutex = new(sync.Mutex)
	Basename    string      = "app"
	LogLimit    int64       = 1 << 25
	MaxLogFiles int         = int(math.MaxInt64)
	logGzNum    int
	logCount    int64
	rotateQueue = make(chan int)

	//logTimePrefix string = "2006/01/02 15:04:05.000000000"
	logTimePrefix string = time.RFC3339Nano
	logFilePrefix string = "%[3]s %[1]s:%[2]d:" // [1] is file, [2] is line num, [3] is function
	useFilePrefix bool   = true
	outlog        *os.File

	QueueLogRate = 10 //secs

	LogProfile = false
)

func Run() {
	runProfile()
	runRotator()
	err := createLog()
	outlogStat, err := outlog.Stat()
	if err != nil {
		panic(err)
	}
	logCount = outlogStat.Size()
	logGzNum = 1
	for {
		_, err = os.Stat(fmt.Sprintf(Basename+".log.%d.gz", logGzNum))
		if os.IsNotExist(err) {
			break
		} else if err != nil {
			panic(err)
		} else {
			logGzNum++
		}
	}
}

func getLevelStr(level int) string {
	// Make sure to update with new log levels
	switch level {
	case LOG_FATAL:
		return "FATAL"
	case LOG_ERROR:
		return "ERROR"
	case LOG_WARN:
		return "WARN"
	case LOG_INFO:
		return "INFO"
	case LOG_DEBUG:
		return "DEBUG"
	case LOG_TRACE:
		return "TRACE"
	default:
		return ""
	}
}

func GetLevelString(level int) string {
	return getLevelStr(level)
}

func ParseLogLevel(lvl string) int {
	switch lvl {
	case "fatal", "FATAL":
		return LOG_FATAL
	case "error", "ERROR":
		return LOG_ERROR
	case "warn", "WARN":
		return LOG_WARN
	case "info", "INFO":
		return LOG_INFO
	case "debug", "DEBUG":
		return LOG_DEBUG
	case "trace", "TRACE":
		return LOG_TRACE
	default:
		return LOG_INFO
	}
}

func ResetLog() {
	outlog = nil
}

func createLog() (err error) {
	outlog, err = os.OpenFile(Basename+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0777)
	redirectStderr(outlog)
	return err
}

func gzipOldLog(f string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("gzipOldLog Fatal Error:", r)
		}
	}()
	inlog, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer inlog.Close()
	defer os.Remove(f)
	gziplog, err := os.Create(f + ".gz")
	if err != nil {
		panic(err)
	}
	defer gziplog.Close()
	gzipper, err := gzip.NewWriterLevel(gziplog, 8)
	if err != nil {
		panic(err)
	}
	defer gzipper.Close()
	io.Copy(gzipper, inlog)
	gzipper.Flush()
}

// Pass logfile path with non-numbered filename prefix (e.g. "./engie.log").
// Returns number of backups after rotation.
func rotateLogFiles(path string, maxbackups, ngzip int) (int, error) {
	dir := filepath.Dir(path)
	logbasename := filepath.Base(path)
	numbered_logbasename := logbasename + "."
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	var numbers = make([]string, ngzip)
	var todelete = make([]string, 0)
	for _, f := range files {
		if strings.HasPrefix(f.Name(), numbered_logbasename) {
			gzidx := strings.Index(f.Name(), ".gz")
			if gzidx > len(numbered_logbasename) {
				fn, err := strconv.Atoi(f.Name()[len(numbered_logbasename):gzidx])
				if err != nil {
					return 0, err
				}
				if fn > ngzip {
					todelete = append(todelete, f.Name())
					continue
				}
				numbers[fn-1] = f.Name()
			}
		}
	}
	for _, fname := range todelete {
		os.Remove(filepath.Join(dir, fname))
	}
	var newset = make([]string, 0, maxbackups)
	i := ngzip - 1
	for len(newset) < maxbackups && i >= 0 {
		if len(numbers[i]) > 0 {
			newset = append(newset, numbers[i])
		}
		i--
	}
	for i = len(newset) - 1; i >= 0; i-- {
		fname := newset[i]
		os.Rename(filepath.Join(dir, fname), filepath.Join(dir, fmt.Sprintf("%s.%d.gz", logbasename, maxbackups-i)))
	}
	return len(newset), nil

	// for i, fname := range numbers {
	// 	fmt.Println("idx:", i, "fname:", fname)
	// 	// if i < ngzip-maxbackups && len(fname) > 0 {
	// 	// 	os.Remove(filepath.Join(dir, fname))
	// 	// }
	// 	if i == len(numbers)-1 {
	// 		break
	// 	}
	// 	for j, fnamep := range numbers[i+1:] {
	// 		if len(fnamep) > 0 {
	// 			numbers[i] = fnamep
	// 			numbers[j] = ""
	// 			break
	// 		}
	// 	}
	// }
	// for i := 0; {
	// 	if len(fname) > 0 {
	// 		os.Rename(filepath.Join(dir, fname), filepath.Join(dir, fmt.Sprintf("%s.%d.gz", logbasename, i)))
	// 	}
	// }
	// return nil
}

func runProfile() {
	if LogProfile {
		out := "profile.out"
		f, err := os.Create(out)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
}

func ShouldLog(level int) bool {
	if level > LogLevel {
		return false
	}
	return true
}

func runRotator() {
	go func() {
		for ngzip := range rotateQueue {
			n, err := rotateLogFiles(Basename+".log", MaxLogFiles, ngzip)
			if err == nil {
				// remove num rotated from logger gznum
				LogLock.Lock()
				logGzNum -= (ngzip - n)
				LogLock.Unlock()
			} else {
				fmt.Println("rotate log files error:", err)
			}
		}
	}()
}

func doLog(level int, args ...interface{}) {
	// Apply LogLevel filter
	if level > LogLevel {
		return
	}
	// Apply prefix args
	if useFilePrefix {
		pc, file, line, ok := runtime.Caller(2) // 0 is doLog, 1 is internal logger.go func, 2 is caller
		file = filepath.Base(file)
		if ok {
			f := runtime.FuncForPC(pc)
			args = append(args, nil)
			copy(args[1:], args)
			args[0] = fmt.Sprintf(logFilePrefix, file, line, f.Name())
		}
	}
	args = append(args, nil)
	copy(args[1:], args)
	args[0] = time.Now().Format(logTimePrefix)
	levelStr := getLevelStr(level)
	if levelStr != "" {
		args = append(args, nil)
		copy(args[1:], args)
		args[0] = levelStr
	}
	// Do the needful
	LogLock.Lock()
	defer LogLock.Unlock()
	if outlog == nil {
		fmt.Println(args...)
		return
	}
	n, err := fmt.Fprintln(outlog, args...)
	if err != nil {
		n, err = fmt.Fprintln(outlog, args...)
		if err != nil {
			panic(err)
		}
	}
	logCount += int64(n)
	if logCount >= LogLimit {
		newf := fmt.Sprintf(Basename+".log.%d", logGzNum)
		outlog.Close()
		err := os.Rename(Basename+".log", newf)
		createLog()
		logCount = 0
		if err == nil {
			logGzNum++
			go func(gznum int) {
				gzipOldLog(newf)
				if gznum > MaxLogFiles {
					rotateQueue <- gznum
				}
			}(logGzNum)
		} else {
			fmt.Println("Failed to copy old logs. Erasing data instead:", err)
			truncateLog()
		}
	}
}

func DoLog(level int, v ...interface{}) {
	doLog(level, v...)
}

func Log(v ...interface{}) {
	doLog(LogLevel, v...)
}

func Logf(fmts string, v ...interface{}) {
	doLog(LogLevel, fmt.Sprintf(fmts, v...))
}

func LogError(v ...interface{}) {
	doLog(LOG_ERROR, v...)
}

func Error(v ...interface{}) {
	doLog(LOG_ERROR, v...)
}

func LogWarn(v ...interface{}) {
	doLog(LOG_WARN, v...)
}

func Warn(v ...interface{}) {
	doLog(LOG_WARN, v...)
}

func LogInfo(v ...interface{}) {
	doLog(LOG_INFO, v...)
}

func Info(v ...interface{}) {
	doLog(LOG_INFO, v...)
}

func LogDebug(v ...interface{}) {
	doLog(LOG_DEBUG, v...)
}

func Debug(v ...interface{}) {
	doLog(LOG_DEBUG, v...)
}

func LogTrace(v ...interface{}) {
	doLog(LOG_TRACE, v...)
}

func Trace(v ...interface{}) {
	doLog(LOG_TRACE, v...)
}

func StackTrace() {
	out := string(debug.Stack())
	for _, line := range strings.Split(out, "\n") {
		doLog(LOG_ERROR, "Stack trace: ", line)
	}
}

func LogFatal(v ...interface{}) {
	doLog(LOG_FATAL, v...)
	os.Exit(1)
}

func Fatal(v ...interface{}) {
	doLog(LOG_FATAL, v...)
	os.Exit(1)
}

func Panic(v ...interface{}) {
	panic(fmt.Sprintln(v...))
}

// For log functions that perform additional computation,
// only generate the log arguments when required.
func LogExec(level int, f func() []interface{}) {
	if level > LogLevel {
		return
	}
	doLog(level, f()...)
}

func LogTime(level int, fname string, f func()) {
	if level > LogLevel {
		f()
		return
	}
	start := time.Now()
	f()
	doLog(level, "time-log", fname, "time:", time.Now().Sub(start))
}

func LogStackTrace() {
	out := string(debug.Stack())
	for _, line := range strings.Split(out, "\n") {
		doLog(LOG_ERROR, "Stack trace: ", line)
	}
}

func CodeRefString() string {
	_, fn, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", fn, line)
}

func LogMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	doLog(LOG_DEBUG, "Memory usage:", m.Alloc, "Alloc,", m.TotalAlloc, "TotalAlloc")
}

type memstats struct {
	Alloc      uint64
	TotalAlloc uint64
	Sys        uint64
	Mallocs    uint64
	Lookups    uint64
}

func LogGC() {
	before, after := goutil.NewMemStats(), goutil.NewMemStats()
	before.Read()
	// Go 1.12 changed the default behavior on linux from `madv_dontneed` to `madv_free`,
	// effectively reducing the rate at which garbage collected memory is reclaimed by the OS.
	// https://golang.org/doc/go1.12#runtime
	//
	// FreeOSMemory forces a garbage collection followed by an attempt to return as much memory to the operating system as possible.
	debug.FreeOSMemory()
	after.Read()
	doLog(LOG_DEBUG, "Garbage Collection:", "before ->", before, "\t\t\t\t", "after ->", after)
}

func LogRequest(level int, r *http.Request) {
	if level > LogLevel {
		return
	}
	r.ParseForm()
	doLog(level, fmt.Sprintf("Request: %s %s %d\nUserAgent: %s\nRemote Addr: %s\nReferer: %s\nForm: %v", r.Method, r.RequestURI, r.ContentLength, r.UserAgent(), r.RemoteAddr, r.Referer(), r.Form))
}

func LogRequestComplete(level int, r *http.Request, start time.Time) {
	if level > LogLevel {
		return
	}
	r.ParseForm()
	doLog(level, fmt.Sprintf("Request: %s %s %d\nUserAgent: %s\nRemote Addr: %s\nReferer: %s\nForm: %v Command time: %v", r.Method, r.RequestURI, r.ContentLength, r.UserAgent(), r.RemoteAddr, r.Referer(), r.Form, time.Now().Sub(start)))
}

func NewError(v ...interface{}) error {
	return errors.New(fmt.Sprint(v...))
}

type LogManager struct {
	lvl int
}

func NewLogManager(lvl int) LogManager {
	return LogManager{lvl}
}

func (l LogManager) ShouldLog() bool {
	return l.lvl <= LogLevel
}

func (l LogManager) Log(v ...interface{}) {
	doLog(l.lvl, v...)
}

type TimeTracker struct {
	s time.Time
}

func NewTimeTracker() *TimeTracker {
	return &TimeTracker{time.Now()}
}

func (t *TimeTracker) Reset() {
	t.s = time.Now()
}

func (t *TimeTracker) Dur() time.Duration {
	return time.Now().Sub(t.s)
}

func (t *TimeTracker) DurReset() time.Duration {
	defer t.Reset()
	return time.Now().Sub(t.s)
}
