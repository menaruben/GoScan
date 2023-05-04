package main

import (
	"GoScan"
	"fmt"
)

func main() {
	port_range := [2]int{1, 1024}

	// ScanHostFast concurrently scans all ports of a host
	result, runtime := GoScan.ScanHostFast("localhost", port_range)
	fmt.Printf("Port scanning finished in %f seconds\n", runtime.Seconds())

	GoScan.ResultOutput(result) // prints out result table to terminal
}
