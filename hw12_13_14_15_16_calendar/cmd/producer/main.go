package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/app"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/rabbit"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()
	logg := logger.New("producer")

	cfg, err := LoadConfig(configFile)
	if err != nil {
		logg.Error("failed to load config: " + err.Error())
		//fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}
	logg.Info(fmt.Sprintf("✅ Config loaded: %+v", cfg))

	appInstance := app.NewWithConfig(cfg.Config, logg)
	if appInstance == nil {
		logg.Error("application is not initialized")
		return
	}
	logg.Info("✅ app started\n")
	// Build AMQP URI from config
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.Rabbit.User,
		cfg.Rabbit.Password,
		cfg.Rabbit.Host,
		cfg.Rabbit.Port,
	)

	producer, err := rabbit.NewProducer(appInstance, amqpURI, "", cfg.Rabbit.Queue)
	if err != nil {
		logg.Error("failed to load config: " + err.Error())
	}
	defer producer.Shutdown()

	quit := make(chan struct{})
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		close(quit)
	}()

	producer.Start(quit)
}
