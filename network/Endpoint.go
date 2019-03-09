package network

import (
    "net"

    "github.com/vishvananda/netlink"
)

type Endpoint struct {
    ID           string           `json:"id"`
    Device       netlink.Veth     `json:"dev"`
    IP           net.IP           `json:"ip"`
    Mac          net.HardwareAddr `json:"mac"`
    PortMappings []string         `json:"port_mappings"`
    Network      *Network
}
