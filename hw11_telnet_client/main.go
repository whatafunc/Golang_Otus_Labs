package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// Define and parse --timeout flag
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout (e.g. 10s, 1m)")
	flag.Parse()

	// After flags, expect exactly two positional args: host and port
	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--timeout=duration] host port\n", os.Args[0])
		os.Exit(1)
	}
	host := args[0]
	port := args[1]
	address := host + ":" + port

	// Create client with stdin/stdout and timeout
	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	// Connect to server
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	// Run Receive in a goroutine (reads from server, writes to stdout)
	go func() {
		if err := client.Receive(); err != nil {
			log.Printf("Receive error: %v", err)
			os.Exit(1)
		}
	}()

	// Run Send in main goroutine (reads from stdin, writes to server)
	if err := client.Send(); err != nil {
		log.Printf("Send error: %v", err)
	}
}
