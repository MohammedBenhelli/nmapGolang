package arg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Arg() FilterArg {
	lastCommand := os.Args[1]
	filter := FilterArg{true, []string{"test"}, "test", 1024, "tcp", false}
	if os.Args[1] == "--help" {
		fmt.Println("help")
		return filter
	} else {
		filter.SetHelp(false)
	}
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i][0] == 45 && os.Args[i][1] == 45 {
			lastCommand = os.Args[i]
			i++
		}
		switch lastCommand {
		case "--ports":
			filter.SetPorts(filterPorts(os.Args[i]))
		case "--ip":
			filter.SetIp(os.Args[i])
		case "--file":
			filter.SetIp(os.Args[i])
			filter.SetFile(true)
		case "--speedup":
			tmp, _ := strconv.Atoi(os.Args[i])
			filter.SetSpeedup(tmp)
		case "--scan":
			filter.SetScan(strings.ToLower(os.Args[i]))
		}
	}
	if !filter.file {
		fmt.Printf("Target Ip-Address : %s\n", filter.Ip())
	} else {
		fmt.Printf("Target File for Ip-Addresses : %s\n", filter.Ip())
	}
	fmt.Printf("No of Threads : %d\n", filter.Speedup())
	fmt.Printf("Scan to be performed : %s\n", filter.Scan())
	return filter
}

func filterPorts(str string) []string {
	num := 0
	if strings.Contains(str, "-") {
		ret := strings.Split(str, "-")
		for i := 0; i < len(ret); i += 2 {
			if i+2 < len(ret) || i == 0 {
				sum1, _ := strconv.Atoi(ret[i+1])
				sum2, _ := strconv.Atoi(ret[i])
				num += sum1 - sum2
			} else {
				num++
			}
		}
		fmt.Println("No of Ports to scan : " + strconv.Itoa(num))
		return ret
	} else {
		fmt.Println("No of Ports to scan : 1")
		return []string{str}
	}
}
