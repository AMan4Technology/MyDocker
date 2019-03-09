package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"

	"MyDocker/container"
	"MyDocker/templates"
)

func ListContainers() {
	files, err := ioutil.ReadDir(container.DefaultInfoDir)
	if err != nil {
		log.Errorf("Read dir %s failed. %v", container.DefaultInfoDir, err)
		return
	}
	infos := make([]*container.Info, 0, len(files))
	for _, file := range files {
		if file.Name() == "network" {
			continue
		}
		info, err := getContainerInfo(file.Name())
		if err == nil {
			infos = append(infos, info)
		}
	}
	w := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)
	templates.FPrintText(w, container.InfosID, container.InfosName, infos)
	if err := w.Flush(); err != nil {
		log.Errorf("Flush failed. %v", err)
		return
	}
}

func getContainerInfo(containerName string) (info *container.Info, err error) {
	configFileURL := fmt.Sprintf(container.DefaultInfoLocation, containerName) + container.ConfigName
	content, err := ioutil.ReadFile(configFileURL)
	if err != nil {
		log.Errorf("Read file %s failed. %v", configFileURL, err)
		return nil, err
	}
	info = new(container.Info)
	if err = json.Unmarshal(content, info); err != nil {
		log.Errorf("container %s json unmarshal failed. %v", containerName, err)
		return nil, err
	}
	return
}
