// Package GoScan is a blanzingly fast network/port scanner written in Go
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
func ValidateIpv4(ipaddr string) (bool, error) {
	octets := strings.Split(ipaddr, ".")

	if len(octets) != 4 {
		return false, nil
	}

	for i := 0; i < 4; i++ {
		num, err := strconv.Atoi(octets[i])
		if err != nil {
			return false, err
		}

		if num < 0 || num > 255 {
			return false, nil
		}
	}

	return true, nil
}

// GetSubnetMask returns the subnet mask in the x.x.x.x format and takes the subnet suffix as input.
func GetSubnetMask(suffix int) string {
	if suffix < 0 || suffix > 32 {
		return ""
	}

	remainder := suffix % 8
	lastOctet := 256 - int(math.Pow(2, float64(8-remainder)))
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
func ScanPort(hostname string, port int, timeout time.Duration) (ScanResult, error) {
	result := ScanResult{Port: port}
	address := hostname + ":" + strconv.Itoa(port)

	conn, err := net.DialTimeout("tcp", address, timeout)

	if err != nil {
		result.State = false
		return result, err
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	result.State = true
	return result, nil
}

// ScanHost scans all ports inside the portRange argument and returns all open ports.
func ScanHost(hostname string, portRange [2]int, scanInterval time.Duration, timeout time.Duration) ([]ScanResult, error) {
	var result []ScanResult
	start := portRange[0]
	end := portRange[1]

	for i := start; i <= end; i++ {
		portResult, err := ScanPort(hostname, i, timeout)
		if err != nil {
			return result, err
		}
		if portResult.State {
			result = append(result, portResult)
		}
		time.Sleep(scanInterval)
	}

	return result, nil
}

// ScanHostFast scans all ports inside the portRange argument concurrently and returns all open ports.
func ScanHostFast(hostname string, portRange [2]int, timeout time.Duration) []ScanResult {
	var wg sync.WaitGroup

	// create a channel to receive the scan results
	results := make(chan ScanResult)

	// launch goroutine for each port in the range
	for port := portRange[0]; port <= portRange[1]; port++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			portResult, _ := ScanPort(hostname, port, timeout)
			results <- portResult
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

	return openPorts
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

	// append results to table
	for _, res := range results {
		state := "closed"
		if res.State {
			state = "open"
		}
		table.Append([]string{strconv.Itoa(res.Port), state, GetService(res.Port)})
	}

	table.Render()
}

// NetworkInfo contains information about a network.
type NetworkInfo struct {
	NetworkIP    string
	SubnetMask   string
	SubnetSuffix int
	Hosts        []string
}

// IsIPReachable returns if an IP address is reachable
func IsIPReachable(ipAddr string, timeout time.Duration) (bool, error) {
	conn, err := net.DialTimeout("tcp", ipAddr+":80", timeout)
	if err != nil {
		return false, err
	}

	err = conn.Close()
	if err != nil {
		return false, err
	}
	return true, nil
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
func ScanNetwork(netaddr string, timeout time.Duration) (NetworkInfo, error) {
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
		return network, err
	}

	hostChannel := make(chan string)

	// wait for all goroutines to complete
	var wg sync.WaitGroup
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incrementHostIp(ip) {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			ipReachable, _ := IsIPReachable(ip, timeout)
			if ipReachable {
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

	return network, nil
}

func getPortsFromResults(scanResults []ScanResult) []int {
	var portNumbers []int
	for _, result := range scanResults {
		portNumbers = append(portNumbers, result.Port)
	}

	return portNumbers
}

// ScanNetHosts returns an array of an array of scan results of all hosts inside a network
func ScanNetHosts(network NetworkInfo, portRange [2]int, scanInterval time.Duration, timeout time.Duration) ([][2]any, []error) {
	var scanResults [][2]any
	var errors []error

	for i := 0; i < len(network.Hosts); i++ {
		result, err := ScanHost(network.Hosts[i], portRange, scanInterval, timeout)
		if err != nil {
			errors = append(errors, err)
		} else {
			resultPort := getPortsFromResults(result)
			resultHost := [2]any{network.Hosts[i], resultPort}
			scanResults = append(scanResults, resultHost)
		}
		time.Sleep(scanInterval)
	}

	return scanResults, errors
}

// ScanNetHostsFast returns an array of an array of scan results of all hosts inside a network
func ScanNetHostsFast(network NetworkInfo, portRange [2]int, timeout time.Duration) [][2]any {
	var scanResults [][2]any
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, host := range network.Hosts {
		wg.Add(1)

		go func(hostname string) {
			defer wg.Done()

			result := ScanHostFast(hostname, portRange, timeout)

			// lock current variable to the go routine
			mu.Lock()
			resultPort := getPortsFromResults(result)
			resultHost := [2]any{hostname, resultPort}
			scanResults = append(scanResults, resultHost)
			mu.Unlock()
		}(host)
	}

	wg.Wait()
	return scanResults
}
