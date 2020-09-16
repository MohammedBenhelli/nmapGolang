package scanner

import (
	"../arg"
	"net"
	"strconv"
	"sync"
	"time"
)

type ScanResult struct {
	Port    string
	State   string
	Service string
	Type    string
}

func (s *ScanResult) test() {
	println(s.Port)
}

func ScanPort(protocol, hostname string, port int, results *[]ScanResult, wg *sync.WaitGroup) {
	result := ScanResult{Port: strconv.Itoa(port), Type: protocol}
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, 60*time.Second)

	defer wg.Done()

	if err != nil {
		result.State = "Closed"
		*results = append(*results, result)
		return
	}
	defer conn.Close()
	result.State = "Open"
	*results = append(*results, result)
}

//func GetService(port int, protocol string) string {
//	return netdb.GetServByPort(port, netdb.GetProtoByName(protocol)).Name
//}

func InitialScan(arguments arg.FilterArg, results *[]ScanResult) {
	if begin, err := strconv.Atoi(arguments.Ports()[0]); err == nil {
		if end, err := strconv.Atoi(arguments.Ports()[1]); err == nil {
			var wg sync.WaitGroup
			for i := begin; i <= end; i++ {
				if end - i <= arguments.Speedup() {
					wg.Add(end - i + 1)
					for j := i; j <= end; j++ {
						go ScanPort("tcp", arguments.Ip(), j, results, &wg)
					}
					i = end + 1
					wg.Wait()
				} else {
					wg.Add(arguments.Speedup())
					for j := i; j < i + arguments.Speedup(); j++ {
						go ScanPort("tcp", arguments.Ip(), j, results, &wg)
					}
					i += arguments.Speedup()
					wg.Wait()
				}
			}
		}
	}

}
