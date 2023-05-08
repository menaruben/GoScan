package GoScan_test

import (
	"fmt"
	GoScan "github.com/menaruben/GoScan"
	"time"
)

func ExampleValidateIpv4() {
	validCheck, err := GoScan.ValidateIpv4("192.168.100.29")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(validCheck)
	// Output:
	// true
}

func ExampleGetSubnetMask() {
	subnetMask := GoScan.GetSubnetMask(25)
	fmt.Println(subnetMask)
	// Output:
	// 255.255.255.128
}

func ExampleScanPort() {
	sshResult, err := GoScan.ScanPort("localhost", 22, 12*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sshResult.Port, sshResult.State)
	// Output:
	// 22 true
}

func ExampleScanHost() {
	// scan ports 20 to 30
	portRange := [2]int{20, 30}

	// scan each port with 2 seconds interval
	result, err := GoScan.ScanHost("localhost", portRange, 2*time.Second, 12*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
	// Port scanning finished in 45.630600 seconds
	// [{22 true}]
}

func ExampleScanHostFast() {
	// scan ports 20 to 30
	portRange := [2]int{20, 30}

	// scan ports concurrently
	result := GoScan.ScanHostFast("localhost", portRange, 12*time.Second)
	fmt.Println(result)
	// Output:
	// [{22 true}]
}

func ExampleGetService() {
	service := GoScan.GetService(22)
	fmt.Println(service)
	// Output:
	// SSH (Secure Shell)
}

func ExampleResultOutput() {
	// scan ports 20 to 30
	portRange := [2]int{20, 30}

	// scan ports concurrently
	result := GoScan.ScanHostFast("localhost", portRange, 12*time.Second)
	GoScan.ResultOutput(result)
	// Output:
	// +------+-------+--------------------+
	// | PORT | STATE |      SERVICE       |
	// +------+-------+--------------------+
	// |   22 | open  | SSH (Secure Shell) |
	// +------+-------+--------------------+
}

func ExampleIsIPReachable() {
	validCheck, err := GoScan.IsIPReachable("142.250.203.100", 12*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(validCheck)
	// Output:
	// true
}

func ExampleScanNetwork() {
	myNetwork, err := GoScan.ScanNetwork("192.168.1.0/24", 12*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(myNetwork)
	// Output:
	// {192.168.1.0 255.255.255.0 24 [192.168.1.19 192.168.1.4 192.168.1.101]}
}

func ExampleScanNetHosts() {
	portRange := [2]int{20, 30}
	myNetwork, _ := GoScan.ScanNetwork("192.168.1.0/24", 12*time.Second)
	myResult, err := GoScan.ScanNetHosts(myNetwork, portRange, 0*time.Second, 12*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(myResult)
	// Output:
	// [[192.168.1.74 []] [192.168.1.19 [{22 true}]]
}

func ExampleScanNetHostsFast() {
	portRange := [2]int{20, 30}
	myNetwork, err := GoScan.ScanNetwork("192.168.1.0/24", 12*time.Second)
	if err != nil {
		fmt.Println(err)
	}

	myResult := GoScan.ScanNetHostsFast(myNetwork, portRange, 12*time.Second)

	fmt.Println(myResult)
	// Output:
	// [[192.168.1.19 [{22 true}]] [192.168.1.74 []]
}
