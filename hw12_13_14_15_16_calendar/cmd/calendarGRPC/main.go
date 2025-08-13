package main

import (
	"flag"
	"log"
	"net"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/calendarGRPC"
	calendarpb "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/calendarGRPC/pb"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/app"

	//"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/config"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := LoadConfig(configFile)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	logg := logger.New(cfg.Logger.Level)
	application := app.NewWithConfig(cfg, logg)

	lis, err := net.Listen("tcp", cfg.GRPC.ListenGrpc)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcSrv := grpc.NewServer()
	calendarpb.RegisterCalendarServiceServer(grpcSrv, calendarGRPC.NewEventServer(application))
	// Enable reflection so grpcurl/Postman can inspect services
	reflection.Register(grpcSrv)
	log.Printf("gRPC server listening on %s", cfg.GRPC.ListenGrpc)
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
