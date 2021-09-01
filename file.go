package goutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func SaveFile(dat []byte, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(dat), 0644)
}

func CheckFileExists(filepath string) (bool, error) {
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CheckFilesExist(filepaths []string) (bool, error) {
	for _, fpath := range filepaths {
		if ok, err := CheckFileExists(fpath); !ok {
			return ok, err
		}
	}
	return true, nil
}

// Copy Pasta'd from https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func GetFileNumBytes(path string) (size int64, exists bool, err error) {
	exists, err = CheckFileExists(path)
	if !exists || err != nil {
		return
	}
	var file *os.File
	if file, err = os.Open(path); err != nil {
		return
	}
	stat, _ := file.Stat()
	size = stat.Size()
	return
}

func GetFileInfo(path string) (stat os.FileInfo, exists bool, err error) {
	exists, err = CheckFileExists(path)
	if !exists || err != nil {
		return
	}
	var file *os.File
	if file, err = os.Open(path); err != nil {
		return
	}
	stat, _ = file.Stat()
	return
}

func RemoveDirectoryContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
