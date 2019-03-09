package cgroups

import (
    "MyDocker/cgroups/subsystems"
)

func NewManager(path string) *Manager {
    return &Manager{Path: path}
}

type Manager struct {
    Path     string
    Resource *subsystems.ResourceConfig
}

func (m *Manager) Apply(pid int) error {
    subsystems.Range(func(subsystem subsystems.Subsystem) (string, error) {
        return "apply", subsystem.Apply(m.Path, pid)
    })
    return nil
}

func (m *Manager) Set(res *subsystems.ResourceConfig) error {
    subsystems.Range(func(subsystem subsystems.Subsystem) (string, error) {
        return "set", subsystem.Set(m.Path, res)
    })
    return nil
}

func (m *Manager) Destroy() error {
    subsystems.Range(func(subsystem subsystems.Subsystem) (string, error) {
        return "remove", subsystem.Remove(m.Path)
    })
    return nil
}
