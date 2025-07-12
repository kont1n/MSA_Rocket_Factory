package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryV1API "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/api/v1"
	inventoryRepository "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/inmemory"
	inventoryService "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/service/part"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func main() {
	log.Printf("Inventory service starting...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer()

	// Регистрируем сервис
	repo := inventoryRepository.NewRepository()
	service := inventoryService.NewService(repo)
	api := inventoryV1API.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
