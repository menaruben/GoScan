package main

import (
	"fmt"
	"time"

	GoScan "github.com/menaruben/GoScan"
)

func main() {
	portRange := [2]int{20, 30}
	myNetwork := GoScan.ScanNetwork("192.168.1.0/24", 12*time.Second)
	myResult := GoScan.ScanNetHostsFast(myNetwork, portRange, 12*time.Second)
	fmt.Println(myNetwork.Hosts)
	fmt.Println(myResult)

	// fmt.Println(myNetwork)
	// // scan port range
	// port_range := [2]int{1, 1024}

	// result, runtime := GoScan.ScanHostFast("localhost", port_range)
	// fmt.Printf("finished in %f seconds\n", runtime.Seconds())
	// GoScan.ResultOutput(result)

	// // scan single port(s)
	// sshResult := GoScan.ScanPort("localhost", 22)
	// httpResult := GoScan.ScanPort("localhost", 80)

	// fmt.Println(sshResult.Port, sshResult.State)
	// fmt.Println(httpResult.Port, httpResult.State)
}
