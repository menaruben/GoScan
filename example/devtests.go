package main

import (
	"GoScan"
	"fmt"
)

func main() {
	port_range := [2]int{1, 30}

	// ScanHostFast concurrently scans all ports of a host
	result, runtime := GoScan.ScanHost("localhost", port_range, 2)
	fmt.Printf("Port scanning finished in %f seconds\n", runtime.Seconds())

	GoScan.ResultOutput(result) // prints out result table to terminal
}
