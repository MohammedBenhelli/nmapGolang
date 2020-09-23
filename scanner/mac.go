package scanner

import (
	"fmt"
	"log"
	"net"
)

type MacResult struct {
	MAC string
	IP  net.IP
	Mask string
}

func GetMacAddr() ([]MacResult, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var results []MacResult
	for _, ifa := range ifaces {
		address, _ := ifa.Addrs()
		for _, addr := range address {
			switch v := addr.(type) {
				case *net.IPNet:
					a := MacResult{ifa.HardwareAddr.String(), v.IP, v.IP.DefaultMask().String()}
					if a.MAC != "" {
						results = append(results, a)
					}
				case *net.IPAddr:
					a := MacResult{ifa.HardwareAddr.String(), v.IP, v.IP.DefaultMask().String()}
					if a.MAC != "" {
						results = append(results, a)
					}
			}
		}

	}
	return results, nil
}

func LocalIp()  {
	as, err := GetMacAddr()
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range as {
		fmt.Println(a)
	}
}
