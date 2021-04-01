package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLogAutoRotate(t *testing.T) {
	LogBasename = "logtest"
	LogLimit = 9
	nfiles := 6
	mbackups := 3
	expected := mbackups + 1
	dir, err := ioutil.TempDir("./", "logtest")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = os.RemoveAll(dir); err != nil {
			t.Log("remove dir error:", err.Error())
		}
	}()

	testlog := NewTmbLoggerWithDir(dir, LogBasename, LOG_DEBUG, 9)
	testlog.SetMaxBackups(mbackups)
	for i := 0; i < nfiles-1; i++ {
		testlog.LogError(fmt.Sprintf("%d_3456789", i))
	}
	t.Log("waiting 2s for all logs to gzip...")
	time.Sleep(2 * time.Second)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != expected {
		t.Error("incorrect number of files logged")
	}
	for _, f := range files {
		t.Log("found file:", f.Name())
	}
}

func TestLogRotate(t *testing.T) {
	LogBasename = "logtest"
	LogLimit = 9
	nfiles := 6
	mbackups := 3
	expected := mbackups + 2 // base file and log.5 (not gzipped)
	dir, err := ioutil.TempDir("./", "logtest")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = os.RemoveAll(dir); err != nil {
			t.Log("remove dir error:", err.Error())
		}
	}()

	testlog := NewTmbLoggerWithDir(dir, LogBasename, LOG_DEBUG, 9)
	for i := 0; i < nfiles-1; i++ {
		testlog.LogError("123456789")
	}
	t.Log("waiting 2s for all logs to gzip...")
	time.Sleep(2 * time.Second)

	// create much larger file to simulate past runs
	os.Create(filepath.Join(dir, LogBasename+".log.100.gz"))
	os.Create(filepath.Join(dir, LogBasename+".log.5"))

	// check log file count correct
	t.Log("Checking files before rotation...")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != nfiles+2 {
		t.Error("incorrect number of files logged")
	}
	for _, f := range files {
		t.Log("found file:", f.Name())
	}

	// do rotation
	t.Log("rotating files...")
	if _, err := rotateLogFiles(filepath.Join(dir, LogBasename+".log"), mbackups, nfiles-1); err != nil {
		t.Error(err)
	}

	// verify results
	t.Log("Checking files after rotation...")
	files, err = ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != expected {
		t.Error("incorrect number of files after rotation")
	}
	for _, f := range files {
		t.Log("found file:", f.Name())
	}
}
