package subsystems

import log "github.com/sirupsen/logrus"

// 通过不同的Subsystem初始化创建资源限制处理链数组
var subsystems = []Subsystem{
    new(Memory),
    new(Cpu),
    new(CpuSet)}

func Range(callback func(subsystem Subsystem) (string, error)) {
    for _, subsystem := range subsystems {
        if option, err := callback(subsystem); err != nil {
            log.Warnf("%s cgroup fail %v", option, err)
        }
    }
}

// Subsystem接口，每个Subsystem应该具备以下4个method
type Subsystem interface {
    // 返回SubSystem的名字，例如：memory、cpu
    Name() string
    // 设置某个cgroup在这个SubSystem中的资源限制
    Set(cgroupPath string, res *ResourceConfig) error
    // 将进程添加到某个cgroup中
    Apply(cgroupPath string, pid int) error
    // 移除某个cgroup
    Remove(cgroupPath string) error
}
