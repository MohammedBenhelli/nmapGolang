package scanner

import (
	"../arg"
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
)

func FileScan(arguments arg.FilterArg, results *[]ScanResult) {
	ipList := fileToTab(arguments.Ip())
	for _, ip := range ipList {
		arguments.SetIp(ip)
		color.Green("Ip tested %s\n", ip)
		InitialScan(arguments, results)
		printResult(*results)
		*results = []ScanResult{}
	}
}

func fileToTab(path string) []string {
	var result []string
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return []string{}
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			result = append(result, scanner.Text())
		}
	}
	defer file.Close()
	return result
}

func printResult(results []ScanResult) {
	for i := 0; i < len(results); i++ {
		fmt.Printf("%+v\n", results[i])
	}
}
