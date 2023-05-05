package GoScan_test

import (
	"fmt"
	"github.com/menaruben/GoScan"
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
	sshResult := GoScan.ScanPort("localhost", 22)
	fmt.Println(sshResult.Port, sshResult.State)
	// Output:
	// 22 true
}

func ExampleScanHost() {
	// scan ports 20 to 30
	port_range := [2]int{20, 30}

	// scan each port with 2 seconds interval
	result, runtime := GoScan.ScanHost("localhost", port_range, 2)
	fmt.Printf("Port scanning finished in %f seconds\n", runtime.Seconds())
	fmt.Println(result)
	// Output:
	// Port scanning finished in 45.630600 seconds
	// [{22 true}]
}

func ExampleScanHostFast() {
	// scan ports 20 to 30
	port_range := [2]int{20, 30}

	// scan ports concurrently
	result, runtime := GoScan.ScanHostFast("localhost", port_range)
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
	port_range := [2]int{20, 30}

	// scan ports concurrently
	result, _ := GoScan.ScanHostFast("localhost", port_range)

	GoScan.ResultOutput(result)
	// Output:
	// +------+-------+--------------------+
	// | PORT | STATE |      SERVICE       |
	// +------+-------+--------------------+
	// |   22 | open  | SSH (Secure Shell) |
	// +------+-------+--------------------+
}
