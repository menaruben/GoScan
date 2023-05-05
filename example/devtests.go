package main

import (
	"fmt"

	GoScan "github.com/menaruben/GoScan"
)

func main() {
	var myNetwork NetworkInfo = GoScan.ScanNetwork("192.168.1.0/24", 0)
	fmt.Println(myNetwork)
}
