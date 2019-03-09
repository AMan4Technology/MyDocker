package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func RunContainerInitProcess() (err error) {
	commands := readUserCommand()
	if len(commands) == 0 {
		return fmt.Errorf("run container get user command error: command is nil")
	}
	setUpMount()
	path, err := exec.LookPath(commands[0]) // 在系统的PATH里获取该命令的绝对路径
	if err != nil {
		log.Errorf("exec look path error %v", err)
		return
	}
	log.Infof("find path: %s", path)
	if err = syscall.Exec(path, commands, os.Environ()); err != nil {
		log.Error(err)
		return
	}
	return nil
}

func readUserCommand() []string {
	/* 每个进程默认都有3个文件描述符，分别是标准输入、输出、错误
	   通过cmd.ExtraFiles保存的新的文件描述符，index从3开始 */
	pipe := os.NewFile(uintptr(3), "pipe")
	defer pipe.Close()
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	return strings.Split(string(msg), " ")
}

// init挂载点
func setUpMount() {
	pwd, err := os.Getwd() // 获取当前路径
	if err != nil {
		log.Errorf("get current location error %v", err)
		return
	}
	log.Infof("current location is %s", pwd)
	pivotRoot(pwd)
	mountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(mountFlags), "")
	syscall.Mount("tmpfs", "/dev", "tmpfs",
		syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivotRoot(root string) (err error) {
	/* 为了使当前的root的oldRoot和newRoot不在同一文件系统下，我们把root重新mount了一次，
	   bind mount是把相同的内容换了一个挂载点的挂载方法 */
	if err = syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error %v", err)
	}
	// 创建rootfs/.pivot_root存储oldRoot
	pivotDir := filepath.Join(root, ".pivot_root")
	if err = mkdir(pivotDir, 0777); err != nil {
		return
	}
	/* pivotRoot到新的rootfs，oldRoot现在挂载在rootfs/.pivot_root上
	   挂载点目前仍然可以在mount命令中看到 */
	if err = syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}
	if err = syscall.Chdir("/"); err != nil { // 修改当前工作目录到根目录
		return fmt.Errorf("chdir / %v", err)
	}
	pivotDir = filepath.Join("/", ".pivot_root")
	if err = syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}
	return os.Remove(pivotDir) // 删除临时文件夹
}
