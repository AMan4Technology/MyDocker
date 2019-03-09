package network

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "net"
    "os"
    "path/filepath"

    log "github.com/sirupsen/logrus"

    "MyDocker/datastructure"
)

const ipamDefaultAllocatorURL = "/var/run/mydocker/network/ipam/subnet.json"

var ipAllocator = IPAM{
    SubnetAllocatorURL: ipamDefaultAllocatorURL,
    Subnets:            make(map[string]*datastructure.Bitmap)}

type IPAM struct {
    SubnetAllocatorURL string
    Subnets            map[string]*datastructure.Bitmap
}

func (ipam *IPAM) Allocate(subnet *net.IPNet) (ip net.IP, err error) {
    if err = ipam.load(); err != nil {
        log.Errorf("Load allocation %s info failed. %v", ipam.SubnetAllocatorURL, err)
    }
    var (
        one, size = subnet.Mask.Size()
        subnetStr = subnet.String()
        bm        = ipam.Subnets[subnetStr]
    )
    if bm == nil {
        newBm := datastructure.NewBitmap(1 << (uint8(size - one)))
        bm = &newBm
        bm.Save(0)
        ipam.Subnets[subnetStr] = bm
    }
    if bm.Full() {
        err = errors.New("subnet can't provide more ip")
        log.Errorf("Allocate ip from %s failed. %v", subnetStr, err)
        return nil, err
    }
    for i := 1; i < bm.Cap(); i++ {
        if !bm.Have(i) {
            bm.Save(i)
            ip = make([]byte, 4)
            for j := len(ip) - 1; j >= 0; j-- {
                ip[j] = subnet.IP[j] + uint8(i)
                i = i >> 8
            }
            break
        }
    }
    if err = ipam.save(); err != nil {
        log.Errorf("Save to file failed. %v", err)
        return ip, nil
    }
    return
}

func (ipam *IPAM) Release(subnet *net.IPNet, ip net.IP) (err error) {
    if err = ipam.load(); err != nil {
        log.Errorf("Load allocation %s info failed. %v", ipam.SubnetAllocatorURL, err)
    }
    ip = ip.To4()
    var value int
    for i, length := 0, len(ip); i < length; i++ {
        value = value<<8 + int(ip[i]-subnet.IP[i])
    }
    ipam.Subnets[subnet.String()].Remove(value)
    if err = ipam.save(); err != nil {
        log.Errorf("Save to file failed. %v", err)
    }
    return nil
}

func (ipam *IPAM) load() (err error) {
    if _, err = os.Stat(ipam.SubnetAllocatorURL); err != nil {
        if os.IsNotExist(err) {
            return nil
        }
        return
    }
    content, err := ioutil.ReadFile(ipam.SubnetAllocatorURL)
    if err != nil {
        return
    }
    if err = json.Unmarshal(content, &ipam.Subnets); err != nil {
        log.Errorf("unmarshal %s to ipam failed. %v", content, err)
        return
    }
    return nil
}

func (ipam IPAM) save() (err error) {
    ipamDirURL, _ := filepath.Split(ipam.SubnetAllocatorURL)
    if _, err = os.Stat(ipamDirURL); err != nil {
        if !os.IsNotExist(err) {
            return
        }
        if err = os.MkdirAll(ipamDirURL, 0644); err != nil {
            return
        }
    }
    subnetsJSON, err := json.Marshal(ipam)
    if err != nil {
        return
    }
    ipamFile, err := os.OpenFile(ipam.SubnetAllocatorURL,
        os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        return
    }
    defer ipamFile.Close()
    if _, err = ipamFile.Write(subnetsJSON); err != nil {
        return
    }
    return nil
}
