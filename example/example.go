package main

import (
	GoScan "github.com/menaruben/GoScan"
	"fmt"
)

func main() {
	// scan port range
	port_range := [2]int{1, 1024}

	result, runtime := GoScan.ScanHostFast("localhost", port_range)
	fmt.Printf("finished in %f seconds\n", runtime.Seconds())
	GoScan.ResultOutput(result)

	// scan single port(s)
	sshResult := GoScan.ScanPort("localhost", 22)
	httpResult := GoScan.ScanPort("localhost", 80)

	fmt.Println(sshResult.Port, sshResult.State)
	fmt.Println(httpResult.Port, httpResult.State)
}
