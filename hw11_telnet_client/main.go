package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout (e.g. 10s, 1m)")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--timeout=duration] host port\n", os.Args[0])
		os.Exit(1)
	}

	host := args[0]
	port := args[1]
	address := host + ":" + port

	// Create context that cancels on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Printf("Failed to connect: %v", err)
	}
	defer client.Close()
	errCh := make(chan error, 2)

	// Run Receive with context in a goroutine
	go func() {
		errCh <- client.Receive(ctx)
	}()

	// Run Send with context in a goroutine
	go func() {
		errCh <- client.Send(ctx)
	}()

	// Wait for either Send or Receive to return an error or finish
	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			log.Printf("Client error: %v", err)
			break
		}
	}

	log.Println("Client shutdown complete")
}
