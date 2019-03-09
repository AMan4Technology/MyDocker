package datastructure

import (
	"errors"
	"fmt"
	"net"
	"testing"
)

func TestBitmap(t *testing.T) {
	ip1, subnet, _ := net.ParseCIDR("192.168.0.5/24")
	fmt.Println(ip1.String(), subnet.IP.String())
	subnetStr := subnet.String()
	one, size := subnet.Mask.Size()
	bm := NewBitmap(1 << uint(size-one))
	bm.Save(0)
	if bm.Full() {
		err := errors.New("subnet can't provide more ip")
		fmt.Printf("Allocate ip from %s failed. %v", subnetStr, err)
	}
	for i := 1; i < bm.Cap(); i++ {
		if !bm.Have(i) {
			val := i
			bm.Save(val)
			ip := make([]byte, 4)
			for j := len(ip) - 1; j >= 0; j-- {
				ip[j] = subnet.IP[j] + uint8(val)
				val >>= 8
			}
			fmt.Println(net.IP(ip).String())
			fmt.Println(bm.Have(i))
			break
		}
	}
	ip, _, _ := net.ParseCIDR("192.168.0.1/24")
	ip = ip.To4()
	var value int
	for i, length := 0, len(ip); i < length; i++ {
		value = value<<8 + int(ip[i]-subnet.IP[i])
	}
	fmt.Println(value)
	fmt.Println(bm.Have(value))
	bm.Remove(value)
	fmt.Println(bm.Have(value))
}
