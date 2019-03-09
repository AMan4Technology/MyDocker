package container

import (
    "fmt"
    "os"
    "os/exec"
    "syscall"

    log "github.com/sirupsen/logrus"
)

const (
    RootURL       = "/root"
    MntURL        = RootURL + "/mnt/%s"
    WriteLayerURL = RootURL + "/writeLayer/%s"
    ImageURL      = RootURL + "/image/%s"
)

/* 这里是父进程，也就是当前进程执行的内容
   1. 这里的/proc/self/exe调用中，/proc/self指定是当前运行进程自己的环境，exe其实就是自己
      调用了自己，使用这种方式对创建出来的进程进行初始化
   2. 后面的args是参数，其中init是传递给本进程的第一个参数，在本例中，其实就是会去调用
      initCommand去初始化进程的一些环境和资源
   3. 下面的clone参数就是去fork出来一个新进程，并且使用了namespace隔离新创建的进程和外部环境
   4. 如果用户指定了-ti参数，就需要把当前进程的输入输出导入到标准输入输出上 */
func NewParentProcess(tty bool, volume, containerName, imageName string, envs []string) (*exec.Cmd, *os.File) {
    readPipe, writePipe, err := NewPipe()
    if err != nil {
        log.Errorf("new pipe error %v", err)
        return nil, nil
    }
    initCmd, err := os.Readlink("/proc/self/exe")
    if err != nil {
        log.Errorf("Get init process failed. %v", err)
        return nil, nil
    }
    cmd := exec.Command(initCmd, "init")
    cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWUTS |
      syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET |
      syscall.CLONE_NEWIPC}
    if tty {
        cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
    } else {
        dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
        if err := mkdir(dirURL, 0622); err != nil {
            return nil, nil
        }
        logFileURL := dirURL + LogFile
        logFile, err := os.Create(logFileURL)
        if err != nil {
            log.Errorf("NewParentProcess create file %s failed. %v", logFileURL, err)
            return nil, nil
        }
        cmd.Stdout = logFile
    }
    cmd.ExtraFiles = []*os.File{readPipe}
    cmd.Dir = fmt.Sprintf(MntURL, containerName)
    cmd.Env = append(os.Environ(), envs...)
    NewWorkspace(volume, containerName, imageName)
    return cmd, writePipe
}

func NewPipe() (readPipe, writePipe *os.File, err error) {
    if readPipe, writePipe, err = os.Pipe(); err != nil {
        return nil, nil, err
    }
    return
}
