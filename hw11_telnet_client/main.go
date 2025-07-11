package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--timeout=duration] host port\n", os.Args[0])
		os.Exit(1)
	}

	address := args[0] + ":" + args[1]

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	// Setup signal channel
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Run Send and Receive concurrently without context
	errCh := make(chan error, 2)

	go func() {
		errCh <- client.Receive()
	}()

	go func() {
		errCh <- client.Send()
	}()

	select {
	case sig := <-sigCh:
		log.Printf("Received signal %v, exiting", sig)
	case err := <-errCh:
		if err != nil {
			log.Printf("Client error: %v", err)
		}
	}

	log.Println("Client shutdown complete")
}
