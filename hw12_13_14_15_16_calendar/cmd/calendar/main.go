//nolint:depguard // not finished app
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

	cfg, err := LoadConfig(configFile)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	logg := logger.New(cfg.Logger.Level)

	calendar := app.NewWithConfig(cfg, logg)

	server := internalhttp.NewServer(logg, calendar, cfg.HTTP.Listen)

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
	// logg.Error("No error for calendar is running...")
	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
