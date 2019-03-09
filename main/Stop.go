package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"syscall"

	log "github.com/sirupsen/logrus"

	"MyDocker/container"
)

func stopContainer(containerName string) {
	info, err := getContainerInfo(containerName)
	if err != nil {
		log.Errorf("Stop container %s getContainerInfo failed. %v", containerName, err)
		return
	}
	pid, err := strconv.Atoi(info.Pid)
	if err != nil {
		log.Errorf("convert pid %s from string to int failed. %v", info.Pid, err)
		return
	}
	if err = syscall.Kill(pid, syscall.SIGTERM); err != nil {
		log.Errorf("Stop container %s process %s failed. %v", containerName, info.Pid, err)
		return
	}
	info.Status, info.Pid = container.Stop, " "
	infoJSON, err := json.Marshal(info)
	if err != nil {
		log.Errorf("Json marshal container %s info failed. %v", containerName, err)
		return
	}
	configFileURL := fmt.Sprintf(container.DefaultInfoLocation, containerName) + container.ConfigName
	if err = ioutil.WriteFile(configFileURL, infoJSON, 0622); err != nil {
		log.Errorf("Write to file %s failed. %v", configFileURL, err)
	}
}
