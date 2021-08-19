package goutil

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/minio/minio/pkg/disk"
)

func GetDiskAvailable() (uint64, error) {
	path, _ := os.Getwd()
	di, err := disk.GetInfo(path)
	if err != nil {
		return 0, err
	}
	return di.Free, nil
}

func GetPercentDiskAvailable() (float64, error) {
	path, _ := os.Getwd()
	di, err := disk.GetInfo(path)
	if err != nil {
		return 0, err
	}
	return (float64(di.Total-di.Free) / float64(di.Total)) * 100, nil
}

func GetDiskInfoString() string {
	path, _ := os.Getwd()
	di, err := disk.GetInfo(path)
	if err != nil {
		return fmt.Sprint("GetDiskInfoString err:", err)
	}
	percentage := (float64(di.Total-di.Free) / float64(di.Total)) * 100
	return fmt.Sprintf("%s of %s disk space used (%0.2f%%)\n",
		humanize.Bytes(di.Total-di.Free),
		humanize.Bytes(di.Total),
		percentage,
	)
}

func RunCmd(crashOnError bool, env map[string]string, dir string, command string) (string, error) {
	cmd, args := GenCmdArgs(command)
	return runCmd(crashOnError, env, dir, cmd, args...)
}

// No crash mode
func SafeRunCmd(env map[string]string, dir string, command string) (string, error) {
	cmd, args := GenCmdArgs(command)
	return runCmd(false, env, dir, cmd, args...)
}

func StartCmd(command string) error {
	cmd, args := GenCmdArgs(command)
	return startCmd(cmd, args...)
}

func StartRemoteCmd(ignoreErr bool, addr, sshKeyFilePath, remoteCmd string) string {
	err := startCmd("ssh", "-i", sshKeyFilePath, "-o", "StrictHostKeyChecking=no", addr, remoteCmd)
	if err != nil {
		if !ignoreErr {
			panic(err)
		}
		return err.Error()
	}
	return ""
}

// addr in format ubuntu@127.0.0.1
func ExecRemoteCmd(crashOnError bool, addr, sshKeyFilePath, remoteCmd string) string {
	out, _ := runCmd(crashOnError, nil, "", "ssh", "-i", sshKeyFilePath, "-o", "StrictHostKeyChecking=no", addr, remoteCmd)
	return out
}

func SafeExecRemoteCmd(addr, sshKeyFilePath, remoteCmd string) (string, error) {
	return runCmd(false, nil, "", "ssh", "-i", sshKeyFilePath, "-o", "StrictHostKeyChecking=no", addr, remoteCmd)
}

func GenCmdArgs(cmd string) (string, []string) {
	vals := strings.Split(cmd, " ")
	if len(vals) > 1 {
		return vals[0], vals[1:]
	}
	return vals[0], []string{}
}

func PipeCmds(cmd1 string, cmd2 string) error {
	args1 := strings.Split(cmd1, " ")
	args2 := strings.Split(cmd2, " ")
	c1 := exec.Command(args1[0], args1[1:]...)
	c2 := exec.Command(args2[0], args2[1:]...)

	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()
	io.Copy(os.Stdout, &b2)
	return nil
}

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

func runCmd(crashOnError bool, env map[string]string, dir string, name string, cmds ...string) (string, error) {
	cmd := exec.Command(name, cmds...)
	for key, value := range env {
		cmd.Env = append(cmd.Env, key+"="+value)
	}
	if dir != "" {
		cmd.Dir = dir
	}
	output, err := cmd.CombinedOutput()
	var resp string
	buffer := bufio.NewScanner(bytes.NewReader(output))
	for buffer.Scan() {
		t := buffer.Text()
		resp += "\n" + t
	}
	// ssh might fail on first call when making remote call
	if err != nil && !strings.Contains(strings.ToLower(resp), "connection refused") && !strings.Contains(strings.ToLower(resp), "cannot stat") {
		if crashOnError {
			panic(err)
		}
		return resp, err
	}

	if err = buffer.Err(); err != io.EOF && err != nil {
		if crashOnError {
			panic(err)
		}
		return resp, err
	}
	return resp, nil
}

func startCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Start()
}
