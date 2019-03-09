package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"MyDocker/container"
)

func removeContainer(containerName string) {
	info, err := getContainerInfo(containerName)
	if err != nil {
		log.Errorf("Remove container %s getContainerInfo failed. %v", containerName, err)
		return
	}
	if info.Status != container.Stop {
		log.Errorf("Container %s status is not stop.", containerName)
		return
	}
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err = os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove dir %s failed. %v", dirURL, err)
		return
	}
	container.DeleteWorkspace(info.Volume, containerName)
}
