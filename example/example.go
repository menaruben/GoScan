package main

import (
	"GoScan"
	"fmt"
)

func main() {
	// port_range := [2]int{1, 1024}

	// result, runtime := GoScan.ScanHostFast("localhost", port_range)
	// fmt.Printf("finished in %f seconds\n", runtime.Seconds())
	// GoScan.ResultOutput(result)

	sshResult := GoScan.ScanPort("localhost", 22)
	httpResult := GoScan.ScanPort("localhost", 80)

	fmt.Println(sshResult.Port, sshResult.State)
	fmt.Println(httpResult.Port, httpResult.State)
}
