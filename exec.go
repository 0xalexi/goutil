package goutil

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
)

// No crash mode
func RunCmd(env map[string]string, dir string, command string) (string, error) {
	cmd, args := GenCmdArgs(command)
	return runCmd(false, env, dir, cmd, args...)
}

func UnsafeRunCmd(env map[string]string, dir string, command string) (string, error) {
	cmd, args := GenCmdArgs(command)
	return runCmd(true, env, dir, cmd, args...)
}

func StartCmd(command string) error {
	cmd, args := GenCmdArgs(command)
	return startCmd(cmd, args...)
}

func StartRemoteCmd(ignoreErr bool, addr, sshKeyFilePath, remoteCmd string) error {
	return startCmd("ssh", "-i", sshKeyFilePath, "-o", "StrictHostKeyChecking=no", addr, remoteCmd)
}

func UnsafeStartRemoteCmd(ignoreErr bool, addr, sshKeyFilePath, remoteCmd string) {
	err := startCmd("ssh", "-i", sshKeyFilePath, "-o", "StrictHostKeyChecking=no", addr, remoteCmd)
	if err != nil {
		panic(err)
	}
}

// addr in format ubuntu@127.0.0.1
func RunRemoteCmd(addr, sshKeyFilePath, remoteCmd string) string {
	out, _ := runCmd(false, nil, "", "ssh", "-i", sshKeyFilePath, "-o", "StrictHostKeyChecking=no", addr, remoteCmd)
	return out
}

func UnsafeRunRemoteCmd(addr, sshKeyFilePath, remoteCmd string) (string, error) {
	return runCmd(true, nil, "", "ssh", "-i", sshKeyFilePath, "-o", "StrictHostKeyChecking=no", addr, remoteCmd)
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
