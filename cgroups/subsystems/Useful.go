package subsystems

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
    "strings"
)

func GetCgroupPath(subsystem, cgroupPath string, autoCreate bool) (string, error) {
    cgroupRoot := FindCgroupMountPoint(subsystem)
    if cgroupRoot == "" {
        return "", fmt.Errorf("subsystem %s not exist", subsystem)
    }
    cgroup := path.Join(cgroupRoot, cgroupPath)
    if _, err := os.Stat(cgroup); err != nil {
        if !os.IsNotExist(err) || !autoCreate {
            return "", fmt.Errorf("cgroup path error %v", err)
        }
        if err = os.Mkdir(cgroup, 0755); err != nil {
            return "", fmt.Errorf("create cgroup error %v", err)
        }
    }
    return cgroup, nil
}

/* 通过/proc/self/mountinfo找出挂载了某个subsystem的hierarchy cgroup根节点所在的目录
   FindCgroupMountPoint("memory") */
func FindCgroupMountPoint(subsystem string) string {
    f, err := os.Open("/proc/self/mountinfo")
    if err != nil {
        return ""
    }
    defer f.Close()
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        var (
            text   = scanner.Text()
            fields = strings.Split(text, " ")
        )
        if strings.Contains(","+fields[len(fields)-1]+",", ","+subsystem+",") {
            return fields[4]
        }
    }
    scanner.Err()
    return ""
}

func addCgroupLimit(cgroupPath, option, name, limit, fileName string) error {
    if limit == "" {
        return nil
    }
    subsystemCgroupPath, err := GetCgroupPath(name, cgroupPath, true)
    if err != nil {
        return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
    }
    err = ioutil.WriteFile(filepath.Join(subsystemCgroupPath, fileName),
        []byte(limit), 0644)
    if err != nil {
        return fmt.Errorf("%s cgroup %s fail %v", name, option, err)
    }
    return nil
}

func deleteCgroupLimit(cgroupPath, name string) error {
    subsystemCgroupPath, err := GetCgroupPath(name, cgroupPath, false)
    if err != nil {
        return err
    }
    return os.Remove(subsystemCgroupPath)
}
