package network

import (
    "fmt"
    "net"
    "os/exec"
    "strings"

    log "github.com/sirupsen/logrus"
    "github.com/vishvananda/netlink"
)

type BridgeDriver struct{}

func (d BridgeDriver) Name() string {
    return "bridge"
}

func (d BridgeDriver) Create(name, subnet string) (nw *Network, err error) {
    ip, ipRange, _ := net.ParseCIDR(subnet)
    ipRange.IP = ip
    nw = &Network{name, d.Name(), ipRange}
    if err = d.initBridge(nw); err != nil {
        log.Errorf("Init bridge failed. %v", err)
        return nil, err
    }
    return
}

func (d BridgeDriver) Delete(nw *Network) error {
    link, err := netlink.LinkByName(nw.Name)
    if err != nil {
        return err
    }
    return netlink.LinkDel(link)
}

func (d BridgeDriver) Connect(nw *Network, ep *Endpoint) (err error) {
    link, err := netlink.LinkByName(nw.Name)
    if err != nil {
        return
    }
    la := netlink.NewLinkAttrs()
    la.Name = ep.ID[:5]
    la.MasterIndex = link.Attrs().Index
    ep.Device = netlink.Veth{
        LinkAttrs: la,
        PeerName:  "cif-" + la.Name,
    }
    if err = netlink.LinkAdd(&ep.Device); err != nil {
        return fmt.Errorf("add endpoint device veth %s error: %v", ep.Device.PeerName, err)
    }
    if err = netlink.LinkSetUp(&ep.Device); err != nil {
        return fmt.Errorf("set endpoint device veth %s up error: %v", ep.Device.PeerName, err)
    }
    return nil
}

func (d BridgeDriver) Disconnect(nw *Network, ep *Endpoint) error {
    panic("implement me")
}

func (d BridgeDriver) initBridge(nw *Network) (err error) {
    if err = createBridgeInterface(nw.Name); err != nil {
        return fmt.Errorf("add bridge %s error: %v", nw.Name, err)
    }
    ipRangeStr := nw.IPRange.String()
    if err = setInterfaceIP(nw.Name, ipRangeStr); err != nil {
        return fmt.Errorf("assigning address %s on bridge %s error: %v",
            ipRangeStr, nw.Name, err)
    }
    if err = setInterfaceUp(nw.Name); err != nil {
        return fmt.Errorf("set up bridge %s error: %v", nw.Name, err)
    }
    if err = setupIPTables(nw.Name, nw.IPRange); err != nil {
        return fmt.Errorf("setting iptables for %s error: %v", nw.Name, err)
    }
    return
}

func setupIPTables(bridgeName string, subnet *net.IPNet) error {
    iptablesCmd := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE",
        subnet.String(), bridgeName)
    if output, err := exec.Command("iptables", strings.Split(iptablesCmd, " ")...).
        Output(); err != nil {
        log.Errorf("iptables Output: %v", output)
        return err
    }
    return nil
}

func createBridgeInterface(bridgeName string) (err error) {
    if _, err = net.InterfaceByName(bridgeName); err == nil ||
      !strings.Contains(err.Error(), "no such network interface") {
        return fmt.Errorf("bridge %s is exist, error: %v", bridgeName, err)
    }
    lA := netlink.NewLinkAttrs()
    lA.Name = bridgeName
    if err = netlink.LinkAdd(&netlink.Bridge{LinkAttrs: lA}); err != nil {
        return fmt.Errorf("create bridge %s error: %v", bridgeName, err)
    }
    return
}

func setInterfaceIP(name, rawIP string) (err error) {
    link, err := netlink.LinkByName(name)
    if err != nil {
        return fmt.Errorf("get interface %s error: %v", name, err)
    }
    ipNet, err := netlink.ParseIPNet(rawIP)
    if err != nil {
        return
    }
    return netlink.AddrAdd(link, &netlink.Addr{
        IPNet:     ipNet,
        Label:     "",
        Flags:     0,
        Scope:     0,
        Broadcast: nil})
}

func setInterfaceUp(name string) (err error) {
    link, err := netlink.LinkByName(name)
    if err != nil {
        return fmt.Errorf("retrieving a link named [%s] error: %v", name, err)
    }
    if err = netlink.LinkSetUp(link); err != nil {
        return fmt.Errorf("enabling link %s error: %v", name, err)
    }
    return
}
