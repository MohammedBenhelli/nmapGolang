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
	//begin, _ := strconv.Atoi(arguments.Ports()[0])
	//end, _ := strconv.Atoi(arguments.Ports()[1])
	scanner.InitialScan(arguments, &results)
	for i := 0; i < len(results); i++ {
		if results[i].State == "Open" {
			fmt.Printf("%+v\n", results[i])
		}
	}
}
