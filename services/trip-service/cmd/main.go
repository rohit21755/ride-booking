package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/infrastructure/grpc"
	"ride-sharing/services/trip-service/infrastructure/repository"
	"ride-sharing/services/trip-service/service"
	"syscall"

	grpcserver "google.golang.org/grpc"
)

const (
	GrpcAddr = ":9093"
)

func main() {
	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
		<-sigch
		cancel()

	}()

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpcserver.NewServer()
	//TODO initialize our grpc server
	grpc.NewGRPCHandler(grpcServer, svc)
	log.Printf("server started")

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// wait for the shudown

	<-ctx.Done()
	grpcServer.GracefulStop()
}
