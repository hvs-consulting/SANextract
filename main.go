package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type result struct {
	Target string   `json:"target"`
	SANs   []string `json:"SANs"`
}

var workers int
var jsonOutput bool
var tlsTimeout time.Duration

func init() {
	flag.IntVar(&workers, "workers", 250, "Number of workers.")
	flag.BoolVar(&jsonOutput, "json", false, "Output JSON.")
	flag.DurationVar(&tlsTimeout, "timeout", 2500*time.Millisecond, "Connection timeout as duration, e.g. 2s or 800ms")
	flag.Parse()
}

func main() {
	printDisclaimer()
	// Start Workers
	scanQueue := make(chan func() result, 20)
	resultsQueue := make(chan result, 20)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(queue <-chan func() result, results chan<- result) {
			for queuedTask := range queue {
				results <- queuedTask()
			}
			wg.Done()
		}(scanQueue, resultsQueue)
	}

	// Fill input queue with tasks
	go func(scanQueue chan<- func() result) {
		// Read from stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			target := scanner.Text()
			// Filter input
			sanitized := sanitizeInput(target)
			// This queues the scan
			scanQueue <- getScanFunc(sanitized)
		}
		err := scanner.Err()
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Error on input processing: %s", err.Error()))
		}
		close(scanQueue)
	}(scanQueue)

	// Close results queue after all workers are finished
	go func(resultsQueue chan result) {
		wg.Wait()
		close(resultsQueue)
	}(resultsQueue)

	// Process results
	for result := range resultsQueue {
		processOutput(jsonOutput, result)
	}
}

func printDisclaimer() {
	fmt.Fprintln(os.Stderr, "SANextract - fetch TLS certificates from endpoints and extract Subject Alternative Names")
	fmt.Fprintln(os.Stderr, "By Michael Eder, HvS-Consulting AG")
	fmt.Fprintln(os.Stderr, "https://www.hvs-consulting.de/  https://twitter.com/michael_eder_")
	fmt.Fprintln(os.Stderr, "")
}

// This returns a closure that simply needs to be called by a worker to perform the check against a target
func getScanFunc(target string) func() result {
	return func() result {
		// Connect to the remote server
		conn, err := net.DialTimeout("tcp", target, tlsTimeout)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: Failed to connect: %s", target, err.Error()))
			return result{}
		}
		defer conn.Close()
		// For SNI, we need the host name. The input has already been sanitized, so we assume something like <host>:<port> here.
		// This is a best effort approach
		host := strings.Split(target, ":")[0]

		// Do the TLS handshake
		tlsConn := tls.Client(conn, &tls.Config{ServerName: host, InsecureSkipVerify: true}) // Skip the certificate verification as this is not of interest
		err = tlsConn.Handshake()
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: Failed to connect: %s", target, err.Error()))
			return result{}
		}
		cs := tlsConn.ConnectionState()
		r := result{
			Target: target,
			SANs:   []string{},
		}

		r.SANs = cs.PeerCertificates[0].DNSNames
		return r
	}
}

func processOutput(printJSON bool, r result) {
	if r.Target == "" {
		return
	}
	if printJSON {
		marshalled, err := json.Marshal(r)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Error marshalling result for %s", r.Target))
		}
		fmt.Println(string(marshalled))
	} else {
		fmt.Println(strings.Join(r.SANs, "\n"))
	}
}
