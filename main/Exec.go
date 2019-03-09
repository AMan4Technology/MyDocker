package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"

	_ "MyDocker/nsenter"
)

const (
	EnvExecPid = "mydocker_pid" // 环境变量名mydocker_pid
	EnvExecCmd = "mydocker_cmd" // 环境变量名mydocker_cmd
)

func ExecContainer(containerName string, commands []string) {
	info, err := getContainerInfo(containerName)
	if err != nil {
		log.Errorf("Exec container %s getContainerInfo failed. %v", containerName, err)
		return
	}
	var (
		pid    = info.Pid
		cmdStr = strings.Join(commands, " ")
		cmd    = exec.Command("/proc/self/exe", "exec")
	)
	log.Infof("Container pid: %s | cmdStr: %s", pid, cmdStr)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.Env = append(os.Environ(), getEnvsByPid(pid)...)
	os.Setenv(EnvExecPid, pid)
	os.Setenv(EnvExecCmd, cmdStr)
	if err = cmd.Run(); err != nil {
		log.Errorf("Exec container %s failed. %v", containerName, err)
	}
}

func getEnvsByPid(pid string) []string {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Read file %s failed. %v", path, err)
		return nil
	}
	return strings.Split(string(contentBytes), "\u0000")
}
