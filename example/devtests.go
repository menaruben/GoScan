package main

import (
	"fmt"
	"time"

	"github.com/menaruben/GoScan"
)

func main() {
	// var myNetwork NetworkInfo = GoScan.ScanNetwork("192.168.1.0/24", 0, 12*time.Second)
	// fmt.Println(myNetwork)
	myPorts := GoScan.ScanHostFast("localhost", [2]int{20, 30}, 9*time.Second)
	fmt.Println(myPorts)

	sshResult, err := GoScan.ScanPort("localhost", 22, 2*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sshResult)

}
