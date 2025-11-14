package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/rabbit"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Build AMQP URI from config
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	consumer, err := rabbit.NewConsumer(amqpURI, cfg.Queue, cfg.ConsumerTag)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}
	defer consumer.Shutdown()

	quit := make(chan struct{})
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		close(quit)
	}()

	consumer.Start(quit)
}
