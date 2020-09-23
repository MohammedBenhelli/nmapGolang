package scanner

import (
	"../arg"
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"net"
	"os"
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

func ScanPort(protocol, hostname string, port int, results *[]ScanResult, wg *sync.WaitGroup) {
	result := ScanResult{Port: strconv.Itoa(port), Type: protocol}
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, 200*time.Millisecond)
	defer wg.Done()
	if err != nil {
		result.State = "Closed"
		return
	}
	defer conn.Close()
	result.State = "Open"
	result.Service = GetService(strconv.Itoa(port), protocol)
	*results = append(*results, result)
}

func GetService(port string, protocol string) string {
	csvFile, err := os.Open("scanner/service-names-port-numbers.csv")
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			return "No Service Name found"
		} else if err != nil {
			log.Fatal(err)
		} else if line[1] == port && line[2] == protocol {
			return line[0]
		} else {
			portCsv, _ := strconv.Atoi(line[1])
			portArg, _ := strconv.Atoi(port)
			if portCsv > portArg {
				return "No Service Name found"
			}
		}
	}
}

func InitialScan(arguments arg.FilterArg, results *[]ScanResult) {
	srcIp, sPort := localIPPort(net.ParseIP(arguments.Ip()))
	if begin, err := strconv.Atoi(arguments.Ports()[0]); err == nil {
		if end, err := strconv.Atoi(arguments.Ports()[1]); err == nil {
			var wg sync.WaitGroup
			for i := begin; i <= end; i++ {
				if end-i <= arguments.Speedup() {
					wg.Add(end - i + 1)
					for j := i; j <= end; j++ {
						if arguments.Scan() == "syn" {
							go SynScan(arguments.Ip(), j, results, &wg, srcIp, sPort)
						} else {
							go ScanPort(arguments.Scan(), arguments.Ip(), j, results, &wg)
						}
					}
					i = end + 1
					wg.Wait()
				} else {
					wg.Add(arguments.Speedup())
					for j := i; j < i+arguments.Speedup(); j++ {
						if arguments.Scan() == "syn" {
							go SynScan(arguments.Ip(), j, results, &wg, srcIp, sPort)
						} else {
							go ScanPort(arguments.Scan(), arguments.Ip(), j, results, &wg)
						}
					}
					i += arguments.Speedup()
					wg.Wait()
				}
			}
		}
		printResult(*results)
	}
}
