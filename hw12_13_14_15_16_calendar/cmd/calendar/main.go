package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	calendarpb "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/calendarGRPC/pb"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/app"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/logger"
	calendarGRPC "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/server/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to config file")
}

func main() {
	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}
	cfg, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	logg := logger.New(cfg.Logger.Level)

	appInstance := app.NewWithConfig(cfg, logg)

	grpcAddr := cfg.GRPC.ListenGrpc
	httpAddr := cfg.HTTP.Listen

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logg.Error("failed to listen for gRPC: " + err.Error())
			cancel()
			return
		}

		grpcServer := calendarGRPC.NewGRPCServer(appInstance, logg)
		reflection.Register(grpcServer)

		logg.Info("gRPC server listening on " + grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			logg.Error("failed to serve gRPC: " + err.Error())
		}
	}()

	// Start HTTP gateway server
	go func() {
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		if err := calendarpb.RegisterCalendarServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
			logg.Error("failed to start HTTP gateway: " + err.Error())
			cancel()
			return
		}

		logg.Info("HTTP gateway listening on " + httpAddr)
		srv := &http.Server{
			Addr:         httpAddr,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		if err := srv.ListenAndServe(); err != nil {
			logg.Error("failed to serve HTTP: " + err.Error())
		}
	}()

	// Gracefully handle termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	logg.Info("Shutdown signal received, exiting")
	cancel()
	time.Sleep(time.Second) // Wait briefly for cleanup if needed
}
