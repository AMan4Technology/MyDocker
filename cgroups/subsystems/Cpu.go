package subsystems

import "strconv"

type Cpu struct{}

func (*Cpu) Name() string {
    return "cpu"
}

func (s *Cpu) Set(cgroupPath string, res *ResourceConfig) error {
    return addCgroupLimit(cgroupPath, "set", s.Name(), res.CpuShare, "cpu.shares")
}

func (s *Cpu) Apply(cgroupPath string, pid int) error {
    return addCgroupLimit(cgroupPath, "apply proc", s.Name(), strconv.Itoa(pid), "tasks")
}

func (s *Cpu) Remove(cgroupPath string) error {
    return deleteCgroupLimit(cgroupPath, s.Name())
}
