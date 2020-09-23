package main

import (
	"./arg"
	"./scanner"
	"fmt"
)

var (
	results []scanner.ScanResult
)

func main() {
	arguments := arg.Arg()
	fmt.Printf("%+v\n", arguments)
	if !arguments.File() {
		scanner.InitialScan(arguments, &results)
	} else {
		scanner.FileScan(arguments, &results)
	}
	//scanner.LocalIp()
	//ifs, _ := pcap.FindAllDevs()
	//for i := range ifs {
	//	fmt.Println(ifs[i])
	//}
}
