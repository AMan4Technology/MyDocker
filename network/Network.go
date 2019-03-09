package network

import (
    "encoding/json"
    "fmt"
    htmlTmp "html/template"
    "io/ioutil"
    "net"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"
    "text/tabwriter"
    "text/template"

    log "github.com/sirupsen/logrus"
    "github.com/vishvananda/netlink"
    "github.com/vishvananda/netns"

    "MyDocker/templates"

    "MyDocker/container"
)

func init() {
    templates.RegisterTextTmp(NetworksID, template.Must(template.ParseFiles(
        filepath.Join(TemplateDir, NetworksName))))
    templates.RegisterHtmlTmp(NetworksID, htmlTmp.Must(htmlTmp.ParseFiles(
        filepath.Join(TemplateDir, NetworksHTMLName))))
}

const (
    TemplateDir      = templates.BaseURL + "network/"
    NetworksID       = "network.Networks"
    NetworksName     = "Networks"
    NetworksHTMLName = NetworksName + templates.HTMLName
)

var (
    defaultNetworkDirURL = "/var/run/mydocker/network/network/"
    drivers              = make(map[string]Driver)
    networks             = make(map[string]*Network)
)

func Cerate(name, driver, subnet string) error {
    _, ipNet, _ := net.ParseCIDR(subnet)
    gatewayIP, err := ipAllocator.Allocate(ipNet)
    if err != nil {
        return err
    }
    ipNet.IP = gatewayIP
    nw, err := drivers[driver].Create(name, ipNet.String())
    if err != nil {
        return err
    }
    return nw.saveTo(defaultNetworkDirURL)
}

func Init() error {
    bridge := new(BridgeDriver)
    drivers[bridge.Name()] = bridge
    if err := initDir(defaultNetworkDirURL); err != nil {
        return err
    }
    return filepath.Walk(defaultNetworkDirURL, func(path string, info os.FileInfo, err error) error {
        if info.IsDir() {
            return nil
        }
        nw, err := from(path)
        if err != nil {
            log.Errorf("Load network file %s failed. %v", info.Name(), err)
            return nil
        }
        networks[nw.Name] = nw
        return nil
    })
}

func Delete(networkName string) (err error) {
    nw := networks[networkName]
    if nw == nil {
        return fmt.Errorf("no such network: %s", networkName)
    }
    if err = ipAllocator.Release(nw.IPRange, nw.IPRange.IP); err != nil {
        return fmt.Errorf("remove network gateway ip failed. error: %s", err)
    }
    if err = drivers[nw.Driver].Delete(nw); err != nil {
        return fmt.Errorf("remove network failed. error: %s", err)
    }
    return nw.removeFrom(defaultNetworkDirURL)
}

func Connect(networkName string, info *container.Info) (err error) {
    nw := networks[networkName]
    if nw == nil {
        return fmt.Errorf("no such network: %s", networkName)
    }
    ip, err := ipAllocator.Allocate(nw.IPRange)
    if err != nil {
        return
    }
    ep := &Endpoint{
        ID:           fmt.Sprintf("%s-%s", info.ID, networkName),
        IP:           ip,
        Network:      nw,
        PortMappings: info.PortMappings,
    }
    if err = drivers[nw.Driver].Connect(nw, ep); err != nil {
        return
    }
    if err = configEndpointIPAndRoute(ep, info); err != nil {
        return
    }
    return configPortMapping(ep, info)
}

func List() {
    tw := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
    templates.FPrintText(tw, NetworksID, NetworksName, networks)
    if err := tw.Flush(); err != nil {
        log.Errorf("Flush failed. %v", err)
        return
    }
}

type Network struct {
    Name    string     // 网络名
    Driver  string     // 网络驱动名
    IPRange *net.IPNet // 地址段
}

func (nw *Network) saveTo(dirURL string) (err error) {
    if err = initDir(dirURL); err != nil {
        return err
    }
    nwPath := filepath.Join(dirURL, nw.Name)
    nwFile, err := os.OpenFile(nwPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        log.Errorf("Dump network OpenFile %s failed. %v", nwPath, err)
        return
    }
    defer nwFile.Close()
    nwJSON, err := json.Marshal(nw)
    if err != nil {
        log.Errorf("Dump network marshal %s failed. %v", nw.Name, err)
        return
    }
    if _, err = nwFile.Write(nwJSON); err != nil {
        log.Errorf("Dump network write %s to file %s failed. %v", nwJSON, nwPath, err)
        return
    }
    return nil
}

func (nw *Network) removeFrom(dirURL string) error {
    nwPath := filepath.Join(dirURL, nw.Name)
    if _, err := os.Stat(nwPath); err != nil {
        if os.IsNotExist(err) {
            return nil
        }
        return err
    }
    return os.Remove(nwPath)
}

func from(nwFilePath string) (nw *Network, err error) {
    nw = &Network{}
    contentOfNw, err := ioutil.ReadFile(nwFilePath)
    if err != nil {
        log.Errorf("Read file %s failed. %v", nwFilePath, err)
        return
    }
    if err = json.Unmarshal(contentOfNw, nw); err != nil {
        log.Errorf("Unmarshal content %s failed. %v", contentOfNw, err)
        return
    }
    networks[nw.Name] = nw
    return
}

func configEndpointIPAndRoute(ep *Endpoint, info *container.Info) (err error) {
    veth, err := netlink.LinkByName(ep.Device.PeerName)
    if err != nil {
        return fmt.Errorf("config endpoint %s error: %v", ep.ID, err)
    }
    defer enterContainerNetNs(veth, info)()
    interfaceIP := *ep.Network.IPRange
    interfaceIP.IP = ep.IP
    if err = setInterfaceIP(ep.Device.PeerName, interfaceIP.String()); err != nil {
        return fmt.Errorf("%v error: %v", ep.Network, err)
    }
    if err = setInterfaceUp("lo"); err != nil {
        return
    }
    _, cidr, _ := net.ParseCIDR("0.0.0.0/0")
    if err = netlink.RouteAdd(&netlink.Route{
        LinkIndex: veth.Attrs().Index,
        Gw:        ep.Network.IPRange.IP,
        Dst:       cidr});
      err != nil {
        return
    }
    return nil
}

func configPortMapping(ep *Endpoint, info *container.Info) error {
    for _, pairOfPort := range ep.PortMappings {
        ports := strings.Split(pairOfPort, ":")
        if len(ports) != 2 {
            log.Errorf("Port mapping %s format error", pairOfPort)
            continue
        }
        ipTablesCmd := fmt.Sprintf("-t nat -A PREROUTING -p tcp -m tcp --dport %s -j DNAT --to-destination %s:%s",
            ports[0], ep.IP.String(), ports[1])
        if output, err := exec.Command("iptables", strings.Split(ipTablesCmd, " ")...).
            Output(); err != nil {
            log.Errorf("iptables output: %v", output)
        }
    }
    return nil
}

func initDir(dirURL string) (err error) {
    if _, err = os.Stat(dirURL); err == nil {
        return nil
    }
    if os.IsNotExist(err) {
        return os.MkdirAll(dirURL, 0644)
    }
    return
}

func enterContainerNetNs(link netlink.Link, info *container.Info) (exit func()) {
    exit = func() {}
    f, err := os.OpenFile(fmt.Sprintf("/proc/%s/ns/net", info.Pid), os.O_RDONLY, 0)
    if err != nil {
        log.Errorf("Get container %s net namespace failed. %v", info.Name, err)
        return
    }
    defer f.Close()
    nsFd := f.Fd()
    runtime.LockOSThread()
    if err = netlink.LinkSetNsFd(link, int(nsFd)); err != nil {
        log.Errorf("Set link %s net namespace failed. %v", link.Attrs().Name, err)
        return
    }
    origin, err := netns.Get()
    if err != nil {
        log.Errorf("Get current net namespace failed. %v", err)
        return
    }
    if err = netns.Set(netns.NsHandle(nsFd)); err != nil {
        log.Errorf("Set net namespace failed. %v", err)
        return
    }
    return func() {
        netns.Set(origin)
        origin.Close()
        runtime.UnlockOSThread()
    }
}
