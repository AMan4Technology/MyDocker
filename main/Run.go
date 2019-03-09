package main

import (
    "encoding/json"
    "fmt"
    "math/rand"
    "os"
    "strconv"
    "strings"
    "time"

    log "github.com/sirupsen/logrus"

    "MyDocker/network"

    "MyDocker/cgroups"
    "MyDocker/cgroups/subsystems"
    "MyDocker/container"
)

/* Run 这里的Start方法是真正开始前面创建好的command的调用，它首先会clone出来一个namespace隔离的
   进程，然后在子进程中，调用/proc/self/exe，也就是调用自己，发送init参数，调用我们写的init
   方法，去初始化容器的一些资源 */
func Run(tty bool, commands, envs, portMappings []string, res *subsystems.ResourceConfig, volume, containerName, imageName, nw string) {
    parent, writePipe := container.NewParentProcess(tty, volume, containerName, imageName, envs)
    if parent == nil {
        log.Error("New parent process error")
        return
    }
    if err := parent.Start(); err != nil {
        log.Error(err)
    }
    info, err := recordContainerInfo(parent.Process.Pid, containerName, volume, commands, portMappings)
    if err != nil {
        return
    }
    /* use mydocker-cgroup as cgroup name
       创建cgroup manager，并通过调用set和apply设置资源限制 */
    cgroupManager := cgroups.NewManager(info.ID)
    defer cgroupManager.Destroy()
    cgroupManager.Set(res)
    cgroupManager.Apply(parent.Process.Pid)
    if nw != "" {
        network.Init()
        if err = network.Connect(nw, &info); err != nil {
            log.Errorf("Connect to network %s failed. %v", nw, err)
            return
        }
    }
    sendInitCommand(commands, writePipe)
    if tty {
        parent.Wait()
        deleteContainerInfo(info.Name)
        container.DeleteWorkspace(volume, info.Name)
        os.Exit(0)
    }
}

func recordContainerInfo(pid int, containerName, volume string, commands, portMappings []string) (info container.Info, err error) {
    id := randStringBytes(10)
    if containerName == "" {
        containerName = id
    }
    info = container.Info{
        Pid:          strconv.Itoa(pid),
        ID:           id,
        Name:         containerName,
        Command:      strings.Join(commands, ""),
        CreatedTime:  time.Now().Format("2006-01-02 15:04:05"),
        Status:       container.Running,
        Volume:       volume,
        PortMappings: portMappings}
    jsonBytes, err := json.Marshal(info)
    if err != nil {
        log.Errorf("Record container %s info failed. %v", containerName, err)
        return
    }
    dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
    if err = os.MkdirAll(dirURL, 0622); err != nil {
        log.Errorf("Mkdir %s failed. %v", dirURL, err)
        return
    }
    fileURL := dirURL + container.ConfigName
    file, err := os.Create(fileURL)
    if err != nil {
        log.Errorf("Create file %s failed. %v", fileURL, err)
        return
    }
    defer file.Close()
    if _, err = file.Write(jsonBytes); err != nil {
        log.Errorf("Write info to file %s failed. %v", fileURL, err)
        return
    }
    return info, nil
}

func sendInitCommand(commands []string, writePipe *os.File) {
    defer writePipe.Close()
    command := strings.Join(commands, " ")
    log.Infof("Command is %s", command)
    writePipe.WriteString(command)
}

func deleteContainerInfo(containerName string) {
    dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
    if err := os.RemoveAll(dirURL); err != nil {
        log.Errorf("Remove dir %s failed. %v", dirURL, err)
    }
}

func randStringBytes(n int) string {
    var (
        letterBytes = "0123456789"
        max         = len(letterBytes)
    )
    rand.Seed(time.Now().UnixNano())
    randBytes := make([]byte, n)
    for i := range randBytes {
        randBytes[i] = letterBytes[rand.Intn(max)]
    }
    return string(randBytes)
}
