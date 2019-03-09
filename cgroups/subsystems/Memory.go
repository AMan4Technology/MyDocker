package subsystems

import (
    "strconv"
)

type Memory struct{}

func (*Memory) Name() string {
    return "memory"
}

func (s *Memory) Set(cgroupPath string, res *ResourceConfig) error {
    return addCgroupLimit(cgroupPath, "set", s.Name(), res.MemoryLimit, "memory.limit_in_bytes")
}

func (s *Memory) Apply(cgroupPath string, pid int) error {
    return addCgroupLimit(cgroupPath, "apply proc", s.Name(), strconv.Itoa(pid), "tasks")
}

func (s *Memory) Remove(cgroupPath string) error {
    return deleteCgroupLimit(cgroupPath, s.Name())
}
