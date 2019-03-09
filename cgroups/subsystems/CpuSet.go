package subsystems

import "strconv"

type CpuSet struct{}

func (*CpuSet) Name() string {
    return "cpuset"
}

func (s *CpuSet) Set(cgroupPath string, res *ResourceConfig) error {
    return addCgroupLimit(cgroupPath, "set", s.Name(), res.CpuSet, "cpuset.cpus")
}

func (s *CpuSet) Apply(cgroupPath string, pid int) error {
    return addCgroupLimit(cgroupPath, "apply proc", s.Name(), strconv.Itoa(pid), "tasks")
}

func (s *CpuSet) Remove(cgroupPath string) error {
    return deleteCgroupLimit(cgroupPath, s.Name())
}
