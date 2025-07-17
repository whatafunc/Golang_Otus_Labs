package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/app"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage/memory"
	redisstorage "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage/redis"
	postgresstorage "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := LoadConfig(configFile)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	logg := logger.New(config.Logger.Level)

	var storage app.Storage
	switch config.Storage.Type {
	case "memory":
		storage = memorystorage.New()
	case "redis":
		storage = redisstorage.New(config.Storage.Redis)
	case "postgres":
		storage = postgresstorage.New(config.Storage.Postgres)
	default:
		logg.Error("unknown storage type: " + config.Storage.Type)
		os.Exit(1)
	}
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, config.HTTP.Listen)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")
	logg.Error("No error for calendar is running...")
	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
