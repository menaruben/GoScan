package GoScan

import (
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
)

var defaultServices = map[int]string{
	1:   "tcpmux",
	5:   "rje",
	7:   "echo",
	18:  "msp",
	20:  "ftp-data",
	21:  "ftp",
	22:  "ssh",
	23:  "telnet",
	25:  "smtp",
	29:  "msg-icp",
	37:  "time",
	42:  "nameserver",
	43:  "whois",
	49:  "tacacs",
	53:  "domain",
	69:  "tftp",
	70:  "gopher",
	79:  "finger",
	80:  "http",
	103: "x400",
	108: "sna",
	109: "pop2",
	110: "pop3",
	115: "sftp",
	118: "sqlserv",
	119: "nntp",
	135: "rpc",
	137: "netbios-ns",
	139: "netbios-ssn",
	143: "imap",
	150: "sql-net",
	156: "sqlsrv",
	161: "snmp",
	179: "bgp",
	190: "gacp",
	194: "irc",
	197: "dls",
	389: "ldap",
	396: "netware-ip",
	443: "https",
	444: "snpp",
	445: "microsoft-ds",
	458: "appleqtc",
	546: "dhcp-client",
	547: "dhcp-server",
	563: "snews",
	569: "msrpc",
}

// The ScanResult struct contains the port and its state (wether it is open or closed)
type ScanResult struct {
	Port  int
	State bool
}

// ScanPort scans a single port
func ScanPort(hostname string, port int) ScanResult {
	result := ScanResult{Port: port}
	address := hostname + ":" + strconv.Itoa(port)

	conn, err := net.DialTimeout("tcp", address, 12*time.Second)

	if err != nil {
		result.State = false
		return result
	}

	defer conn.Close()

	result.State = true
	return result
}

// ScanHost scans all ports inside the port_range argument and returns all open ports.
func ScanHost(hostname string, port_range [2]int) ([]ScanResult, time.Duration) {
	start_time := time.Now()
	var result []ScanResult
	start := port_range[0]
	end := port_range[1]

	for i := start; i <= end; i++ {
		result = append(result, ScanPort(hostname, i))
	}

	return result, time.Since(start_time)
}

// ScanHostFast scans all ports inside the port_range argument concurrently and returns all open ports.
func ScanHostFast(hostname string, port_range [2]int) ([]ScanResult, time.Duration) {
	start_time := time.Now()
	var wg sync.WaitGroup

	// create a channel to receive the scan results
	results := make(chan ScanResult)

	// launch goroutine for each port in the range
	for port := port_range[0]; port <= port_range[1]; port++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			results <- ScanPort(hostname, port)
		}(port)
	}

	// wait for all tasks to complete
	go func() {
		wg.Wait()
		close(results) // close results channel
	}()

	// receive results from channel and add them to openPorts
	var openPorts []ScanResult
	for result := range results {
		if result.State {
			openPorts = append(openPorts, result)
		}
	}

	return openPorts, time.Since(start_time)
}

// The getService returns the service for the given port.
func getService(port int) string {
	_, ok := defaultServices[port]
	if ok {
		return defaultServices[port]
	} else {
		return ""
	}
}

// ResultOutput prints the output to the terminal as a table.
func ResultOutput(results []ScanResult) {
	// initialize table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Port", "State", "Service"})

	// append restults to table
	for _, res := range results {
		state := "closed"
		if res.State {
			state = "open"
		}
		table.Append([]string{strconv.Itoa(res.Port), state, getService(res.Port)})
	}

	table.Render()
}
