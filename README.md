# GoScan
GoScan is a blanzingly fast port scanner written in Go.

## Examples
### scan single port(s)
In order to scan a single port(s) you can use the ```ScanPort``` function:
```go
package main

import (
	"GoScan"
	"fmt"
)

func main() {
	sshResult := GoScan.ScanPort("localhost", 22)
	httpResult := GoScan.ScanPort("localhost", 80)

	fmt.Println(sshResult.Port, sshResult.State)
	fmt.Println(httpResult.Port, httpResult.State)
}
```
Now we can look at the terminal output in order to look at the results:
```
22 true
80 false
```
This means that port 22 is open and port 80 is closed.

### scan port range
In order to scan a port range concurrently you will need to use the ```ScanHostFast``` function. If you don't want the concurrent aspect of ```ScanHostFast``` you can replace it with the ```ScanHost``` function:

```go
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
```
After running the code above your terminal output should look like this:
```
Port scanning finished in 2.498907 seconds
+------+-------+--------------+
| PORT | STATE |   SERVICE    |
+------+-------+--------------+
|   22 | open  | ssh          |
|  445 | open  | microsoft-ds |
|  135 | open  | rpc          |
|  912 | open  |              |
|  902 | open  |              |
+------+-------+--------------+
```
