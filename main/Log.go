package main

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	"MyDocker/container"
)

func logOfContainer(containerName string) {
	file, err := os.Open(fmt.Sprintf(container.DefaultInfoLocation, containerName) + container.LogFile)
	if err != nil {
		log.Errorf("Open container %s log file failed. %v", containerName, err)
		return
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("Read container %s log file failed. %v", containerName, err)
		return
	}
	fmt.Fprintln(os.Stdout, string(content))
}
