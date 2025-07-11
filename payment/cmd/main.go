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

	paymentV1API "github.com/kont1n/MSA_Rocket_Factory/payment/internal/api/payment/v1"
	paymentService "github.com/kont1n/MSA_Rocket_Factory/payment/internal/service/payment"
	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

func main() {
	log.Printf("Payment service starting...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer()

	// Регистрируем сервис
	service := paymentService.NewService()
	api := paymentV1API.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(s, api)

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
