# GoScan
GoScan is a blanzingly fast port scanner written in Go.

## Installation
> ⚠️ **_NOTE:_**  This module is under heavy development. If any errors occur please file an [issue](https://github.com/menaruben/GoScan/issues).
Move to your folder containing your go.mod and type:
```
go get github.com/menaruben/GoScan
```

## Examples
### ValidateIpv4
We always need to pass in ip addresses or hostnames. When we parse out the IP addresses from somewhere we can try to vaildate the parsed string with ```ValidateIpv4``` so that we can avoid errors when passing in wrong ip addresses. Here is an example on how to use it:
```go
package main

import (
	"fmt"
	"github.com/menaruben/GoScan"
)

func main() {
	validCheck := GoScan.ValidateIpv4("192.168.100.29")
	fmt.Println(validCheck)
}
```
Output:
```
true
```
This tells us that the IPv4 address is in the correct format.

### GetSubnetMask
When scanning a network it is important to get certain information about the network such as the network address and the subnet mask in order to know all the possible IPs a network can have. The ```GetSubnetMask``` returns the subnet mask for the given suffix. Here is an example on how to use the ```GetSubnetMask``` function:
```go
package main

import (
	"fmt"
	"github.com/menaruben/GoScan"
)

func main() {
	subnetMask := GoScan.GetSubnetMask(25)
	fmt.Println(subnetMask)
}
```
Output:
```
255.255.255.128
```

### ScanPort
In order to scan a single port(s) you can use the ```ScanPort``` function:
```go
package main

import (
	"github.com/menaruben/GoScan"
	"fmt"
	"time"
)

func main() {
	timeout := 12*time.Second
	sshResult := GoScan.ScanPort("localhost", 22, timeout)
	httpResult := GoScan.ScanPort("localhost", 80, timeout)

	fmt.Println(sshResult.Port, sshResult.State)
	fmt.Println(httpResult.Port, httpResult.State)
}
```
Output:
```
22 true
80 false
```
This means that port 22 is open and port 80 is closed.

### ScanHostFast
In order to scan a port range concurrently you will need to use the ```ScanHostFast``` function. If you don't want the concurrent aspect of ```ScanHostFast``` you can replace it with the ```ScanHost``` function:

```go
package main

import (
	"github.com/menaruben/GoScan"
	"fmt"
	"time"
)

func main() {
    port_range := [2]int{1, 1024}
	timeout := 12*time.Second

    // ScanHostFast concurrently scans all ports of a host
    result, runtime := GoScan.ScanHostFast("localhost", port_range, timeout)
    fmt.Printf("Port scanning finished in %f seconds\n", runtime.Seconds())

    GoScan.ResultOutput(result) // prints out result table to terminal
}
```
Output:
```
finished in 2.490884 seconds
+------+-------+--------------------------------+
| PORT | STATE |            SERVICE             |
+------+-------+--------------------------------+
|  445 | open  | Microsoft-DS (Directory        |
|      |       | Services)                      |
|   22 | open  | SSH (Secure Shell)             |
|  135 | open  |                                |
|  912 | open  |                                |
|  902 | open  | VMware ESXi                    |
+------+-------+--------------------------------+
```

### ScanHost
In order to less obvious when scanning a network we can use the ```ScanHost``` function that lets us define the time interval between the scanning:
```go
package main

import (
	"github.com/menaruben/GoScan"
	"fmt"
	"time"
)

func main() {
	// scan ports 20 to 30
	timeout := 12*time.Second
	port_range := [2]int{20, 30}

	// scan each port with 2 seconds interval
	result, runtime := GoScan.ScanHost("localhost", port_range, 2, timeout)
	fmt.Printf("Port scanning finished in %f seconds\n", runtime.Seconds())
	fmt.Println(result)
}
```
Output:
```
Port scanning finished in 45.630600 seconds
[{22 true}]
```

### ExampleGetService
This function returns the service that is mapped to the port given as the argument:
```go
package main

import (
	"github.com/menaruben/GoScan"
	"fmt"
	"time"
)

func main() {
	service := GoScan.GetService(22)
	fmt.Println(service)
}
```
Output:
```
SSH (Secure Shell)
```

### ResultOutput
The ```ResultOutput``` function prints out the ```[]ScanResult``` given as an argument. Here is a small example:
```go
package main

import (
	"github.com/menaruben/GoScan"
	"fmt"
	"time"
)

func main() {
	port_range := [2]int{20, 30}

	// ScanHost scans all ports of a host with interval of 2 seconds between scans
	result, runtime := ScanHost("localhost", port_range, 2, 12*time.Second)
	fmt.Printf("Port scanning finished in %f seconds\n", runtime.Seconds())

	ResultOutput(result) // prints out result table to terminal
}
```
Output:
```
Port scanning finished in 45.561080 seconds
+------+-------+--------------------+
| PORT | STATE |      SERVICE       |
+------+-------+--------------------+
|   22 | open  | SSH (Secure Shell) |
+------+-------+--------------------+
```

### IsIPReachable
This function checks wether or not an IP is reachable or not and returns the boolean.
```go
package main

import (
	"github.com/menaruben/GoScan"
	"fmt"
)

func main() {
	var validCheck bool = IsIPReachable("142.250.203.100", 12*time.Second)
	fmt.Println(validCheck)
}
```
Output:
```
true
```

### ScanNetwork
This function scans a network and returns a variable of type ```NetworkInfo``` which is defined like this:
```go
type NetworkInfo struct {
	NetworkIP    string
	SubnetMask   string
	SubnetSuffix int
	Hosts        []string
}
```

In order to scan a network we need to use the following code:
```go
package main

import (
	"github.com/menaruben/GoScan"
	"fmt"
)

func main() {
	myNetwork := GoScan.ScanNetwork("192.168.1.0/24", 0, 12*time.Second)
	fmt.Println(myNetwork)
}
```
Output:
```
{192.168.1.0 255.255.255.0 24 [192.168.1.19 192.168.1.4 192.168.1.101]}
```

## License
```
MIT License

Copyright (c) 2023 Rubén Mena

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
