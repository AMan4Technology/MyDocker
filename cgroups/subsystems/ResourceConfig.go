package subsystems

// 用于传递资源限制配置的结构体，包含内存限制、CPU时间片权重、CPU核心数
type ResourceConfig struct {
    MemoryLimit, CpuShare, CpuSet string
}
