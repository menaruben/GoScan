package main

import (
	"fmt"

	GoScan "github.com/menaruben/GoScan"
)

// "GoScan"
// "fmt"

func main() {
	// scan ports 20 to 30
	port_range := [2]int{20, 30}

	// scan each port concurrently
	result, runtime := GoScan.ScanHostFast("localhost", port_range)
	fmt.Printf("Port scanning finished in %f seconds\n", runtime.Seconds())
	fmt.Println(result)
}
