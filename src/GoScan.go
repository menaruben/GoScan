package main

// package GoScan

import (
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
)

var portServices = map[int]string{
	21:   "FTP (File Transfer Protocol)",
	22:   "SSH (Secure Shell)",
	23:   "Telnet",
	25:   "SMTP (Simple Mail Transfer Protocol)",
	53:   "DNS (Domain Name System)",
	67:   "BOOTP / DHCP",
	68:   "BOOTP / DHCP",
	80:   "HTTP (Hypertext Transfer Protocol)",
	110:  "POP3 (Post Office Protocol version 3)",
	119:  "NNTP (Network News Transfer Protocol)",
	123:  "NTP (Network Time Protocol)",
	143:  "IMAP (Internet Message Access Protocol)",
	161:  "SNMP (Simple Network Management Protocol)",
	194:  "IRC (Internet Relay Chat)",
	220:  "IMAP version 3",
	445:  "Microsoft-DS (Directory Services)",
	443:  "HTTPS (HTTP Secure)",
	465:  "SMTPS (SMTP Secure)",
	587:  "SMTP (Mail Submission Agent)",
	749:  "Kerberos administration",
	751:  "Kerberos authentication",
	752:  "Kerberos password (kpasswd) server",
	902:  "VMware ESXi",
	903:  "VMware ESXi",
	993:  "IMAPS (IMAP Secure)",
	995:  "POP3S (POP3 Secure)",
	1433: "Microsoft SQL Server",
	3306: "MySQL",
	5432: "PostgreSQL",
	8080: "HTTP (Alternative Port)",
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
func ScanHost(hostname string, port_range [2]int, scan_interval int) ([]ScanResult, time.Duration) {
	start_time := time.Now()
	var result []ScanResult
	start := port_range[0]
	end := port_range[1]

	for i := start; i <= end; i++ {
		port_result := ScanPort(hostname, i)
		if port_result.State {
			result = append(result, port_result)
		}
		time.Sleep(time.Duration(scan_interval) * time.Second)
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
	_, ok := portServices[port]
	if ok {
		return portServices[port]
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
