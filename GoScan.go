// GoScan is a blanzingly fast network/port scanner written in Go.
package GoScan

import (
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
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

// ValidateIpv4 returns a bool and checks wether the input is a valid IPv4 address.
func ValidateIpv4(ipaddr string) bool {
	octets := strings.Split(ipaddr, ".")

	if len(octets) != 4 {
		return false
	}

	for i := 0; i < 4; i++ {
		num, err := strconv.Atoi(octets[i])
		if err != nil {
			return false
		}

		if num < 0 || num > 255 {
			return false
		}
	}

	return true
}

// GetSubnetMask returns the subnet mask in the x.x.x.x format and takes the subnet suffix as input.
func GetSubnetMask(suffix int) string {
	if suffix < 0 || suffix > 32 {
		return ""
	}

	remainder := suffix % 8
	lastOctet := 256 - int(math.Pow(2, float64((8-remainder))))
	previousOctetsNum := (suffix - remainder) / 8

	var octets []string
	for i := 0; i < previousOctetsNum; i++ {
		octets = append(octets, "255")
	}

	octets = append(octets, strconv.Itoa(lastOctet))
	var octetsMissing int = 4 - len(octets)

	for i := 0; i < octetsMissing; i++ {
		octets = append(octets, "0")
	}

	subnetMask := strings.Join(octets, ".")

	return subnetMask
}

// ScanResult struct contains the port and its state as a boolean (open = true, closed = false).
type ScanResult struct {
	Port  int
	State bool
}

// ScanPort scans a single port.
func ScanPort(hostname string, port int, timeout time.Duration) ScanResult {
	result := ScanResult{Port: port}
	address := hostname + ":" + strconv.Itoa(port)

	conn, err := net.DialTimeout("tcp", address, timeout)

	if err != nil {
		result.State = false
		return result
	}

	defer conn.Close()

	result.State = true
	return result
}

// ScanHost scans all ports inside the port_range argument and returns all open ports.
func ScanHost(hostname string, port_range [2]int, scan_interval int, timeout time.Duration) ([]ScanResult, time.Duration) {
	start_time := time.Now()
	var result []ScanResult
	start := port_range[0]
	end := port_range[1]

	for i := start; i <= end; i++ {
		port_result := ScanPort(hostname, i, timeout)
		if port_result.State {
			result = append(result, port_result)
		}
		time.Sleep(time.Duration(scan_interval) * time.Second)
	}

	return result, time.Since(start_time)
}

// ScanHostFast scans all ports inside the port_range argument concurrently and returns all open ports.
func ScanHostFast(hostname string, port_range [2]int, timeout time.Duration) ([]ScanResult, time.Duration) {
	start_time := time.Now()
	var wg sync.WaitGroup

	// create a channel to receive the scan results
	results := make(chan ScanResult)

	// launch goroutine for each port in the range
	for port := port_range[0]; port <= port_range[1]; port++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			results <- ScanPort(hostname, port, timeout)
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

// The GetService returns the service for the given port.
func GetService(port int) string {
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
		table.Append([]string{strconv.Itoa(res.Port), state, GetService(res.Port)})
	}

	table.Render()
}

// Contains information about a network.
type NetworkInfo struct {
	NetworkIP    string
	SubnetMask   string
	SubnetSuffix int
	Hosts        []string
}

// returns wether or not an IP address is reachacble
func IsIPReachable(ipAddr string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", ipAddr+":80", timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func incrementHostIp(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// ScanNetwork returns NetworkInfo about a specific network. It contains the network IP, subnet mask/suffix and all hosts.
func ScanNetwork(netaddr string, timeout time.Duration) NetworkInfo {
	var network NetworkInfo
	netaddrFields := strings.Split(netaddr, "/")

	network.NetworkIP = netaddrFields[0]
	SubnetSuffix, err := strconv.Atoi(netaddrFields[1])
	if err != nil {
		log.Fatal(err)
	}

	network.SubnetSuffix = SubnetSuffix
	network.SubnetMask = GetSubnetMask(network.SubnetSuffix)

	// get all hosts inside the network
	ip, ipNet, err := net.ParseCIDR(netaddr)
	if err != nil {
		log.Fatal(err)
	}

	hostChannel := make(chan string)

	// wait for all goroutines to complete
	var wg sync.WaitGroup
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incrementHostIp(ip) {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			if IsIPReachable(ip, timeout) {
				hostChannel <- ip
			}
		}(ip.String())
	}

	// Wait for all goroutines to complete before closing the channel
	go func() {
		wg.Wait()
		close(hostChannel)
	}()

	// Read from the channel until it's closed
	for host := range hostChannel {
		if host != "" {
			network.Hosts = append(network.Hosts, host)
		}
	}

	// remove network (first) and broadcast (last) address
	if len(network.Hosts) >= 3 {
		network.Hosts = network.Hosts[1 : len(network.Hosts)-1]
	}

	return network
}
