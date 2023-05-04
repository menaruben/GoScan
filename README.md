# GoScan
GoScan is a blanzingly fast port scanner written in Go.

## Installation
Move to your folder containing your go.mod and type:
```
go get github.com/menaruben/GoScan
```

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
