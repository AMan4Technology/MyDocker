package main

import (
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"MyDocker/container"
)

func commitContainer(containerName, imageName string) {
	var (
		mntURL   = fmt.Sprintf(container.MntURL+"/", containerName)
		imageTar = fmt.Sprintf(container.ImageURL+".tar", imageName)
	)
	log.Infof("Commit to %s", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").
		CombinedOutput(); err != nil {
		log.Errorf("Tar folder %s failed. %v", mntURL, err)
	}
}
