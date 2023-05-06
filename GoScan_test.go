package GoScan_test

import (
	"fmt"
	"github.com/menaruben/GoScan"
	"time"
)

func ExampleValidateIpv4() {
	validCheck := GoScan.ValidateIpv4("192.168.100.29")
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
	sshResult := GoScan.ScanPort("localhost", 22, 12*time.Second)
	fmt.Println(sshResult.Port, sshResult.State)
	// Output:
	// 22 true
}

func ExampleScanHost() {
	// scan ports 20 to 30
	portRange := [2]int{20, 30}

	// scan each port with 2 seconds interval
	result, runtime := GoScan.ScanHost("localhost", portRange, 2*time.Second, 12*time.Second)
	fmt.Printf("Port scanning finished in %f seconds\n", runtime.Seconds())
	fmt.Println(result)
	// Output:
	// Port scanning finished in 45.630600 seconds
	// [{22 true}]
}

func ExampleScanHostFast() {
	// scan ports 20 to 30
	portRange := [2]int{20, 30}

	// scan ports concurrently
	result, runtime := GoScan.ScanHostFast("localhost", portRange, 12*time.Second)
	fmt.Printf("Port scanning finished in %f seconds\n", runtime.Seconds())
	fmt.Println(result)
	// Output:
	// Port scanning finished in 2.352606 seconds
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
	result, _ := GoScan.ScanHostFast("localhost", portRange, 12*time.Second)

	GoScan.ResultOutput(result)
	// Output:
	// +------+-------+--------------------+
	// | PORT | STATE |      SERVICE       |
	// +------+-------+--------------------+
	// |   22 | open  | SSH (Secure Shell) |
	// +------+-------+--------------------+
}

func ExampleIsIPReachable() {
	var validCheck bool = GoScan.IsIPReachable("142.250.203.100", 12*time.Second)
	fmt.Println(validCheck)
	// Output:
	// true
}

func ExampleScanNetwork() {
	myNetwork := GoScan.ScanNetwork("192.168.1.0/24", 12*time.Second)
	fmt.Println(myNetwork)
	// Output:
	// {192.168.1.0 255.255.255.0 24 [192.168.1.19 192.168.1.4 192.168.1.101]}
}

func ExampleScanNetHosts() {
	portRange := [2]int{20, 30}
	myNetwork := GoScan.ScanNetwork("192.168.1.0/24", 12*time.Second)
	myResult := GoScan.ScanNetHosts(myNetwork, portRange, 0*time.Second, 12*time.Second)
	fmt.Println(myResult)
	// Output:
	// [[192.168.1.74 []] [192.168.1.19 [{22 true}]]
}

func ExampleScanNetHostsFast() {
	portRange := [2]int{20, 30}
	myNetwork := GoScan.ScanNetwork("192.168.1.0/24", 12*time.Second)
	myResult := GoScan.ScanNetHostsFast(myNetwork, portRange, 12*time.Second)
	fmt.Println(myResult)
	// Output:
	// [[192.168.1.19 [{22 true}]] [192.168.1.74 []]
}
