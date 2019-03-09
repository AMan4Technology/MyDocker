package container

import (
    "fmt"
    "os"
    "os/exec"
    "strings"

    log "github.com/sirupsen/logrus"
)

// Create a AUFS filesystem as container root workspace
func NewWorkspace(volume, containerName, imageName string) {
    CreateReadOnlyLayer(imageName)
    CreateWriteLayer(containerName)
    CreateMountPoint(containerName, imageName)
    if volume == "" {
        return
    }
    volumeURLs := strings.Split(volume, ":")
    if len(volumeURLs) != 2 || volumeURLs[0] == "" || volumeURLs[1] == "" {
        log.Warn("Volume argument input is not correct.")
        return
    }
    MountVolume(containerName, volumeURLs)
    log.Infof("Mount volume %q", volumeURLs)
}

// 将busybox.tar解压到busybox目录下，作为容器的只读层
func CreateReadOnlyLayer(imageName string) (err error) {
    var (
        imageURL    = fmt.Sprintf(ImageURL+"/", imageName)
        imageTarURL = fmt.Sprintf(ImageURL+".tar", imageName)
    )
    exist, err := PathExists(imageURL)
    if err != nil {
        log.Infof("Fail to judge whether dir %s exists. %v", imageURL, err)
        return
    }
    if !exist {
        if err = mkdir(imageURL, 0777); err != nil {
            return
        }
        if _, err = exec.Command("tar", "-xvf", imageTarURL, "-C", imageURL).
            CombinedOutput(); err != nil {
            log.Errorf("unTar dir %s failed. %v", imageTarURL, err)
            return
        }
    }
    return nil
}

func CreateWriteLayer(containerName string) error {
    return mkdir(fmt.Sprintf(WriteLayerURL, containerName), 0777)
}

func CreateMountPoint(containerName, imageName string) (err error) {
    // 创建mnt文件夹作为挂载点
    mntURL := fmt.Sprintf(MntURL, containerName)
    if err = mkdir(mntURL, 0777); err != nil {
        return
    }
    // 把writeLayer目录和busybox目录mount到mnt目录下
    if _, err = exec.Command("mount", "-t", "aufs", "-o",
        fmt.Sprintf("dirs=%s:%s", fmt.Sprintf(WriteLayerURL, containerName),
            fmt.Sprintf(ImageURL, imageName)),
        "none", mntURL).CombinedOutput(); err != nil {
        log.Errorf("mount to %s failed. %v", mntURL, err)
        return
    }
    return nil
}

func MountVolume(containerName string, volumeURLs []string) error {
    var (
        parentURL    = volumeURLs[0]                                      // 宿主机文件目录
        containerURL = fmt.Sprintf(MntURL, containerName) + volumeURLs[1] // 容器文件系统内的挂载点
    )
    mkdir(parentURL, 0777)
    mkdir(containerURL, 0777)
    if _, err := exec.Command("mount", "-t", "aufs", "-o",
        "dirs="+parentURL, "none", containerURL).CombinedOutput(); err != nil {
        log.Errorf("Mount volume %s to %s failed. %v", parentURL, containerURL, err)
        return err
    }
    return nil
}

func PathExists(url string) (exist bool, err error) {
    if _, err = os.Stat(url); err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return
}

// Delete the AUFS filesystem while container exit
func DeleteWorkspace(volume, containerName string) {
    if volume != "" {
        volumeURLs := strings.Split(volume, ":")
        if len(volumeURLs) == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
            UnmountVolumeMountPoint(fmt.Sprintf(MntURL, containerName), volumeURLs[1])
        }
    }
    DeleteMountPoint(containerName)
    DeleteWriteLayer(containerName)
}

func UnmountVolumeMountPoint(mntURL, volumeURL string) error {
    return unmountDir(mntURL + volumeURL)
}

func DeleteMountPoint(containerName string) error {
    mntURL := fmt.Sprintf(MntURL, containerName)
    if err := unmountDir(mntURL); err != nil {
        return err
    }
    return removeDir(mntURL)
}

func DeleteWriteLayer(containerName string) error {
    return removeDir(fmt.Sprintf(WriteLayerURL, containerName))
}

func mkdir(url string, perm os.FileMode) error {
    if err := os.MkdirAll(url, perm); err != nil {
        log.Errorf("Mkdir %s failed. %v", url, err)
        return err
    }
    return nil
}

func unmountDir(url string) error {
    cmd := exec.Command("umount", url)
    cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
    if err := cmd.Run(); err != nil {
        log.Errorf("Unmount %s failed. %v", url, err)
        return err
    }
    return nil
}

func removeDir(url string) error {
    if err := os.RemoveAll(url); err != nil {
        log.Errorf("Remove dir %s failed. %v", url, err)
        return err
    }
    return nil
}
